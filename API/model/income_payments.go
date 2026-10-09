package model

import (
	"mysql/model/base"
)

type IncomePayment struct {
	base.ModelBase
	IncomeID    uint64  `json:"income_id" gorm:"type:bigint unsigned;not null;index:idx_income_payments_income_id"`
	PaymentCode string  `json:"payment_code" gorm:"type:varchar(50);not null;uniqueIndex:uk_income_payment_code"`
	PaymentDate string  `json:"payment_date" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;index:idx_income_payments_payment_date"`
	Amount      float64 `json:"amount" gorm:"type:decimal(12,2);not null;default:0.00"`
	Method      string  `json:"method" gorm:"type:varchar(50);not null;default:CASH"`
	Reference   *string `json:"reference" gorm:"type:varchar(191)"`
	Description *string `json:"description" gorm:"type:varchar(500)"`
	CreatedBy   *uint64 `json:"created_by" gorm:"type:bigint unsigned"`
}

func (IncomePayment) TableName() string {
	return "income_payments"
}
