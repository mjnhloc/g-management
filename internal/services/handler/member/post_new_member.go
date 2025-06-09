package member

import (
	baseDto "g-management/internal/services/pkg/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *HTTPHandler) PostNewMember(c *gin.Context) {
	input, err := h.GetInputsAsMap(c)
	if err != nil {
		h.SetGenericErrorResponse(c, err)
		return
	}

	validation, err := h.Validator.Validate(PostNewMember, input)
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
				member: post_new_member {
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
