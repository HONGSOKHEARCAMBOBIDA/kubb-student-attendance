package request

import "mysql/model"

type SchoolarshipRequestCreate struct {
	Code               string             `json:"code" gorm:"column:code"`
	Name               string             `json:"name" gorm:"column:name"`
	DiscountType       model.DiscountType `gorm:"type:enum('percentage','amount');not null" json:"discount_type"`
	DiscountPercentage *float64           `json:"discount_percentage" gorm:"column:discount_percentage"`
	DiscountAmount     *float64           `json:"discount_amount" gorm:"column:discount_amount"`
	Description        *string            `json:"description" gorm:"column:description"`
}

type SchoolarshipRequestUpdate struct {
	Code               *string             `json:"code" gorm:"column:code"`
	Name               *string             `json:"name" gorm:"column:name"`
	DiscountType       *model.DiscountType `gorm:"type:enum('percentage','amount');not null" json:"discount_type"`
	DiscountPercentage *float64            `json:"discount_percentage" gorm:"column:discount_percentage"`
	DiscountAmount     *float64            `json:"discount_amount" gorm:"column:discount_amount"`
	Description        *string             `json:"description" gorm:"column:description"`
	Active             *bool               `json:"active" gorm:"column:active"`
}
