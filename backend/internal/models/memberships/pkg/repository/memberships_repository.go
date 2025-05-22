package repository

import (
	"context"

	"g-management/internal/models/memberships/pkg/entity"

	"gorm.io/gorm"
)

type MembershipsRepositoryInterface interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Memberships, error)
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
	err := cdb.Model(&membership).Where(conditions).Take(&entity.Memberships{}).Error
	return membership, err
}
