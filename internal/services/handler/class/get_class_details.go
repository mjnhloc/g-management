package class

import (
	baseDto "g-management/internal/services/pkg/dto"
	"g-management/pkg/shared/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *HTTPHandler) GetClassDetails(c *gin.Context) {
	classID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.SetBadRequestErrorResponse(c, map[string]string{
			"id": utils.ErrorInputFail,
		})
		return
	}

	result := graphql.Do(graphql.Params{
		Schema: h.graphql,
		VariableValues: map[string]interface{}{
			"id": classID,
		},
		Context: c,
		RequestString: `
			query ($id: BigInt!) {
				class: get_class_details (class_id: $id) {
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
