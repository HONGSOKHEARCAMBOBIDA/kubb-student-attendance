package response

import (
	"mysql/model"
	"mysql/model/base"
)

type UserClass struct {
	base.ModelBase
	UserID          int             `json:"user_id"`
	NameKH          string          `json:"name_kh" gorm:"column:name_kh"`
	NameEN          string          `json:"name_en" gorm:"column:name_en"`
	Gender          int             `json:"gender" gorm:"column:gender"`
	Code            string          `json:"code" gorm:"column:code"`
	ClassID         int             `json:"class_id"`
	ClassName       string          `json:"class_name"`
	Type            model.ClassType `json:"type"`
	Active          bool            `json:"active"`
	MajorID         *int64          `gorm:"column:major_id" json:"major_id"`
	MajorName       string          `json:"major_name"`
	ShiftID         *int64          `gorm:"column:shift_id" json:"shift_id"`
	ShiftName       string          `json:"shift_name"`
	GenerationID    *int64          `gorm:"column:generation_id" json:"generation_id"`
	GenerationName  string          `json:"generation_name"`
	GenerationStart string          `json:"generation_start"`
	GenerationEnd   string          `json:"generation_end"`
	Year            *int8           `gorm:"column:year" json:"year"`
	// Semester           *int8                `gorm:"column:semester" json:"semester"`
	Group              *int8                `gorm:"column:group_name" json:"group_name"`
	Term               *int8                `gorm:"column:term" json:"term"`
	ProgrammeID        *int64               `gorm:"column:programme_id" json:"programme_id"`
	ProgrammeName      string               `json:"programme_name"`
	FeeID              uint64               `json:"fee_id" gorm:"column:fee_id"`
	Schoolarship       string               `json:"schoolarship"`
	Feeschedule        string               `json:"fee_schedule" gorm:"column:fee_schedule"`
	Amount             float64              `json:"amount"`
	Discount           float64              `json:"discount"`
	Total              float64              `json:"total"`
	PaidAmount         float64              `gorm:"column:paid_amount" json:"paid_amount"`
	InstallmentRespone []InstallmentRespone `json:"installments" gorm:"-"`
}

type InstallmentRespone struct {
	base.ModelBase
	FeeID               uint64  `gorm:"not null;index;uniqueIndex:uk_installment_sequence" json:"fee_id"`
	SequenceNo          int     `gorm:"not null;uniqueIndex:uk_installment_sequence" json:"sequence_no"`
	DueDate             string  `gorm:"type:date;not null;index" json:"due_date"`
	Amount              float64 `gorm:"type:decimal(12,2);not null;default:0.00" json:"amount"`
	Status              string  `gorm:"type:enum('pending','invoiced','partial','paid','overdue','cancelled');not null;default:'pending';index" json:"status"`
	FeetransactionID    int     `gorm:"column:fee_transacntion_id" json:"fee_transacntion_id"`
	FeetransactionTotal float64 `gorm:"column:fee_transaction_total" json:"fee_transaction_total"`
}

type PrintInvoiceResponse struct {
	NameKH          string  `json:"name_kh" gorm:"column:name_kh"`
	NameEN          string  `json:"name_en" gorm:"column:name_en"`
	Gender          int     `json:"gender" gorm:"column:gender"`
	Code            string  `json:"code" gorm:"column:code"`
	ProgrammeName   string  `json:"programme_name"`
	GenerationName  string  `json:"generation_name"`
	MajorName       string  `json:"major_name"`
	Term            *int8   `gorm:"column:term" json:"term"`
	Year            *int8   `gorm:"column:year" json:"year"`
	TransactionCode string  `json:"transacntion_code" gorm:"column:transacntion_code"`
	Date            string  `gorm:"not null;default:CURRENT_TIMESTAMP;index" json:"date"`
	DueDate         *string `json:"due_date"`
	Amount          float64 `gorm:"type:decimal(12,2);not null;default:0.00" json:"amount"`
	Discount        float64 `gorm:"type:decimal(12,2);not null;default:0.00" json:"discount"`
	Tax             float64 `gorm:"type:decimal(12,2);not null;default:0.00" json:"tax"`
	Total           float64 `gorm:"type:decimal(12,2);not null;default:0.00" json:"total"`
	Reference       *string `gorm:"type:varchar(191)" json:"reference"`
	Method          *string `gorm:"type:varchar(50)" json:"method"`
	Message         *string `gorm:"type:longtext" json:"message"`
	Description     *string `gorm:"type:longtext" json:"description"`
	SequenceNo      int     `gorm:"not null;uniqueIndex:uk_installment_sequence" json:"sequence_no"`
}
