package request

type TeacherRate struct {
	UserID     int64   `gorm:"not null;index:idx_teacher_rates_user" json:"user_id"`
	HourlyRate float64 `gorm:"type:decimal(12,2);not null;default:0.00" json:"hourly_rate"`
}
