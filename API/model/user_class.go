package model

import "mysql/model/base"

type UserClassStatus string

const (
	UserClassStatusSTUDY    UserClassStatus = "STUDY"
	UserClassStatusSUSPEND  UserClassStatus = "SUSPEND"
	UserClassStatusTRANSFER UserClassStatus = "TRANSFER"
	UserClassStatusDROPPED  UserClassStatus = "DROPPED"
)

type UserClass struct {
	base.ModelBase
	UserID   int64           `gorm:"column:user_id;not null" json:"user_id"`
	ClassID  int64           `gorm:"column:class_id;not null" json:"class_id"`
	IsActive bool            `gorm:"column:is_active;not null;default:true" json:"is_active"`
	Status   UserClassStatus `json:"status"`
	Class    Class
}

func (UserClass) TableName() string {
	return "user_class"
}
