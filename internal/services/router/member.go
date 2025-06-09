package router

import (
	"g-management/internal/services/handler/member"

	"github.com/gin-gonic/gin"
)

func BindMemberRoutes(
	router *gin.RouterGroup,
	handler *member.HTTPHandler,
) {
	router.GET("/", handler.GetAllMembers)
	router.GET("/:id", handler.GetMemberDetails)
	router.POST("/", handler.PostNewMember)
	router.PUT("/:id", handler.PutMemberInfo)
	router.DELETE("/:id", handler.DeleteMember)
}
