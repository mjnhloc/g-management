package validator

import (
	"fmt"
	"os"
	"path/filepath"

	"g-management/pkg/shared/utils"

	"github.com/xeipuuv/gojsonschema"
)

var cachedGoJsonSchemaInstances = map[string]*gojsonschema.Schema{}

type JsonSchemaValidator struct {
	basePath string
	schemas  map[string]*gojsonschema.Schema
}

func NewJsonSchemaValidator() (*JsonSchemaValidator, error) {
	validator := &JsonSchemaValidator{
		basePath: os.Getenv("GM_SCHEMAS_PATH"),
		schemas:  make(map[string]*gojsonschema.Schema),
	}

	gojsonschema.FormatCheckers.Add("not-fullwidth-or-halfwidth", NotFullwidthOrHalfwidthFormatChecker{})
	gojsonschema.FormatCheckers.Add("date-time", NonStandardDateTimeFormatChecker{})
	gojsonschema.FormatCheckers.Add("password", PasswordChecker{})
	gojsonschema.FormatCheckers.Add("strong-password", StrongPasswordChecker{})
	gojsonschema.FormatCheckers.Add("auth0-password", Auth0PasswordChecker{})
	gojsonschema.FormatCheckers.Add("domain", DomainChecker{})
	gojsonschema.FormatCheckers.Add("hiragana", HiraganaChecker{})
	gojsonschema.FormatCheckers.Add("google_analytics", GaMeasurementIDChecker{})
	gojsonschema.FormatCheckers.Add("google_tag_manager", GTMChecker{})
	gojsonschema.FormatCheckers.Add("string_with_max_length", MaxLengthChecker{})
	gojsonschema.FormatCheckers.Add("url", UrlChecker{})
	gojsonschema.FormatCheckers.Add("id-sns", IDSnSChecker{})

	err := validator.loadDirSchemas("")
	if err != nil {
		return nil, err
	}

	return validator, nil
}

func (validator *JsonSchemaValidator) loadDirSchemas(path string) error {
	schemasFiles, err := os.ReadDir(validator.basePath + path)
	if err != nil {
		return err
	}

	for _, schemaFile := range schemasFiles {
		if schemaFile.Name() == ".gitkeep" {
			continue
		}

		schemaPath := path + "/" + schemaFile.Name()
		if schemaFile.IsDir() {
			if err := validator.loadDirSchemas(schemaPath); err != nil {
				return err
			}
			continue
		}

		absPath, err := filepath.Abs(validator.basePath + schemaPath)
		if err != nil {
			return err
		}

		goJsonSchemaPath := "file://" + filepath.ToSlash(absPath)
		// goJsonSchemaPath := "file://" + validator.basePath + schemaPath
		var schema *gojsonschema.Schema
		var schemaExists bool
		if schema, schemaExists = cachedGoJsonSchemaInstances[goJsonSchemaPath]; !schemaExists {
			schemaLoader := gojsonschema.NewReferenceLoader(goJsonSchemaPath)
			schema, err = gojsonschema.NewSchema(schemaLoader)
			if err != nil {
				return err
			}

			cachedGoJsonSchemaInstances[goJsonSchemaPath] = schema
		}

		validator.schemas[schemaPath] = schema
	}

	return nil
}

func (validator *JsonSchemaValidator) Validate(
	schemaFile string,
	data interface{},
) (*gojsonschema.Result, error) {
	if schemaFile[0] != '/' {
		schemaFile = "/" + schemaFile
	}

	schema, schemaExists := validator.schemas[schemaFile]
	if !schemaExists {
		return nil, fmt.Errorf("the schema '%v' was not found for json validation", schemaFile)
	}

	dataLoader := gojsonschema.NewGoLoader(data)

	result, err := schema.Validate(dataLoader)
	if err != nil {
		return nil, err
	}
	if len(result.Errors()) == 0 {
		return nil, nil
	}

	return result, nil
}

func (validator *JsonSchemaValidator) GetErrorDetails(
	result gojsonschema.ResultError,
) map[string]interface{} {
	return map[string]interface{}{
		"context":     result.Context(),
		"description": result.Description(),
		"details":     result.Details(),
		"field":       result.Field(),
		"type":        result.Type(),
		"value":       result.Value(),
	}
}

func (validator *JsonSchemaValidator) GetErrorField(
	result gojsonschema.ResultError,
) string {
	field := result.Field()
	errorDetails := result.Details()
	if property, propertyExists := errorDetails["property"]; propertyExists {
		if propertyString, propertyIsString := property.(string); propertyIsString {
			field = field + "." + propertyString
		}
	}

	return field
}

func (validator *JsonSchemaValidator) GetCustomErrorMessage(
	result gojsonschema.ResultError,
) string {
	details := result.Details()
	format, formatExists := details["format"]
	minValue, minExists := details["min"]
	if result.Type() == "format" && formatExists {
		if format == "email" || format == "idn-email" {
			return utils.ErrorEmailFail
		}
		if format == "password" || format == "strong-password" || format == "auth0-password" {
			return utils.ErrorPasswordFail
		}
		if format == "google_analytics" {
			return utils.ErrorGAInvalid
		}
		if format == "google_tag_manager" {
			return utils.ErrorGTMInvalid
		}
		if format == "string_with_max_length" {
			return utils.ErrCheckMaxLengthUnder50Characters
		}
		if format == "domain" {
			return utils.ErrInvalidDomain
		}
		if format == "hiragana" {
			return utils.ErrorInvalidHiragana
		}
	}

	if result.Type() == "required" {
		return "MSGCM001"
	}

	if result.Type() == "string_gte" && minExists {
		if minValue == 1 {
			return utils.ErrorInputFail // Non-empty string
		}
	}

	_, maxExists := details["max"]
	if result.Type() == "string_lte" && maxExists {
		return utils.ErrorInputCharacterLimit
	}

	if result.Type() == maxLengthByte {
		return utils.ErrorInputByteLimit
	}

	return utils.ErrorInputFail
}
