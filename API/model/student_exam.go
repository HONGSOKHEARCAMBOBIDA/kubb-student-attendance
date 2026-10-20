package model

import "mysql/model/base"

type StudentExam struct {
	base.ModelBase
	ExamPlatformDetailID int64 `gorm:"not null;uniqueIndex:unique_exam_table;uniqueIndex:unique_exam_student" json:"exam_platform_detail_id"`
	ClassID              int64 `gorm:"not null" json:"class_id"`
	StudentID            int64 `gorm:"not null;uniqueIndex:unique_exam_student" json:"student_id"`
	TableNumber          int   `gorm:"uniqueIndex:unique_exam_table" json:"table_number"`
}

func (StudentExam) TableName() string {
	return "student_exam"
}
