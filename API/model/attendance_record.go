package model

import (
	"mysql/model/base"
)

const (
	AttendanceSession1 = "session1"
	AttendanceSession2 = "session2"
	AttendanceSession3 = "session3"
	AttendanceSession4 = "session4"
	AttendanceSession5 = "session5"
	AttendanceSession6 = "session6"
)

const (
	StatusPresent    = "PR"
	StatusPermission = "P"
	StatusAbsence    = "A"
)

type AttendanceRecord struct {
	base.ModelBase
	AttendanceID    int     `gorm:"column:attendance_id" json:"attendance_id"`
	UserID          int     `gorm:"column:user_id" json:"user_id"`
	ClassID         int     `gorm:"column:class_id" json:"class_id"`
	ShiftID         int     `gorm:"column:shift_id" json:"shift_id"`
	CheckTime       *string `gorm:"column:check_time;type:time" json:"check_time"`
	Type            string  `gorm:"column:type;type:enum('session1','session2','session3','session4','session5','session6')" json:"type"`
	ClassScheduleID int     `json:"class_schedule_id" gorm:"column:class_schedule_id"`
	SubjectID       int     `json:"subject_id" gorm:"column:subject_id"`
	Inzone          bool    `gorm:"column:inzone;not null;default:0" json:"inzone"`
	Latitude        string  `gorm:"column:latitude;size:255" json:"latitude"`
	Longitude       string  `gorm:"column:longitude;size:255" json:"longitude"`
	Status          string  `gorm:"column:status;type:enum('PR','P','A')" json:"status"`
}

func (AttendanceRecord) TableName() string {
	return "attendance_record"
}
