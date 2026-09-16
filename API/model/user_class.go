package model

import "mysql/model/base"

type UserClass struct {
	base.ModelBase
	UserID   int64 `gorm:"column:user_id;not null" json:"user_id"`
	ClassID  int64 `gorm:"column:class_id;not null" json:"class_id"`
	IsActive bool  `gorm:"column:is_active;not null;default:true" json:"is_active"`
}

func (UserClass) TableName() string {
	return "user_class"
}
