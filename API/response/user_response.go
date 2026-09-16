package response

import (
	"mysql/model/base"
)

type UserResponse struct {
	base.ModelBase
	NameKH string `json:"name_kh" gorm:"column:name_kh"`
	NameEN string `json:"name_en" gorm:"column:name_en"`
	Gender int    `json:"gender" gorm:"column:gender"`
	Code   string `json:"code" gorm:"column:code"`
}

type UserCount struct {
	Total int `json:"total"`
}

type UserApprove struct {
	ID       int    `json:"id"`
	UserName string `json:"user_name"`
}

type StudentWithClass struct {
	UserResponse
	ClassID int `gorm:"column:class_id"`
}
