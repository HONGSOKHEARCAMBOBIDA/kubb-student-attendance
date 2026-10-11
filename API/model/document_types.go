package model

import "mysql/model/base"

type DocumentType struct {
	base.ModelBase
	Code   *string `gorm:"type:varchar(50);uniqueIndex:code;default:null" json:"code"`
	NameKH string  `gorm:"type:varchar(255);not null" json:"name_kh"`
	NameEN *string `gorm:"type:varchar(255);default:null" json:"name_en"`
}

func (DocumentType) TableName() string {
	return "document_types"
}
