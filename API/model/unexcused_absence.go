package model

import (
	"mysql/model/base"
	"time"
)

type UnexcusedAbsence struct {
	base.ModelBase
	UserID    int       `gorm:"not null;index" json:"user_id"`
	ClassID   int       `gorm:"not null;index" json:"class_id"`
	CheckDate time.Time `gorm:"type:date;not null" json:"check_date"`
}

func (UnexcusedAbsence) TableName() string {
	return "unexcused_absence"
}
