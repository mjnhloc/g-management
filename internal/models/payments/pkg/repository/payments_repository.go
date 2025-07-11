package repository

import (
	"context"

	"g-management/internal/models/payments/pkg/entity"
	"g-management/pkg/shared/utils"

	"gorm.io/gorm"
)

type PaymentsRepositoryInterface interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Payments, error)
	CreateWithTransaction(tx *gorm.DB, attributes map[string]interface{}) (entity.Payments, error)
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
	err := pdb.Where(conditions).Take(&payment).Error
	return payment, err
}

// en: CreateWithTransaction function to create a new payment with given attributes within a transaction
func (p *paymentsRepository) CreateWithTransaction(
	tx *gorm.DB,
	attributes map[string]interface{},
) (entity.Payments, error) {
	var payment entity.Payments
	err := utils.MapToStruct(attributes, &payment)
	if err != nil {
		return entity.Payments{}, err
	}

	err = tx.Create(&payment).Error
	return payment, err
}
