package model

import (
	"mysql/model/base"
)

type ExamPlatformDetail struct {
	base.ModelBase
	ExamPlatformID int64  `gorm:"not null" json:"exam_platform_id"`
	SubjectID      int64  `gorm:"not null" json:"subject_id"`
	ExamDate       string `gorm:"type:date;not null" json:"exam_date"`
	StartTime      string `gorm:"type:time;not null" json:"start_time"`
	EndTime        string `gorm:"type:time;not null" json:"end_time"`
	TeacherID      int    `json:"teacher_id"`
	DocumentURL    string `gorm:"type:text" json:"document_url"`
}

func (ExamPlatformDetail) TableName() string {
	return "exam_platform_detail"
}
