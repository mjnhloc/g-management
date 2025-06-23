package handler

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"g-management/pkg/log"
	"g-management/pkg/shared/validator"
	"g-management/pkg/shared/wraperror"

	baseDto "g-management/pkg/dto"

	baseUtils "g-management/pkg/shared/utils"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/xeipuuv/gojsonschema"
	"gorm.io/gorm"
)

// BaseHTTPHandler base handler struct.
type BaseHTTPHandler struct {
	Validator *validator.JsonSchemaValidator
}

// ResponseCSV func
func (h BaseHTTPHandler) ResponseCSV(c *gin.Context, statusCode int, fileName string, data []byte) {
	c.Writer.Header().Set("Content-Description", "File Transfer")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename="+fileName)
	c.Data(statusCode, "text/csv", data)
}

// ResponseZIP func
func (h BaseHTTPHandler) ResponseZIP(c *gin.Context, statusCode int, fileName string, mapContentFile map[string]bytes.Buffer) error {
	c.Writer.Header().Set("Content-Type", "application/zip")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename="+fileName)

	zipW := zip.NewWriter(c.Writer)
	defer func() {
		if err := zipW.Close(); err != nil {
			log.Error(c, baseUtils.ErrorCloseWriter, "error", err)
		}
	}()
	for key, contentFile := range mapContentFile {
		f, err := zipW.Create(key)
		if err != nil {
			return err
		}
		_, err = f.Write(contentFile.Bytes())
		if err != nil {
			return err
		}
	}
	return nil
}

func (h BaseHTTPHandler) GetInputsAsMap(c *gin.Context) (map[string]interface{}, error) {
	contentType := c.ContentType()
	if contentType != "application/json" {
		return nil, wraperror.NewApiDisplayableError(
			http.StatusBadRequest,
			"Expected 'application/json' for Content-Type, got '"+contentType+"'",
			nil,
		)
	}

	// Getting the body as a map
	input := make(map[string]interface{})
	err := c.ShouldBindJSON(&input)
	if err != nil {
		return nil, err
	}

	return input, nil
}

func (h *BaseHTTPHandler) SetGenericErrorResponse(c *gin.Context, finalError error) {
	// Retrieving the original error inside GraphQL's wrapper if there is one
	// If there is none, we keep the error coming from the graphql's engine
	originalError := finalError
	if _, ok := originalError.(gqlerrors.FormattedError); ok {
		err := originalError.(gqlerrors.FormattedError).OriginalError()
		if err != nil {
			originalError = err
		}

		if _, ok := originalError.(*gqlerrors.Error); ok {
			err := originalError.(*gqlerrors.Error).OriginalError
			if err != nil {
				originalError = err
			}
		}
	}

	apiError := &wraperror.ApiDisplayableError{}
	jsonError := &json.SyntaxError{}
	switch {
	case errors.As(originalError, &apiError):
		var debugInfo interface{}
		data := baseDto.BaseErrorResponse{
			Error: &baseDto.ErrorResponse{
				Message:          apiError.Message(),
				DebugInformation: debugInfo,
			},
		}
		c.JSON(apiError.HttpStatus(), data)
		return
	case errors.Is(originalError, gorm.ErrRecordNotFound) || originalError.Error() == gorm.ErrRecordNotFound.Error():
		data := baseDto.BaseErrorResponse{
			Error: &baseDto.ErrorResponse{
				Message: originalError.Error(),
			},
		}
		c.JSON(http.StatusNotFound, data)
		return
	case errors.As(originalError, &jsonError):
		data := &baseDto.BaseErrorResponse{
			Error: &baseDto.ErrorResponse{
				Message: "Invalid json",
				Details: map[string]interface{}{
					"offset": jsonError.Offset,
					"error":  jsonError.Error(),
				},
			},
		}

		c.JSON(http.StatusBadRequest, data)
		return
	default:
		h.SetInternalErrorResponse(c, finalError)
		return
	}
}

func (h *BaseHTTPHandler) SetValidationErrorResponse(c *gin.Context, err error) {
	data := &baseDto.BaseErrorResponse{
		Error: &baseDto.ErrorResponse{
			Message: err,
		},
	}

	c.JSON(http.StatusBadRequest, data)
}

func (h *BaseHTTPHandler) SetJSONValidationErrorResponse(
	c *gin.Context,
	validationResults *gojsonschema.Result,
) {
	h.SetJSONValidationWithCustomErrorResponse(
		c,
		validationResults,
		func(result gojsonschema.ResultError) string {
			return ""
		},
	)
}

func (h *BaseHTTPHandler) SetJSONValidationWithCustomErrorResponse(
	c *gin.Context,
	validationResults *gojsonschema.Result,
	getError func(result gojsonschema.ResultError) string,
) {
	messages := map[string]string{}
	details := make([]map[string]interface{}, 0)
	for _, validationError := range validationResults.Errors() {
		field := h.Validator.GetErrorField(validationError)
		detail := h.Validator.GetErrorDetails(validationError)

		// Getting either a message defined by the handler,
		// or the default customized one
		message := getError(validationError)
		if message == "" {
			message = h.Validator.GetCustomErrorMessage(validationError)
		}

		messages[field] = message
		details = append(details, detail)
	}

	data := &baseDto.BaseErrorResponse{
		Error: &baseDto.ErrorResponse{
			Message: messages,
			Details: details,
		},
	}

	c.JSON(http.StatusBadRequest, data)
}

func (h *BaseHTTPHandler) SetBadRequestErrorResponse(c *gin.Context, messages interface{}) {
	data := &baseDto.BaseErrorResponse{
		Error: &baseDto.ErrorResponse{
			Message: messages,
		},
	}

	c.JSON(http.StatusBadRequest, data)
}

func (h *BaseHTTPHandler) SetCustomErrorAndDetailResponse(c *gin.Context, err error, details interface{}) {
	data := &baseDto.BaseErrorResponse{
		Error: &baseDto.ErrorResponse{
			Message: err.Error(),
			Details: details,
		},
	}

	c.JSON(http.StatusInternalServerError, data)
}

// This outputs a 500 error with a custom message (contrary to SetGenericErrorResponse that hides the real unhandled error)
func (h *BaseHTTPHandler) SetInternalErrorResponse(c *gin.Context, err error) {
	data := &baseDto.BaseErrorResponse{
		Error: &baseDto.ErrorResponse{
			Message:          baseUtils.MessageInternalServerError,
			DebugInformation: err.Error(),
		},
	}

	c.JSON(http.StatusInternalServerError, data)
}

// SetCookie comment.
// en: set cookie.
func (h *BaseHTTPHandler) SetCookie(c *gin.Context, cookieResults []map[string]interface{}, inputMaxAge *int) error {
	for _, cookie := range cookieResults {
		name := cookie["name"].(string)
		value := cookie["value"].(string)
		path := cookie["path"].(string)
		domain := cookie["domain"].(string)
		secure := cookie["secure"].(bool)
		httpOnly := cookie["http_only"].(bool)

		maxAge := baseUtils.DerefInt(inputMaxAge)
		if v, ok := cookie["max_age"].(int); ok && v != 0 {
			maxAge = v
		}

		c.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
	}

	return nil
}
