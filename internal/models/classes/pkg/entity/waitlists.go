package entity

type Waitlist struct {
	ID       int64  `gorm:"column:id;primaryKey;autoIncrement"`
	ClassID  int64  `gorm:"column:class_id;not null"`
	MemberID int64  `gorm:"column:member_id;not null"`
	JoinedAt string `gorm:"column:joined_at"`
	Notified bool   `gorm:"column:notified"`
}
