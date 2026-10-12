package response

import (
	"mysql/model"
	"mysql/model/base"
)

type StudentProfile struct {
	base.ModelBase
	UserID            int                  `gorm:"default:null" json:"user_id"`
	StudentCategoryID int                  `gorm:"default:null;index:idx_student_profile_student_category_id" json:"student_category_id"`
	DateOfBirth       string               `gorm:"type:date;default:null" json:"date_of_birth"`
	Nationality       *string              `gorm:"type:varchar(191);default:null" json:"nationality"`
	Phone             *string              `gorm:"type:varchar(191);default:null" json:"phone"`
	VillageID         int                  `gorm:"default:null" json:"village_id"`
	CommunceID        int                  `json:"commune_id" gorm:"column:commune_id"`
	DistrictID        int                  `json:"district_id" gorm:"column:district_id"`
	ProvinceID        int                  `json:"province_id"`
	Occupation        *string              `gorm:"type:varchar(190);default:null" json:"occupation"`
	AcademicStream    model.AcademicStream `gorm:"type:enum('SCIENCE','SOCIAL_SCIENCE');default:null" json:"academic_stream"`
	TelegramUsername  *string              `gorm:"type:varchar(191);default:null" json:"telegram_username"`
	ExamIn            bool                 `gorm:"type:tinyint;default:null" json:"exam_in"`
	ExamOut           bool                 `gorm:"type:tinyint;default:null" json:"exam_out"`
}
