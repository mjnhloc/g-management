package class

import (
	"net/http"

	baseDto "g-management/internal/services/pkg/dto"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *HTTPHandler) GetSearchClasses(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		h.SetBadRequestErrorResponse(c, map[string]string{
			"keyword": "Keyword is required for search",
		})
		return
	}

	result := graphql.Do(graphql.Params{
		Schema: h.graphql,
		VariableValues: map[string]interface{}{
			"keyword": keyword,
		},
		Context: c,
		RequestString: `
			query ($keyword: String!) {
				classes: get_search_classes (keyword: $keyword) {
					id
					name
					schedule
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
