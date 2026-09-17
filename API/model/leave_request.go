package model

import (
	"mysql/model/base"
	"time"
)

const (
	LeaveStatusPending = "PENDING"
	LeaveStatusApprove = "APPROVE"
)

type LeaveRequest struct {
	base.ModelBase
	UserID         int    `gorm:"not null;index:idx_leave_request_user_id" json:"user_id"`
	ClassID        int    `gorm:"not null;index:idx_leave_request_class_id" json:"class_id"`
	StartDate      string `gorm:"type:date;not null;index:idx_leave_request_date" json:"start_date"`
	EndDate        string `gorm:"type:date;not null;index:idx_leave_request_date" json:"end_date"`
	BackToWorkDate string `gorm:"type:date;not null" json:"back_to_work_date"`

	TotalDay     float64 `gorm:"type:decimal(5,2);not null;default:0.00" json:"total_day"`
	DeductTypeID int     `gorm:"column:deduct_type_id" json:"deduct_type_id"`
	Reason       string  `gorm:"type:text" json:"reason"`

	Status string `gorm:"type:enum('PENDING','APPROVE');index:idx_leave_request_status" json:"status"`

	ApproveBy       *int            `gorm:"index:idx_leave_request_approve_by" json:"approve_by"`
	ApprovedAt      *time.Time      `json:"approved_at"`
	LeaveDeductType LeaveDeductType `gorm:"foreignKey:deduct_type_id"`
}

func (LeaveRequest) TableName() string {
	return "leave_request"
}
