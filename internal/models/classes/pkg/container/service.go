package container

import (
	"g-management/internal/models/classes/pkg/services"
	"g-management/pkg/services/elasticsearch/client"
)

type ServiceContainer struct {
	ClassesService services.ClassesServiceInterface
}

func NewServiceContainer(esClient client.ClientInterface) ServiceContainer {
	return ServiceContainer{
		ClassesService: services.NewClassesService(esClient),
	}
}
