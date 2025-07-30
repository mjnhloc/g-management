package container

import (
	"g-management/internal/models/classes/pkg/repository"

	"gorm.io/gorm"
)

type RepositoryContainer struct {
	ClassesRepository   repository.ClassesRepositoryInterface
	WaitlistsRepository repository.WaitlistsRepositoryInterface
}

func NewRepositoryContainer(db *gorm.DB) RepositoryContainer {
	return RepositoryContainer{
		ClassesRepository:   repository.NewClassesRepository(db),
		WaitlistsRepository: repository.NewWaitlistsRepository(db),
	}
}
