package container

import (
	"g-management/internal/models/memberships/pkg/repository"

	"gorm.io/gorm"
)

type RepositoryContainer struct {
	MembershipsRepository repository.MembershipsRepositoryInterface
}

func NewRepositoryContainer(db *gorm.DB) RepositoryContainer {
	return RepositoryContainer{
		MembershipsRepository: repository.NewMembershipsRepository(db),
	}
}
