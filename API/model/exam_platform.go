package model

import (
	"mysql/model/base"
)

type ExamPlatformType string

const (
	ExamPlatformTypeMIDTERM ExamPlatformType = "MIDTERM"
	ExamPlatformTypeFINAL   ExamPlatformType = "FINAL"
)

type ExamPlatform struct {
	base.ModelBase
	ClassID    int     `gorm:"not null" json:"class_id"`
	Type       string  `gorm:"type:enum('MIDTERM','FINAL');not null" json:"type"`
	StartUse   string  `gorm:"type:date;not null" json:"start_use"`
	ExpireDate string  `gorm:"type:date;not null" json:"expire_date"`
	ShowResult *string `gorm:"type:date" json:"show_result"`
	CreateBy   int     `json:"create_by"`
	CheckBy    int     `json:"check_by"`
}

func (ExamPlatform) TableName() string {
	return "exam_platform"
}
