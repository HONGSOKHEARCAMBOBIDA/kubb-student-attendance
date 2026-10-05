package response

import "mysql/model/base"

type MajorPrice struct {
	base.ModelBase
	MajorID        int64   `json:"major_id" gorm:"not null;index"`
	MajorName      string  `json:"major_name" gorm:"column:major_name"`
	GenerationID   int64   `json:"generation_id" gorm:"not null;index"`
	GenerationName string  `json:"generation_name" gorm:"column:generation_name"`
	ProgrammeID    int64   `json:"programme_id" gorm:"not null;index"`
	ProgrammeName  string  `json:"programme_name" gorm:"column:programme_name"`
	Year           int     `json:"year" gorm:"not null"`
	MonthlyFee     float64 `json:"monthly_fee" gorm:"type:decimal(12,2);not null;default:0.00"`
	QuarterFee     float64 `json:"quarter_fee" gorm:"type:decimal(12,2);not null;default:0.00"`
	SemesterFee    float64 `json:"semester_fee" gorm:"type:decimal(12,2);not null;default:0.00"`
	YearFee        float64 `json:"year_fee" gorm:"type:decimal(12,2);not null;default:0.00"`
	IsActive       bool    `json:"is_active" gorm:"not null;default:true"`
}
