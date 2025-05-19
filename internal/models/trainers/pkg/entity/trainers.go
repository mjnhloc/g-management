package entity

import (
	"time"

	"g-management/pkg/shared/gorm/model"
)

type Trainers struct {
	ID             int       `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id"`
	Name           string    `gorm:"column:name;type:varchar(20);not null" mapstructure:"name"`
	Email          *string   `gorm:"column:email;type:varchar(50)" mapstructure:"email"`
	Phone          string    `gorm:"column:phone;type:varchar(15);unique;not null" mapstructure:"phone"`
	Specialization *string   `gorm:"column:specialization;type:varchar(255)" mapstructure:"specialization"`
	HiredAt        time.Time `gorm:"column:hired_at;type:date;not null" mapstructure:"hired_at"`
	model.BaseModel
}
