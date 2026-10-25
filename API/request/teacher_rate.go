package request

type TeacherRate struct {
	UserID       int64   `gorm:"not null;index:idx_teacher_rates_user" json:"user_id"`
	HourlyRate   float64 `gorm:"type:decimal(12,2);not null;default:0.00" json:"hourly_rate"`
	SubjectID    int64   `gorm:"not null;index:idx_class_schedule_subject" json:"subject_id"`
	DayOfWeek    int8    `gorm:"not null" json:"day_of_week"`
	ScheduleDate string  `gorm:"column:schedule_date" json:"schedule_date"`
	TotalSession int     `gorm:"column:total_session" json:"total_session"`
	TotalStudent int     `json:"total_student" gorm:"column:total_student"`
}
