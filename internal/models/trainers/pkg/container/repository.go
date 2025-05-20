package container

import (
	"g-management/internal/models/trainers/pkg/repository"

	"gorm.io/gorm"
)

type RepositoryContainer struct {
	TrainersRepository repository.TrainersRepositoryInterface
}

func NewRepositoryContainer(db *gorm.DB) RepositoryContainer {
	return RepositoryContainer{
		TrainersRepository: repository.NewTrainersRepository(db),
	}
}
