package trainer

import (
	"net/http"
	"strconv"

	baseDto "g-management/internal/services/pkg/dto"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *HTTPHandler) GetTrainerDetails(c *gin.Context) {
	trainerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.SetBadRequestErrorResponse(c, map[string]string{
			"id": "Invalid trainer ID format",
		})
		return
	}

	result := graphql.Do(graphql.Params{
		Schema: h.graphql,
		VariableValues: map[string]interface{}{
			"id": trainerID,
		},
		Context: c,
		RequestString: `
			query ($id: BigInt!) {
				trainer: get_trainer_details (id: $id) {
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
