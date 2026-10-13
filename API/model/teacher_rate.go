package model

import (
	"mysql/model/base"
)

type TeacherRate struct {
	base.ModelBase
	UserID         int64   `gorm:"not null;index:idx_teacher_rates_user" json:"user_id"`
	ClassID        *int64  `gorm:"index:idx_teacher_rates_class" json:"class_id"`
	ClassSyRate    float64 `gorm:"type:decimal(12,2);not null;default:0.00" json:"hourly_rate"`
	EffeccheduleID *int64  `gorm:"index:idx_teacher_rates_class_schedule" json:"class_schedule_id"`
	HourltiveDate  string  `gorm:"type:date;not null;index:idx_teacher_rates_effective_date" json:"effective_date"`
	EndDate        *string `gorm:"type:date" json:"end_date"`
	Active         bool    `gorm:"type:tinyint;not null;default:1" json:"active"`
	CreatedBy      *int64  `gorm:"index:idx_teacher_rates_created_by" json:"created_by"`
}

func (TeacherRate) TableName() string {
	return "teacher_rates"
}
