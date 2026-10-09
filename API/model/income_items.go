package model

import "mysql/model/base"

type IncomeItem struct {
	base.ModelBase
	IncomeID         uint64  `json:"income_id" gorm:"type:bigint unsigned;not null;index:idx_income_items_income_id"`
	IncomeCategoryID *int    `json:"income_category_id" gorm:"index:idx_income_items_category_id"`
	Description      *string `json:"description" gorm:"type:text"`
	Quantity         float64 `json:"quantity" gorm:"type:decimal(8,2);not null;default:1.00"`
	UnitPrice        float64 `json:"unit_price" gorm:"type:decimal(12,2);not null;default:0.00"`
	SortOrder        int     `json:"sort_order" gorm:"not null;default:1"`
}

func (IncomeItem) TableName() string {
	return "income_items"
}
