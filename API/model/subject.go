package model

import "mysql/model/base"

type Subject struct {
	base.ModelBase
	Code       string  `gorm:"column:code;size:50;not null;uniqueIndex:uq_subject_code" json:"code"`
	NameKh     *string `gorm:"column:name_kh;size:255" json:"name_kh"`
	NameEn     *string `gorm:"column:name_en;size:255" json:"name_en"`
	CreditHour *uint8  `gorm:"column:credit_hour" json:"credit_hour"`
	IsActive   bool    `gorm:"column:is_active;not null;default:1" json:"is_active"`
}

func (Subject) TableName() string {
	return "subject"
}
