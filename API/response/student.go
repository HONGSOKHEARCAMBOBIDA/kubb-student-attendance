package response

import (
	"mysql/model"
	"mysql/model/base"
)

type StudentStudyHistory struct {
	base.ModelBase
	NameKH   string `json:"name_kh" gorm:"column:name_kh"`
	NameEN   string `json:"name_en" gorm:"column:name_en"`
	Gender   int    `json:"gender" gorm:"column:gender"`
	Code     string `json:"code" gorm:"column:code"`
	RoleName string `json:"role_name"`

	StudentCategoryName string               `gorm:"default:null;index:idx_student_profile_student_category_id" json:"student_category_name"`
	DateOfBirth         string               `gorm:"type:date;default:null" json:"date_of_birth"`
	Nationality         *string              `gorm:"type:varchar(191);default:null" json:"nationality"`
	Phone               *string              `gorm:"type:varchar(191);default:null" json:"phone"`
	Province            string               `json:"province"`
	District            string               `json:"district"`
	Communce            string               `json:"communce"`
	Village             string               `json:"village"`
	Occupation          *string              `gorm:"type:varchar(190);default:null" json:"occupation"`
	AcademicStream      model.AcademicStream `gorm:"type:enum('SCIENCE','SOCIAL_SCIENCE');default:null" json:"academic_stream"`
	TelegramUsername    *string              `gorm:"type:varchar(191);default:null" json:"telegram_username"`
	ExamIn              bool                 `gorm:"type:tinyint;default:null" json:"exam_in"`
	ExamOut             bool                 `gorm:"type:tinyint;default:null" json:"exam_out"`
	UserClassRes        []UserClassRes       `json:"user_class_res" gorm:"-"`
}

type UserClassRes struct {
	UserClassID     int             `json:"user_class_id" gorm:"column:user_class_id"`
	UserID          int             `json:"user_id"`
	ClassName       string          `json:"class_name"`
	Type            model.ClassType `json:"type"`
	Active          bool            `json:"active"`
	MajorID         *int64          `gorm:"column:major_id" json:"major_id"`
	MajorName       string          `json:"major_name"`
	ShiftID         *int64          `gorm:"column:shift_id" json:"shift_id"`
	ShiftName       string          `json:"shift_name"`
	GenerationID    *int64          `gorm:"column:generation_id" json:"generation_id"`
	GenerationName  string          `json:"generation_name"`
	GenerationStart string          `json:"generation_start"`
	GenerationEnd   string          `json:"generation_end"`
	Year            *int8           `gorm:"column:year" json:"year"`
	Semester        int             `json:"semester"`
	Group           *int8           `gorm:"column:group_name" json:"group_name"`
	Term            *int8           `gorm:"column:term" json:"term"`
	ProgrammeID     *int64          `gorm:"column:programme_id" json:"programme_id"`
	ProgrammeName   string          `json:"programme_name"`
	SubjectRes      []SubjectRes    `json:"subject_res" gorm:"-"`
}

type SubjectRes struct {
	UserClassID   int    `json:"user_class_id"`
	SubjectID     int64  `json:"subject_id"`
	SubjectCode   string `json:"subject_code"`
	SubjectNameKh string `json:"subject_name_kh"`
	SubjectNameEn string `json:"subject_name_en"`
	CreditHour    int    `json:"credit_hour"`
}
