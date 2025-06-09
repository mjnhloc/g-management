package trainer

import (
	"net/http"

	baseDto "g-management/internal/services/pkg/dto"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *HTTPHandler) PostNewTrainer(c *gin.Context) {
	input, err := h.GetInputsAsMap(c)
	if err != nil {
		h.SetGenericErrorResponse(c, err)
		return
	}

	validation, err := h.Validator.Validate(PostNewTrainer, input)
	if err != nil {
		h.SetInternalErrorResponse(c, err)
		return
	}
	if validation != nil {
		h.SetJSONValidationErrorResponse(c, validation)
		return
	}

	result := graphql.Do(graphql.Params{
		Schema:     h.graphql,
		RootObject: input,
		Context:    c,
		RequestString: `
			mutation {
				trainer: post_new_trainer {
					id
					name
					email
					phone
					specialization
					hired_at
				}
			}
		`,
	})
	if result.HasErrors() {
		h.SetGenericErrorResponse(c, result.Errors[0])
		return
	}

	c.JSON(http.StatusOK, &baseDto.BaseSuccessResponse{
		Data: result.Data,
	})
}
