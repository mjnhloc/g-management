package router

import (
	"g-management/internal/services/handler/trainer"
	"g-management/pkg/shared/middleware"

	"github.com/gin-gonic/gin"
)

func BindTrainerRoutes(
	router *gin.RouterGroup,
	handler *trainer.HTTPHandler,
) {
	router.GET("/", handler.GetAllTrainers)

	router.Use(middleware.RequireRole("admin"))
	{
		router.GET("/:id", handler.GetTrainerDetails)
		router.POST("/", handler.PostNewTrainer)
		router.PUT("/:id", handler.PutTrainerInfo)
		router.DELETE("/:id", handler.DeleteTrainer)
	}
}
