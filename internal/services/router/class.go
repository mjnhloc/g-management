package router

import (
	"g-management/internal/services/handler/class"

	"github.com/gin-gonic/gin"
)

func BindClassRoutes(
	router *gin.RouterGroup,
	handler *class.HTTPHandler,
) {
	router.GET("/all", handler.GetAllClasses)
}
