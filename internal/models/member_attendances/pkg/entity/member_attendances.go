package entity

import (
	"time"

	"g-management/pkg/shared/gorm/model"
)

type MemberAttendances struct {
	ID         int       `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id"`
	MemberID   int       `gorm:"column:member_id;type:bigint;not null" mapstructure:"member_id"`
	ClassID    int       `gorm:"column:class_id;type:bigint;not null" mapstructure:"class_id"`
	AttendedAt time.Time `gorm:"column:attended_at;type:timestamp;not null" mapstructure:"attended_at"`
	model.BaseModel
}
