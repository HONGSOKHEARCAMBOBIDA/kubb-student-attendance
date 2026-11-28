package model

import "mysql/model/base"

type MeetingType string

const (
	MeetingTypeGOOGLE_MEET     MeetingType = "GOOGLE_MEET"
	MeetingTypeZOOM            MeetingType = "ZOOM"
	MeetingTypeMICROSOFT_TEAMS MeetingType = "MICROSOFT_TEAMS"
	MeetingTypeOTHER           MeetingType = "OTHER"
)

type ClassDetail struct {
	base.ModelBase
	ClassID     int          `gorm:"not null" json:"class_id"`
	SubjectID   int          `gorm:"not null" json:"subject_id"`
	Credit      float64      `gorm:"not null;default:0" json:"credit"`
	TeacherID   int          `gorm:"not null" json:"teacher_id"`
	Midterm     bool         `json:"midterm" gorm:"column:midterm"`
	Final       bool         `json:"final" gorm:"column:final"`
	MeetingType *MeetingType `json:"meeting_type" gorm:"column:meeting_type"`
	MeetingLink *string      `json:"meeting_link" gorm:"column:meeting_link"`
}

func (ClassDetail) TableName() string {
	return "class_detail"
}
