package model

import (
	"mysql/model/base"
)

type FeeTransaction struct {
	base.ModelBase
	FeeID         int     `gorm:"not null;index" json:"fee_id"`
	InstallmentID *int    `gorm:"index" json:"installment_id"`
	Code          string  `gorm:"type:varchar(32);uniqueIndex" json:"code"`
	Date          string  `gorm:"not null;default:CURRENT_TIMESTAMP;index" json:"date"`
	DueDate       *string `json:"due_date"`
	Amount        float64 `gorm:"type:decimal(12,2);not null;default:0.00" json:"amount"`
	Discount      float64 `gorm:"type:decimal(12,2);not null;default:0.00" json:"discount"`
	Tax           float64 `gorm:"type:decimal(12,2);not null;default:0.00" json:"tax"`
	Total         float64 `gorm:"type:decimal(12,2);not null;default:0.00" json:"total"`
	Reference     *string `gorm:"type:varchar(191)" json:"reference"`
	Method        *string `gorm:"type:varchar(50)" json:"method"`
	Message       *string `gorm:"type:longtext" json:"message"`
	Description   *string `gorm:"type:longtext" json:"description"`
	Active        bool    `gorm:"not null;default:true" json:"active"`
}

func (FeeTransaction) TableName() string {
	return "fee_transactions"
}
