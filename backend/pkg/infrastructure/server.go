package infrastructure

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewServer(db *gorm.DB) *gin.Engine {
	router := gin.New()

	router.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "Welcome to the g-management API",
		})
	})

	return router
}
