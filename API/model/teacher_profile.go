package model

import (
	"mysql/model/base"
)

type TeacherProfile struct {
	base.ModelBase
	UserID       int64   `gorm:"not null;uniqueIndex:uk_teacher_profile_user" json:"user_id"`
	DateOfBirth  *string `gorm:"type:date" json:"date_of_birth"`
	PlaceOfBirth *string `gorm:"type:varchar(191)" json:"place_of_birth"`
	Nationality  *string `gorm:"type:varchar(191)" json:"nationality"`
	Address      *string `gorm:"type:varchar(191)" json:"address"`
	Phone        *string `gorm:"type:varchar(191)" json:"phone"`
}

func (TeacherProfile) TableName() string {
	return "teacher_profile"
}
