package entity

import "g-management/pkg/shared/gorm/model"

type Payments struct {
	ID            int    `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id"`
	MemberID      int    `gorm:"column:member_id;type:bigint;not null" mapstructure:"member_id"`
	MembershipID  int    `gorm:"column:membership_id;type:bigint;not null" mapstructure:"membership_id"`
	Amount        int    `gorm:"column:amount;type:int(11) unsigned;not null" mapstructure:"amount"`
	PaymentDate   string `gorm:"column:payment_date;type:timestamp;not null" mapstructure:"payment_date"`
	PaymentMethod string `gorm:"column:payment_method;type:enum('cash', 'credit_card', 'bank_transfer');not null;default:'cash'" mapstructure:"payment_method"`
	Status        string `gorm:"column:status;type:enum('completed', 'failed', 'refundeded');not null;default:'completed'" mapstructure:"status"`
	model.BaseModel
}
