package member

import (
	"net/http"
	"strconv"

	baseDto "g-management/internal/services/pkg/dto"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *HTTPHandler) DeleteMember(c *gin.Context) {
	memberID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.SetBadRequestErrorResponse(c, map[string]string{
			"id": "Invalid member ID format",
		})
		return
	}

	result := graphql.Do(graphql.Params{
		Schema: h.graphql,
		VariableValues: map[string]interface{}{
			"id": memberID,
		},
		Context: c,
		RequestString: `
			mutation ($id: BigInt!) {
				delete_member (id: $id) 
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
