package response

import "mysql/model/base"

type ClassResponse struct {
	base.ModelBase
	Name           *string `gorm:"column:name;type:varchar(255)" json:"name"`
	IsActive       bool    `gorm:"column:is_active;default:1" json:"is_active"`
	Latitude       *string `gorm:"column:latitude;type:varchar(255)" json:"latitude"`
	Longitude      *string `gorm:"column:longitude;type:varchar(255)" json:"longitude"`
	Radius         *string `gorm:"column:radius;type:varchar(255)" json:"radius"`
	BotToken       *string `gorm:"column:bot_token;type:varchar(255)" json:"bot_token"`
	GroupChatID    *string `gorm:"column:group_chatID;type:varchar(255)" json:"group_chatID"`
	CanScanOutsize bool    `gorm:"column:can_scan_outsize;default:0" json:"can_scan_outsize"`

	MajorID        *int64         `gorm:"column:major_id" json:"major_id"`
	MajorName      string         `json:"major_name"`
	ShiftID        *int64         `gorm:"column:shift_id" json:"shift_id"`
	ShiftName      string         `json:"shift_name"`
	GenerationID   *int64         `gorm:"column:generation_id" json:"generation_id"`
	GenerationName string         `json:"generation_name"`
	Year           *int8          `gorm:"column:year" json:"year"`
	Semester       *int8          `gorm:"column:semester" json:"semester"`
	Group          *int8          `gorm:"column:group" json:"group"`
	Term           *int8          `gorm:"column:term" json:"term"`
	ProgrammeID    *int64         `gorm:"column:programme_id" json:"programme_id"`
	ProgrammeName  string         `json:"programme_name"`
	UserResponse   []UserResponse `json:"students" gorm:"-"`
}
