package repository

import (
	"context"

	"g-management/internal/models/members/pkg/entity"

	"gorm.io/gorm"
)

type MembersRepositoryInterface interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Members, error)
}

type membersRepository struct {
	DB *gorm.DB
}

func NewMembersRepository(db *gorm.DB) MembersRepositoryInterface {
	return &membersRepository{
		DB: db,
	}
}

// en: TakeByConditions function to get a member by conditions
func (m *membersRepository) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) (entity.Members, error) {
	var member entity.Members
	cdb := m.DB.WithContext(ctx)
	err := cdb.Model(&member).Where(conditions).Take(&entity.Members{}).Error
	return member, err
}
