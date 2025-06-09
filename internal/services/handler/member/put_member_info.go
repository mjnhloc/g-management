package member

import (
	"net/http"
	"strconv"

	baseDto "g-management/internal/services/pkg/dto"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *HTTPHandler) PutMemberInfo(c *gin.Context) {
	memberID, err := strconv.Atoi(c.Param("id"))
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

	validationResult, err := h.Validator.Validate(PutMemberInfo, input)
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
			"id": memberID,
		},
		Context: c,
		RequestString: `
			mutation ($id: BigInt!) {
				member: put_member_info (id: $id) {
					id
					name
					email
					phone
					date_of_birth
					is_active
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
