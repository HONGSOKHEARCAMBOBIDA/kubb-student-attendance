package model

import (
	"mysql/model/base"
)

type Fee struct {
	base.ModelBase
	UserClassID   int     `gorm:"not null;index" json:"user_class_id"`
	UserID        int     `gorm:"not null;index" json:"user_id"`
	ClassID       int     `gorm:"not null;index" json:"class_id"`
	MajorID       int     `gorm:"not null" json:"major_id"`
	GenerationID  int     `gorm:"not null" json:"generation_id"`
	ProgrammeID   int     `gorm:"not null" json:"programme_id"`
	Year          int     `gorm:"type:tinyint;not null" json:"year"`
	Term          *int    `gorm:"type:tinyint" json:"term"`
	ScholarshipID *int    `gorm:"index" json:"scholarship_id"`
	FeeScheduleID int     `gorm:"column:fee_schedule_id" json:"fee_schedule_id"`
	Date          string  `gorm:"type:date;not null" json:"date"`
	Amount        float64 `gorm:"type:decimal(20,2);not null;default:0.00" json:"amount"`
	Discount      float64 `gorm:"type:decimal(20,2);not null;default:0.00" json:"discount"`
	Total         float64 `gorm:"type:decimal(20,2);not null;default:0.00" json:"total"`
	PaidAmount    float64 `gorm:"column:paid_amount" json:"paid_amount"`
	Active        bool    `gorm:"not null;default:true" json:"active"`
}

func (Fee) TableName() string {
	return "fees"
}
