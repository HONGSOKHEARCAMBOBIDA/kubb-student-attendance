package model

import "mysql/model/base"

type MajorPrice struct {
	base.ModelBase
	MajorID      int64   `json:"major_id" gorm:"not null;index"`
	GenerationID int64   `json:"generation_id" gorm:"not null;index"`
	ProgrammeID  int64   `json:"programme_id" gorm:"not null;index"`
	Year         int     `json:"year" gorm:"not null"`
	MonthlyFee   float64 `json:"monthly_fee" gorm:"type:decimal(12,2);not null;default:0.00"`
	QuarterFee   float64 `json:"quarter_fee" gorm:"type:decimal(12,2);not null;default:0.00"`
	SemesterFee  float64 `json:"semester_fee" gorm:"type:decimal(12,2);not null;default:0.00"`
	YearFee      float64 `json:"year_fee" gorm:"type:decimal(12,2);not null;default:0.00"`
	IsActive     bool    `json:"is_active" gorm:"not null;default:true"`
}

func (MajorPrice) TableName() string {
	return "major_price"
}
