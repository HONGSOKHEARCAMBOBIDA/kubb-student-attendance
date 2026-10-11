package model

import (
	"mysql/model/base"
)

type AcademicStream string

const (
	AcademicStreamSCIENCE        AcademicStream = "SCIENCE"
	AcademicStreamSOCIAL_SCIENCE AcademicStream = "SOCIAL_SCIENCE"
)

type StudentProfile struct {
	base.ModelBase
	UserID            int            `gorm:"default:null" json:"user_id"`
	StudentCategoryID int            `gorm:"default:null;index:idx_student_profile_student_category_id" json:"student_category_id"`
	DateOfBirth       string         `gorm:"type:date;default:null" json:"date_of_birth"`
	Nationality       *string        `gorm:"type:varchar(191);default:null" json:"nationality"`
	Phone             *string        `gorm:"type:varchar(191);default:null" json:"phone"`
	VillageID         *uint          `gorm:"default:null" json:"village_id"`
	Occupation        *string        `gorm:"type:varchar(190);default:null" json:"occupation"`
	AcademicStream    AcademicStream `gorm:"type:enum('SCIENCE','SOCIAL_SCIENCE');default:null" json:"academic_stream"`
	TelegramUsername  *string        `gorm:"type:varchar(191);default:null" json:"telegram_username"`
	ExamIn            bool           `gorm:"type:tinyint;default:null" json:"exam_in"`
	ExamOut           bool           `gorm:"type:tinyint;default:null" json:"exam_out"`
}

func (StudentProfile) TableName() string {
	return "student_profile"
}
