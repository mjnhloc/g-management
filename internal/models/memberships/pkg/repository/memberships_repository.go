package repository

import (
	"context"

	"g-management/internal/models/memberships/pkg/entity"
	"g-management/pkg/shared/utils"

	"gorm.io/gorm"
)

type MembershipsRepositoryInterface interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Memberships, error)
	CreateWithTransaction(tx *gorm.DB, attributes map[string]interface{}) (entity.Memberships, error)
}

type membershipsRepository struct {
	DB *gorm.DB
}

func NewMembershipsRepository(db *gorm.DB) MembershipsRepositoryInterface {
	return &membershipsRepository{
		DB: db,
	}
}

// en: TakeByConditions function to get a member by conditions
func (m *membershipsRepository) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) (entity.Memberships, error) {
	var membership entity.Memberships
	cdb := m.DB.WithContext(ctx)
	err := cdb.Where(conditions).Take(&membership).Error
	return membership, err
}

// en: CreateWithTransaction function to create a new membership with given attributes within a transaction
func (m *membershipsRepository) CreateWithTransaction(
	tx *gorm.DB,
	attributes map[string]interface{},
) (entity.Memberships, error) {
	var membership entity.Memberships
	err := utils.MapToStruct(attributes, &membership)
	if err != nil {
		return entity.Memberships{}, err
	}

	err = tx.Create(&membership).Error
	return membership, err
}
