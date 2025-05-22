package entity

import (
	"time"

	"g-management/pkg/shared/gorm/model"
)

type Members struct {
	ID          int        `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id"`
	Name        string     `gorm:"column:name;type:varchar(20);not null" mapstructure:"name"`
	Email       *string    `gorm:"column:email;type:varchar(50)" mapstructure:"email"`
	Phone       string     `gorm:"column:phone;type:varchar(20);not null" mapstructure:"phone"`
	DateOfBirth *time.Time `gorm:"column:date_of_birth;type:date" mapstructure:"date_of_birth"`
	IsActive    bool       `gorm:"column:is_active;type:tinyint(1);not null;default:1" mapstructure:"is_active"`
	model.BaseModel
}
