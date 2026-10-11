package model

import (
	"mysql/model/base"
)

type StudentFamily struct {
	base.ModelBase
	StudentID         int     `gorm:"not null;index:student_families_student_id_index" json:"student_id"`
	FatherName        *string `gorm:"type:varchar(191);default:null" json:"father_name"`
	FatherEnglishName *string `gorm:"type:varchar(191);default:null" json:"father_english_name"`
	FatherAge         *string `gorm:"type:varchar(191);default:null" json:"father_age"`
	FatherIsAlive     bool    `gorm:"type:tinyint;not null;default:1" json:"father_is_alive"`
	FatherPhoneNumber *string `gorm:"type:varchar(191);default:null" json:"father_phone_number"`
	FatherOccupation  *string `gorm:"type:varchar(191);default:null" json:"father_occupation"`
	FatherWorkplace   *string `gorm:"type:varchar(191);default:null" json:"father_workplace"`
	MotherName        *string `gorm:"type:varchar(191);default:null" json:"mother_name"`
	MotherEnglishName *string `gorm:"type:varchar(191);default:null" json:"mother_english_name"`
	MotherAge         *string `gorm:"type:varchar(191);default:null" json:"mother_age"`
	MotherIsAlive     bool    `gorm:"type:tinyint;not null;default:1" json:"mother_is_alive"`
	MotherPhoneNumber *string `gorm:"type:varchar(191);default:null" json:"mother_phone_number"`
	MotherOccupation  *string `gorm:"type:varchar(191);default:null" json:"mother_occupation"`
	MotherWorkplace   *string `gorm:"type:varchar(191);default:null" json:"mother_workplace"`
}

func (StudentFamily) TableName() string {
	return "student_families"
}
