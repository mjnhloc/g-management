package class

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"g-management/internal/models/classes/pkg/entity"
	baseDto "g-management/internal/services/pkg/dto"
	"g-management/pkg/services"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

type ClassUpdateNotification struct {
	Type      string         `json:"type"`
	Action    string         `json:"action"`
	ClassID   int            `json:"class_id"`
	ClassInfo entity.Classes `json:"class_info"`
}

func (h *HTTPHandler) PutClassInfo(c *gin.Context) {
	classID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.SetBadRequestErrorResponse(c, map[string]string{
			"id": "Invalid class ID format",
		})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		h.SetBadRequestErrorResponse(c, map[string]string{
			"body": "Invalid request body",
		})
		return
	}

	// Validate input
	validationResult, err := h.Validator.Validate("put_class_info", input)
	if err != nil {
		h.SetInternalErrorResponse(c, err)
		return
	}
	if validationResult != nil {
		h.SetJSONValidationErrorResponse(c, validationResult)
		return
	}

	// Update class info using GraphQL
	result := graphql.Do(graphql.Params{
		Schema: h.graphql,
		RequestString: `
			mutation ($id: BigInt!, $class: ClassInput!) {
				put_class_info(id: $id) {
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
		VariableValues: map[string]interface{}{
			"id":    classID,
			"class": input,
		},
		Context: c.Request.Context(),
	})

	if len(result.Errors) > 0 {
		h.SetGenericErrorResponse(c, result.Errors[0])
		return
	}

	// Get Redis client from context
	redisClient, exists := c.Get("redisClient")
	if !exists {
		h.SetInternalErrorResponse(c, err)
		return
	}

	// Create notification
	notification := ClassUpdateNotification{
		Type:      "class_update",
		Action:    "update",
		ClassID:   classID,
		ClassInfo: result.Data.(map[string]interface{})["put_class_info"].(entity.Classes),
	}

	// Publish notification to Redis
	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		c.Error(err) // Log error but don't fail request
	} else {
		err = redisClient.(*services.RedisClient).PublishEvent(
			context.Background(),
			"class_updates",
			string(notificationJSON),
		)
		if err != nil {
			c.Error(err) // Log error but don't fail request
		}
	}

	c.JSON(http.StatusOK, &baseDto.BaseSuccessResponse{
		Data: result.Data,
	})
}
