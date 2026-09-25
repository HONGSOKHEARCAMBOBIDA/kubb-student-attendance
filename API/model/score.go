package model

import "mysql/model/base"

type Score struct {
	base.ModelBase
	ClassID      int64 `gorm:"column:class_id;not null"`
	UserID       int64 `gorm:"column:user_id;not null"`
	MajorID      int64 `gorm:"column:major_id;not null"`
	GenerationID int64 `gorm:"column:generation_id;not null"`
	ProgrammeID  int64 `gorm:"column:programme_id;not null"`
	SubjectID    int64 `gorm:"column:subject_id;not null"`
	Year         uint8 `gorm:"column:year;not null"`
	Semester     uint8 `gorm:"column:semester;not null"`
	IsActive     bool  `gorm:"column:is_active;not null;default:true"`
}

func (Score) TableName() string {
	return "score"
}
