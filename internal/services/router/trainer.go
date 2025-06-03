package router

import (
	"g-management/internal/services/handler/trainer"

	"github.com/gin-gonic/gin"
)

func BindTrainerRoutes(
	router *gin.RouterGroup,
	handler *trainer.HTTPHandler,
) {
	router.GET("/", handler.GetAllTrainers)
	router.GET("/:id", handler.GetTrainerDetails)
	router.POST("/", handler.PostNewTrainer)
}
