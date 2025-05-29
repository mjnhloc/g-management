package repository

import (
	"context"

	"g-management/internal/models/trainers/pkg/entity"

	"gorm.io/gorm"
)

type TrainersRepositoryInterface interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Trainers, error)
}

type trainersRepository struct {
	DB *gorm.DB
}

func NewTrainersRepository(db *gorm.DB) TrainersRepositoryInterface {
	return &trainersRepository{
		DB: db,
	}
}

// en: TakeByConditions function to get a trainer by conditions
func (t *trainersRepository) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) (entity.Trainers, error) {
	var trainer entity.Trainers
	cdb := t.DB.WithContext(ctx)
	err := cdb.Where(conditions).Take(&trainer).Error
	return trainer, err
}
