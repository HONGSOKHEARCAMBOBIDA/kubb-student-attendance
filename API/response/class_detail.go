package response

import "mysql/model/base"

type ClassDetail struct {
	base.ModelBase
	ClassID     int     `gorm:"not null" json:"class_id"`
	SubjectID   int     `gorm:"not null" json:"subject_id"`
	SubjectName string  `json:"subject_name" gorm:"column:subject_name"`
	Credit      float64 `gorm:"not null;default:0" json:"credit"`
	Remain      float64 `json:"remaining_hours" gorm:"column:remaining_hours"`
	TeacherID   int     `gorm:"not null" json:"teacher_id"`
	TeacherName string  `json:"teacher_name" gorm:"column:teacher_name"`
}
