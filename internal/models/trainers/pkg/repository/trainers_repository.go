package repository

import (
	"context"

	"g-management/internal/models/trainers/pkg/entity"
	"g-management/pkg/shared/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TrainersRepositoryInterface interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Trainers, error)
	FindByConditions(ctx context.Context, conditions map[string]interface{}) ([]entity.Trainers, error)
	CreateWithTransaction(tx *gorm.DB, attributes map[string]interface{}) (entity.Trainers, error)
	UpsertWithTransaction(tx *gorm.DB, attributes map[string]interface{}) (entity.Trainers, error)
	DeleteByConditions(ctx context.Context, conditions map[string]interface{}) error
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

// en: UpsertWithTransaction function to upsert a trainer with given attributes within a transaction
func (t *trainersRepository) UpsertWithTransaction(
	tx *gorm.DB,
	attributes map[string]interface{},
) (entity.Trainers, error) {
	var trainer entity.Trainers
	err := utils.MapToStruct(attributes, &trainer)
	if err != nil {
		return entity.Trainers{}, err
	}

	err = tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "id"},
		},
		DoUpdates: clause.Assignments(attributes),
	}).Create(&trainer).Error

	return trainer, err
}

// en: DeleteByConditions function to delete trainers by conditions
func (t *trainersRepository) DeleteByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) error {
	cdb := t.DB.WithContext(ctx)
	return cdb.Where(conditions).Delete(&entity.Trainers{}).Error
}
