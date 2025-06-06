package class

import (
	"net/http"
	"strconv"

	baseDto "g-management/internal/services/pkg/dto"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *HTTPHandler) PutClassInfo(c *gin.Context) {
	classID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.SetBadRequestErrorResponse(c, map[string]string{
			"id": "Invalid class ID format",
		})
		return
	}

	input, err := h.GetInputsAsMap(c)
	if err != nil {
		h.SetGenericErrorResponse(c, err)
		return
	}

	validationResult, err := h.Validator.Validate(PutClass, input)
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
		VariableValues: map[string]interface{}{
			"id": classID,
		},
		Context: c,
		RequestString: `
			mutation ($id: BigInt!) {
				class: put_class (id: $id) {
					id
					name
					schedule
					duration
					max_capacity
					description
					trainer {
						id
						name
						specialization	
					}
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
