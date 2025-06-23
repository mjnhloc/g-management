package trainer

import (
	"net/http"
	"strconv"

	baseDto "g-management/internal/services/pkg/dto"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *HTTPHandler) PutTrainerInfo(c *gin.Context) {
	trainerID, err := strconv.Atoi(c.Param("trainer_id"))
	if err != nil {
		h.SetBadRequestErrorResponse(c, map[string]string{
			"id": "Invalid trainer ID format",
		})
		return
	}

	input, err := h.GetInputsAsMap(c)
	if err != nil {
		h.SetGenericErrorResponse(c, err)
		return
	}

	validation, err := h.Validator.Validate(PutTrainerInfo, input)
	if err != nil {
		h.SetInternalErrorResponse(c, err)
		return
	}
	if validation != nil {
		h.SetJSONValidationErrorResponse(c, validation)
		return
	}

	result := graphql.Do(graphql.Params{
		Schema: h.graphql,
		VariableValues: map[string]interface{}{
			"id": trainerID,
		},
		RootObject: input,
		Context:    c,
		RequestString: `
			mutation ($id: BigInt!) {
				trainer: put_trainer_info (id: $id) {
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
