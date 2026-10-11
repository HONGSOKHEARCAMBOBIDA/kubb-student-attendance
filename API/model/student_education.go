package model

import (
	"mysql/model/base"
)

type StudentEducation struct {
	base.ModelBase
	StudentID       int     `gorm:"not null;index:student_education_student_id_index" json:"student_id"`
	Level           *string `gorm:"type:varchar(191);default:null" json:"level"`
	SchoolName      *string `gorm:"type:varchar(191);default:null" json:"school_name"`
	StartDate       *string `gorm:"type:date;default:null" json:"start_date"`
	EndDate         *string `gorm:"type:date;default:null" json:"end_date"`
	CertificateDate *string `gorm:"column:cerificate_date;type:date;default:null" json:"cerificate_date"`
	Score           *string `gorm:"type:varchar(191);default:null" json:"score"`
	GPA             *string `gorm:"type:varchar(191);default:null" json:"gpa"`
	Grade           *string `gorm:"type:varchar(191);default:null" json:"grade"`
}

func (StudentEducation) TableName() string {
	return "student_education"
}
