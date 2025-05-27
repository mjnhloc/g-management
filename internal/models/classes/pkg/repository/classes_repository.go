package repository

import (
	"context"

	"g-management/internal/models/classes/pkg/entity"
	"g-management/pkg/shared/utils"

	"gorm.io/gorm"
)

type ClassesRepositoryInterface interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Classes, error)
	Create(ctx context.Context, attributes map[string]interface{}) (entity.Classes, error)
	FindByConditions(ctx context.Context, conditions map[string]interface{}) ([]entity.Classes, error)
}

type classesRepository struct {
	DB *gorm.DB
}

func NewClassesRepository(db *gorm.DB) ClassesRepositoryInterface {
	return &classesRepository{
		DB: db,
	}
}

// en: TakeByConditions function to get a class by conditions
func (c *classesRepository) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) (entity.Classes, error) {
	var class entity.Classes
	cdb := c.DB.WithContext(ctx)
	err := cdb.Model(&class).Where(conditions).Take(&entity.Classes{}).Error
	return class, err
}

// en: Create function to create a new class with given attributes
func (c *classesRepository) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.Classes, error) {
	var class entity.Classes
	err := utils.MapToStruct(attributes, &class)
	if err != nil {
		return entity.Classes{}, err
	}

	cdb := c.DB.WithContext(ctx)
	err = cdb.Create(&class).Error
	return class, err
}

// en: FindByConditions function to find classes by conditions
func (c *classesRepository) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) ([]entity.Classes, error) {
	var classes []entity.Classes
	cdb := c.DB.WithContext(ctx)
	err := cdb.Model(&entity.Classes{}).Where(conditions).Find(&classes).Error
	return classes, err
}
