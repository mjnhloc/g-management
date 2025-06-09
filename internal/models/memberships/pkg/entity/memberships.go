package entity

import "g-management/pkg/shared/gorm/model"

type Memberships struct {
	ID             int    `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id"`
	MemberID       int    `gorm:"column:member_id;type:bigint;not null" mapstructure:"member_id"`
	MembershipType string `gorm:"column:membership_type;type:enum('monthly', 'quarter', 'annual');not null" mapstructure:"membership_type"`
	StartDate      string `gorm:"column:start_date;type:date;not null" mapstructure:"start_date"`
	EndDate        string `gorm:"column:end_date;type:date;not null" mapstructure:"end_date"`
	model.BaseModel
}
