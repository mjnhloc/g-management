package repository

import (
	"context"
	"g-management/internal/models/classes/pkg/entity"

	"gorm.io/gorm"
)

type WaitlistsRepositoryInterface interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Waitlist, error)
	Create(ctx context.Context, waitlist *entity.Waitlist) error
	CountByConditions(ctx context.Context, conditions map[string]interface{}) (int64, error)
}

type waitlistsRepository struct {
	DB *gorm.DB
}

func NewWaitlistsRepository(db *gorm.DB) WaitlistsRepositoryInterface {
	return &waitlistsRepository{DB: db}
}

func (r *waitlistsRepository) TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Waitlist, error) {
	var waitlist entity.Waitlist
	db := r.DB.WithContext(ctx)
	err := db.Where(conditions).Take(&waitlist).Error
	return waitlist, err
}

func (r *waitlistsRepository) Create(ctx context.Context, waitlist *entity.Waitlist) error {
	db := r.DB.WithContext(ctx)
	return db.Create(waitlist).Error
}

func (r *waitlistsRepository) CountByConditions(ctx context.Context, conditions map[string]interface{}) (int64, error) {
	var count int64
	db := r.DB.WithContext(ctx)
	err := db.Model(&entity.Waitlist{}).Where(conditions).Count(&count).Error
	return count, err
}
