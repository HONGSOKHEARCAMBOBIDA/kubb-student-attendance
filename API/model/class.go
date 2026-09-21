package model

import "mysql/model/base"

type ClassType string

const (
	ClassTypeOnClass ClassType = "onclass"
	ClassTypeOnLine  ClassType = "online"
)

type Class struct {
	base.ModelBase
	Name           *string   `gorm:"column:name;type:varchar(255)" json:"name"`
	Type           ClassType `gorm:"column:type" json:"type"`
	IsActive       bool      `gorm:"column:is_active;default:1" json:"is_active"`
	Latitude       *string   `gorm:"column:latitude;type:varchar(255)" json:"latitude"`
	Longitude      *string   `gorm:"column:longitude;type:varchar(255)" json:"longitude"`
	Radius         string    `gorm:"column:radius;type:varchar(255)" json:"radius"`
	BotToken       *string   `gorm:"column:bot_token;type:varchar(255)" json:"bot_token"`
	GroupChatID    *string   `gorm:"column:group_chatID;type:varchar(255)" json:"group_chatID"`
	CanScanOutsize bool      `gorm:"column:can_scan_outsize;default:0" json:"can_scan_outsize"`

	MajorID      *int64 `gorm:"column:major_id" json:"major_id"`
	ShiftID      *int64 `gorm:"column:shift_id" json:"shift_id"`
	GenerationID *int64 `gorm:"column:generation_id" json:"generation_id"`
	Year         int    `gorm:"column:year" json:"year"`
	Semester     int    `gorm:"column:semester" json:"semester"`
	Group        int    `gorm:"column:group" json:"group"`
	Term         int    `gorm:"column:term" json:"term"`
	ProgrammeID  *int64 `gorm:"column:programme_id" json:"programme_id"`
}

func (Class) TableName() string {
	return "class"
}
