package router

import (
	"g-management/internal/services/handler/class"
	"g-management/pkg/shared/middleware"

	"github.com/gin-gonic/gin"
)

func BindClassRoutes(
	router *gin.RouterGroup,
	handler *class.HTTPHandler,
) {
	router.GET("/", handler.GetAllClasses)
	router.GET("/search", handler.GetSearchClasses)

	router.Use(middleware.RequireRole("admin"))
	{
		router.GET("/:id", handler.GetClassDetails)
		router.POST("/", handler.PostNewClass)
		router.PUT("/:id", handler.PutClassInfo)
		router.DELETE("/:id", handler.DeleteClass)
	}
}
