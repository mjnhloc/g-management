package container

import (
	classesContainer "g-management/internal/models/classes/pkg/container"
	"g-management/pkg/services/elasticsearch/client"
)

type ServiceContainers struct {
	ClassesContainer classesContainer.ServiceContainer
}

func NewServiceContainers(client client.ClientInterface) *ServiceContainers {
	return &ServiceContainers{
		ClassesContainer: classesContainer.NewServiceContainer(client),
	}
}
