package repository

import (
	"context"

	"g-management/internal/models/classes/pkg/entity"

	"gorm.io/gorm"
)

type ClassesRepositoryInterface interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Classes, error)
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
