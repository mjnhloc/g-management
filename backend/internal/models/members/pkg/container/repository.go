package container

import (
	"g-management/internal/models/members/pkg/repository"

	"gorm.io/gorm"
)

type RepositoryContainer struct {
	MembersRepository repository.MembersRepositoryInterface
}

func NewRepositoryContainer(db *gorm.DB) RepositoryContainer {
	return RepositoryContainer{
		MembersRepository: repository.NewMembersRepository(db),
	}
}
