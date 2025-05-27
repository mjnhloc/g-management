package router

import (
	"g-management/internal/services/handler/trainer"

	"github.com/gin-gonic/gin"
)

func BindTrainerRoutes(
	router *gin.RouterGroup,
	handler *trainer.HTTPHandler,
) {
}
