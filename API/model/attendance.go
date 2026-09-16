package model

import "mysql/model/base"

type Attendance struct {
	base.ModelBase
	UserID    int    `gorm:"column:user_id" json:"user_id"`
	ClassID   int    `gorm:"column:class_id" json:"class_id"`
	CheckDate string `gorm:"column:check_date;type:date" json:"check_date"`
	Status    string `gorm:"column:status;size:50" json:"status"`
}

func (Attendance) TableName() string {
	return "attendance"
}
