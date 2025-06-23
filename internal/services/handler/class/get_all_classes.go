package class

import (
	baseDto "g-management/internal/services/pkg/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *HTTPHandler) GetAllClasses(c *gin.Context) {
	result := graphql.Do(graphql.Params{
		Schema:  h.graphql,
		Context: c,
		RequestString: `
			query {
				classes: get_all_classes {
					id
					name
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
