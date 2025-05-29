package trainer

import (
	baseDto "g-management/internal/services/pkg/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *HTTPHandler) GetAllTrainers(c *gin.Context) {
	result := graphql.Do(graphql.Params{
		Schema:  h.graphql,
		Context: c,
		RequestString: `
			query {
				trainer: get_all_trainers {
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
