package model

import "mysql/model/base"

type MajorSubject struct {
	base.ModelBase

	MajorID      int   `gorm:"column:major_id;not null;index:idx_major_subject_major" json:"major_id"`
	GenerationID int   `gorm:"column:generation_id" json:"generation_id"`
	ProgrammeID  int   `gorm:"column:programme_id" json:"programme_id"`
	SubjectID    int   `gorm:"column:subject_id;not null;index:idx_major_subject_subject" json:"subject_id"`
	Year         uint8 `gorm:"column:year;not null" json:"year"`
	Semester     uint8 `gorm:"column:semester;not null" json:"semester"`
	IsActive     bool  `gorm:"column:is_active;not null;default:1" json:"is_active"`
}

func (MajorSubject) TableName() string {
	return "major_subject"
}
