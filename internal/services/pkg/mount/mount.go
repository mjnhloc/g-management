package mount

import (
	"g-management/internal/services/pkg/container"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func MountAll(
	repository *container.RepositoryContainers,
	ginServer *gin.Engine,
	db *gorm.DB,
) {
}
