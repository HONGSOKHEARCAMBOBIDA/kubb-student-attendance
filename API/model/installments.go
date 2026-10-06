package model

import (
	"mysql/model/base"
	"time"
)

type InstallmentStatus string

const (
	InstallmentStatusPending  InstallmentStatus = "pending"
	InstallmentStatusInvoiced InstallmentStatus = "invoiced"
	InstallmentStatusPartial  InstallmentStatus = "partial"
	InstallmentStatusPaid     InstallmentStatus = "paid"
	InstallmentStatusOverdue  InstallmentStatus = "overdue"
	InstallmentStatusCancell  InstallmentStatus = "cancelled"
)

type Installment struct {
	base.ModelBase
	FeeID      uint64    `gorm:"not null;index;uniqueIndex:uk_installment_sequence" json:"fee_id"`
	SequenceNo int       `gorm:"not null;uniqueIndex:uk_installment_sequence" json:"sequence_no"`
	DueDate    time.Time `gorm:"type:date;not null;index" json:"due_date"`
	Amount     float64   `gorm:"type:decimal(12,2);not null;default:0.00" json:"amount"`
	Status     string    `gorm:"type:enum('pending','invoiced','partial','paid','overdue','cancelled');not null;default:'pending';index" json:"status"`
}

func (Installment) TableName() string {
	return "installments"
}
