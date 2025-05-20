package repository

import (
	"context"

	"g-management/internal/models/payments/pkg/entity"

	"gorm.io/gorm"
)

type PaymentsRepositoryInterface interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Payments, error)
}

type paymentsRepository struct {
	DB *gorm.DB
}

func NewPaymentsRepository(db *gorm.DB) PaymentsRepositoryInterface {
	return &paymentsRepository{
		DB: db,
	}
}

// en: TakeByConditions function to get a payment by conditions
func (p *paymentsRepository) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) (entity.Payments, error) {
	var payment entity.Payments
	pdb := p.DB.WithContext(ctx)
	err := pdb.Model(&payment).Where(conditions).Take(&entity.Payments{}).Error
	return payment, err
}
