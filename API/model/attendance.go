package model

import "mysql/model/base"

type Attendance struct {
	base.ModelBase
	UserID            int     `gorm:"column:user_id" json:"user_id"`
	ClassID           int     `gorm:"column:class_id" json:"class_id"`
	CheckDate         string  `gorm:"column:check_date;type:date" json:"check_date"`
	Status            string  `gorm:"column:status;size:50" json:"status"`
	Score             float64 `gorm:"column:score" json:"score"`
	LeaveRequestID    *int    `gorm:"column:leave_request_id" json:"leave_request_id"`
	VerifyBy          *int    `gorm:"column:verify_by" json:"verify_by"`
	UnexcuseAbsenceID *int    `gorm:"column:unexcused_absence_id" json:"unexcused_absence_id"`
}

func (Attendance) TableName() string {
	return "attendance"
}
