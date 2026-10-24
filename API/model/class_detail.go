package model

import "mysql/model/base"

type ClassDetail struct {
	base.ModelBase
	ClassID   int     `gorm:"not null" json:"class_id"`
	SubjectID int     `gorm:"not null" json:"subject_id"`
	Credit    float64 `gorm:"not null;default:0" json:"credit"`
	TeacherID int     `gorm:"not null" json:"teacher_id"`
}

func (ClassDetail) TableName() string {
	return "class_detail"
}
