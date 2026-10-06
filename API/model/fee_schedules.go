package model

import "mysql/model/base"

type FeeScheduleInterval string

const (
	FeeScheduleIntervalMonthlyFee  FeeScheduleInterval = "monthly_fee"
	FeeScheduleIntervalQuarterFee  FeeScheduleInterval = "quarterly_fee"
	FeeScheduleIntervalSemesterFee FeeScheduleInterval = "semesterly_fee"
	FeeScheduleIntervalYearlyFee   FeeScheduleInterval = "yearly_fee"
)

type FeeSchedule struct {
	base.ModelBase
	FeeInterval      FeeScheduleInterval `gorm:"type:enum('monthly_fee','quarterly_fee','semesterly_fee','yearly_fee');not null;uniqueIndex" json:"fee_interval"`
	InstallmentCount int                 `gorm:"not null;default:1" json:"installment_count"`
	Description      *string             `gorm:"type:text" json:"description"`
	Active           bool                `gorm:"not null;default:true" json:"active"`
}

func (FeeSchedule) TableName() string {
	return "fee_schedules"
}
