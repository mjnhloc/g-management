package activity_log

type ActivityLog struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement"`
	UserID     int64  `gorm:"column:user_id"`
	UserRole   string `gorm:"column:user_role"`
	Action     string `gorm:"column:action"`
	TargetType string `gorm:"column:target_type"`
	TargetID   int64  `gorm:"column:target_id"`
	Details    string `gorm:"column:details"`
	CreatedAt  string `gorm:"column:created_at"`
}
