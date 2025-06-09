package member

import (
	"net/http"
	"strconv"

	baseDto "g-management/internal/services/pkg/dto"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *HTTPHandler) GetMemberDetails(c *gin.Context) {
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
			query ($id: BigInt!) {
				member: get_member_details (id: $id) {
					id
					name
					email
					phone
					date_of_birth
					is_active
					membership {
						id
						membership_type
						start_date
						end_date
						payment {
							id
							price
							payment_date
							payment_method
							status
						}
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
