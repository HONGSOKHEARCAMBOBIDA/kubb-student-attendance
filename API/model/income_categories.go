package model

import "mysql/model/base"

type IncomeCategory struct {
	base.ModelBase
	Name        string `json:"name" gorm:"type:mediumtext;not null"`
	Description string `json:"description" gorm:"type:text"`
	Active      bool   `json:"active" gorm:"type:tinyint;not null;default:1"`
}

func (IncomeCategory) TableName() string {
	return "income_categories"
}
