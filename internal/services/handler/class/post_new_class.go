package class

import (
	baseDto "g-management/internal/services/pkg/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *HTTPHandler) PostNewClass(c *gin.Context) {
	input, err := h.GetInputsAsMap(c)
	if err != nil {
		h.SetGenericErrorResponse(c, err)
		return
	}

	validationResult, err := h.Validator.Validate(PostNewClass, input)
	if err != nil {
		h.SetInternalErrorResponse(c, err)
		return
	}

	if validationResult != nil {
		h.SetJSONValidationErrorResponse(c, validationResult)
		return
	}

	result := graphql.Do(graphql.Params{
		Schema:     h.graphql,
		RootObject: input,
		Context:    c,
		RequestString: `
			mutation {
				class: post_new_class {
					id
					name
					trainer {
						id
						name
						specialization	
					}
					schedule
					duration
					max_capacity
					description
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
