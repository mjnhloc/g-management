package repository

import (
	"context"

	"g-management/internal/models/member-attendances/pkg/entity"

	"gorm.io/gorm"
)

type MemberAttendancesRepositoryInterface interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.MemberAttendances, error)
}

type memberAttendancesRepository struct {
	DB *gorm.DB
}

func NewMemberAttendancesRepository(db *gorm.DB) MemberAttendancesRepositoryInterface {
	return &memberAttendancesRepository{
		DB: db,
	}
}

// en: TakeByConditions function to get a member by conditions
func (m *memberAttendancesRepository) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) (entity.MemberAttendances, error) {
	var memberAttendance entity.MemberAttendances
	mdb := m.DB.WithContext(ctx)
	err := mdb.Model(&memberAttendance).Where(conditions).Take(&entity.MemberAttendances{}).Error
	return memberAttendance, err
}
