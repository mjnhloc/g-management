package router

import (
	"g-management/internal/services/handler/class"

	"github.com/gin-gonic/gin"
)

func BindClassRoutes(
	router *gin.RouterGroup,
	handler *class.HTTPHandler,
) {
	router.GET("/", handler.GetAllClasses)
	router.GET("/:id", handler.GetClassDetails)
	router.POST("/", handler.PostNewClass)
	router.PUT("/:id", handler.PutClassInfo)
	router.DELETE("/:id", handler.DeleteClass)
}
