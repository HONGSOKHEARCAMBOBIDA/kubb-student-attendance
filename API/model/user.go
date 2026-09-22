package model

import "mysql/model/base"

type User struct {
	base.ModelBase
	NameKH   string `json:"name_kh" gorm:"column:name_kh"`
	NameEN   string `json:"name_en" gorm:"column:name_en"`
	Gender   int    `json:"gender" gorm:"column:gender"`
	Code     string `json:"code" gorm:"column:code"`
	RoleID   int    `json:"role_id" gorm:"column:role_id"`
	Password string `json:"password" gorm:"column:password"`
	Role     Role
}

func (User) TableName() string {
	return "user"
}
