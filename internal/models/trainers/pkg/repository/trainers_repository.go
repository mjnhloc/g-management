package repository

import (
	"context"

	"g-management/internal/models/trainers/pkg/entity"
	"g-management/pkg/shared/utils"

	"gorm.io/gorm"
)

type TrainersRepositoryInterface interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Trainers, error)
	FindByConditions(ctx context.Context, conditions map[string]interface{}) ([]entity.Trainers, error)
	CreateWithTransaction(tx *gorm.DB, attributes map[string]interface{}) (entity.Trainers, error)
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

// en: FindByConditions function to find trainers by conditions
func (t *trainersRepository) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) ([]entity.Trainers, error) {
	var trainers []entity.Trainers
	cdb := t.DB.WithContext(ctx)
	err := cdb.Where(conditions).Find(&trainers).Error
	return trainers, err
}

// en: CreateWithTransaction function to create a new trainer with given attributes within a transaction
func (t *trainersRepository) CreateWithTransaction(
	tx *gorm.DB,
	attributes map[string]interface{},
) (entity.Trainers, error) {
	var trainer entity.Trainers
	err := utils.MapToStruct(attributes, &trainer)
	if err != nil {
		return entity.Trainers{}, err
	}

	err = tx.Create(&trainer).Error
	return trainer, err
}
