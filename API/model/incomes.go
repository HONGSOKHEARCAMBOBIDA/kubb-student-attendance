package model

import (
	"mysql/model/base"
)

type IncomeStatus string

const (
	IncomeStatusUNPAID    IncomeStatus = "UNPAID"
	IncomeStatusPARTIAL   IncomeStatus = "PARTIAL"
	IncomeStatusPAID      IncomeStatus = "PAID"
	IncomeStatusCANCELLED IncomeStatus = "CANCELLED"
)

type Income struct {
	base.ModelBase
	IncomeCode  string       `json:"income_code" gorm:"type:varchar(50);not null;uniqueIndex:uk_incomes_code"`
	CustomerID  int          `json:"customer_id" gorm:"type:bigint unsigned"`
	IncomeDate  string       `json:"income_date" gorm:"type:date;not null"`
	DueDate     string       `json:"due_date" gorm:"type:date"`
	Subtotal    float64      `json:"subtotal" gorm:"type:decimal(12,2);not null;default:0.00"`
	Tax         float64      `json:"tax" gorm:"type:decimal(12,2);not null;default:0.00"`
	Discount    float64      `json:"discount" gorm:"type:decimal(12,2);not null;default:0.00"`
	Total       float64      `json:"total" gorm:"type:decimal(12,2);not null;default:0.00"`
	Paid        float64      `json:"paid" gorm:"type:decimal(12,2);not null;default:0.00"`
	Status      IncomeStatus `json:"status" gorm:"type:enum('UNPAID','PARTIAL','PAID','CANCELLED');not null;default:UNPAID;index:idx_incomes_status"`
	Description *string      `json:"description" gorm:"type:varchar(500)"`
}

func (Income) TableName() string {
	return "incomes"
}
