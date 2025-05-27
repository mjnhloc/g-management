package mount

import (
	"fmt"

	"g-management/internal/services/pkg/container"
	"g-management/internal/services/router"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func MountAll(
	repositories *container.RepositoryContainers,
	ginServer *gin.Engine,
	db *gorm.DB,
) error {
	graphql, err := container.NewGraphqlSchema(repositories, db)
	if err != nil {
		return fmt.Errorf("failed to create a new GraphQL schema: %w", err)
	}

	routerClass := ginServer.Group("/class")
	routerMember := ginServer.Group("/member")
	routerTrainer := ginServer.Group("/trainer")

	handlerContainer := container.NewHandlerContainer(graphql, db)

	router.BindClassRoutes(routerClass, handlerContainer.Classes)

	router.BindMemberRoutes(routerMember, handlerContainer.Members)

	router.BindTrainerRoutes(routerTrainer, handlerContainer.Trainers)

	return nil
}
