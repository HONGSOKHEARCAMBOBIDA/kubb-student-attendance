package model

import "mysql/model/base"

type ClassSchedule struct {
	base.ModelBase
	ClassID       int64  `gorm:"not null;index:idx_class_schedule_class" json:"class_id"`
	SubjectID     int64  `gorm:"not null;index:idx_class_schedule_subject" json:"subject_id"`
	DayOfWeek     int8   `gorm:"not null" json:"day_of_week"`
	ScheduleDate  string `gorm:"column:schedule_date" json:"schedule_date"`
	TotalSession  int    `gorm:"column:total_session" json:"total_session"`
	IsActive      bool   `gorm:"not null;default:1" json:"is_active"`
	ClassDetailID int    `json:"class_detail_id" gorm:"column:class_detail_id"`
	TotalStudent  int    `json:"total_student" gorm:"column:total_student"`
	Subject       Subject
}

func (ClassSchedule) TableName() string {
	return "class_schedule"
}
