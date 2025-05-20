package container

import (
	"g-management/internal/models/payments/pkg/repository"

	"gorm.io/gorm"
)

type RepositoryContainer struct {
	PaymentsRepository repository.PaymentsRepositoryInterface
}

func NewRepositoryContainer(db *gorm.DB) RepositoryContainer {
	return RepositoryContainer{
		PaymentsRepository: repository.NewPaymentsRepository(db),
	}
}
