package response

import (
	"mysql/model"
	"mysql/model/base"
)

type UserClass struct {
	base.ModelBase
	UserID          int             `json:"user_id"`
	NameKH          string          `json:"name_kh" gorm:"column:name_kh"`
	NameEN          string          `json:"name_en" gorm:"column:name_en"`
	Gender          int             `json:"gender" gorm:"column:gender"`
	Code            string          `json:"code" gorm:"column:code"`
	ClassID         int             `json:"class_id"`
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
	Semester        *int8           `gorm:"column:semester" json:"semester"`
	Group           *int8           `gorm:"column:group" json:"group"`
	Term            *int8           `gorm:"column:term" json:"term"`
	ProgrammeID     *int64          `gorm:"column:programme_id" json:"programme_id"`
	ProgrammeName   string          `json:"programme_name"`
}
