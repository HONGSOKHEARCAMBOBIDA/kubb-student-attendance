package model

import "mysql/model/base"

type MajorExamType string

const (
	MajorExamTypeENTRANCE MajorExamType = "ENTRANCE"
	MajorExamTypeEXIT     MajorExamType = "EXIT"
)

type MajorExam struct {
	base.ModelBase
	MajorID      int    `gorm:"column:major_id;not null;index:idx_major_subject_major" json:"major_id"`
	GenerationID int    `gorm:"column:generation_id" json:"generation_id"`
	ProgrammeID  int    `gorm:"column:programme_id" json:"programme_id"`
	SubjectID    int    `gorm:"column:subject_id;not null;index:idx_major_subject_subject" json:"subject_id"`
	ExamType     string `gorm:"type:enum('ENTRANCE','EXIT')" json:"exam_type"`
	IsActive     bool   `gorm:"column:is_active;not null;default:1" json:"is_active"`
}

func (MajorExam) TableName() string {
	return "major_exam"
}
