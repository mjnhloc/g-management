package mount

import (
	"fmt"

	"g-management/internal/services/pkg/container"
	"g-management/internal/services/router"
	"g-management/pkg/shared/validator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func MountAll(
	repositories *container.RepositoryContainers,
	ginServer *gin.Engine,
	db *gorm.DB,
) error {
	inputValidator, err := validator.NewJsonSchemaValidator()
	if err != nil {
		return fmt.Errorf("failed to create a JSON schema validator: %w", err)
	}

	graphql, err := container.NewGraphqlSchema(repositories, db)
	if err != nil {
		return fmt.Errorf("failed to create a new GraphQL schema: %w", err)
	}

	routerClass := ginServer.Group("/classes")
	routerMember := ginServer.Group("/members")
	routerTrainer := ginServer.Group("/trainers")

	handlerContainer := container.NewHandlerContainer(inputValidator, graphql, db)

	router.BindClassRoutes(routerClass, handlerContainer.Classes)

	router.BindMemberRoutes(routerMember, handlerContainer.Members)

	router.BindTrainerRoutes(routerTrainer, handlerContainer.Trainers)

	return nil
}
