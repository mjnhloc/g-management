package member

import (
	baseDto "g-management/internal/services/pkg/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *HTTPHandler) GetAllMembers(c *gin.Context) {
	result := graphql.Do(graphql.Params{
		Schema:  h.graphql,
		Context: c,
		RequestString: `
			query {
				members: get_all_members {
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
