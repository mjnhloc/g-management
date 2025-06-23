package mount

import (
	"fmt"

	"g-management/internal/services/pkg/container"
	"g-management/internal/services/pkg/container/graphql"
	"g-management/internal/services/router"
	"g-management/pkg/shared/middleware"
	"g-management/pkg/shared/validator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func MountAll(
	repositories *container.RepositoryContainers,
	services *container.ServiceContainers,
	ginServer *gin.Engine,
	db *gorm.DB,
) error {
	inputValidator, err := validator.NewJsonSchemaValidator()
	if err != nil {
		return fmt.Errorf("failed to create a JSON schema validator: %w", err)
	}

	graphql, err := graphql.NewGraphqlSchema(repositories, services, db)
	if err != nil {
		return fmt.Errorf("failed to create a new GraphQL schema: %w", err)
	}

	routerClass := ginServer.Group("/classes")
	routerClass.Use(middleware.CheckAuthentication())
	routerMember := ginServer.Group("/members")
	routerMember.Use(middleware.CheckAuthentication())
	routerTrainer := ginServer.Group("/trainers")
	routerTrainer.Use(middleware.CheckAuthentication())

	handlerContainer := container.NewHandlerContainer(inputValidator, graphql, db)

	router.BindClassRoutes(routerClass, handlerContainer.Classes)

	router.BindMemberRoutes(routerMember, handlerContainer.Members)

	router.BindTrainerRoutes(routerTrainer, handlerContainer.Trainers)

	return nil
}
