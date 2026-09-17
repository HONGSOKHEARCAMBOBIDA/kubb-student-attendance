package model

import "mysql/model/base"

const (
	UnexcusedAbsenceSession1 = "session1"
	UnexcusedAbsenceSession2 = "session2"
	UnexcusedAbsenceSession3 = "session3"
	UnexcusedAbsenceSession4 = "session4"
	UnexcusedAbsenceSession5 = "session5"
)

type UnexcusedAbsenceDetail struct {
	base.ModelBase
	AbsenceID int    `gorm:"not null;uniqueIndex:uk_unexcused_absence_details;index:idx_unexcused_absence_details_absence_id" json:"absence_id"`
	Type      string `gorm:"type:enum('session1','session2','session3','session4','session5');not null;uniqueIndex:uk_unexcused_absence_details;index:idx_unexcused_absence_details_type" json:"type"`
}

func (UnexcusedAbsenceDetail) TableName() string {
	return "unexcused_absence_details"
}
