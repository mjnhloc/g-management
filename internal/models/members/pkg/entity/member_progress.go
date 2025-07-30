package entity

type MemberProgress struct {
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement"`
	MemberID  int64  `gorm:"column:member_id;not null"`
	Goal      string `gorm:"column:goal"`
	Value     int    `gorm:"column:value"`
	UpdatedAt string `gorm:"column:updated_at"`
}
