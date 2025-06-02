package repository

import (
	"context"

	"g-management/internal/models/members/pkg/entity"
	"g-management/pkg/shared/utils"

	"gorm.io/gorm"
)

type MembersRepositoryInterface interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Members, error)
	Create(ctx context.Context, attributes map[string]interface{}) (entity.Members, error)
	FindByConditions(ctx context.Context, conditions map[string]interface{}) ([]entity.Members, error)
	CreateWithTransaction(tx *gorm.DB, attributes map[string]interface{}) (entity.Members, error)
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
	err := cdb.Where(conditions).Take(&member).Error
	return member, err
}

// en: Create function to create a new member with given attributes
func (m *membersRepository) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.Members, error) {
	var member entity.Members
	err := utils.MapToStruct(attributes, &member)
	if err != nil {
		return entity.Members{}, err
	}

	cdb := m.DB.WithContext(ctx)
	err = cdb.Create(&member).Error
	return member, err
}

// en: FindByConditions function to find members by conditions
func (m *membersRepository) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) ([]entity.Members, error) {
	var members []entity.Members
	cdb := m.DB.WithContext(ctx)
	err := cdb.Where(conditions).Find(&members).Error
	return members, err
}

// en: CreateWithTransaction function to create a new member with given attributes within a transaction
func (m *membersRepository) CreateWithTransaction(
	tx *gorm.DB,
	attributes map[string]interface{},
) (entity.Members, error) {
	var member entity.Members
	err := utils.MapToStruct(attributes, &member)
	if err != nil {
		return entity.Members{}, err
	}

	err = tx.Create(&member).Error
	return member, err
}
