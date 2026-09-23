package response

import (
	"mysql/model"
	"mysql/model/base"
)

type AuthResponse struct {
	ID           int                `json:"id"`
	Name         string             `json:"name"`
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
	Permissions  []model.Permission `json:"permissions"`
}

type UserDataResponse struct {
	ID          int                `json:"id"`
	Name        string             `json:"name"`
	RoleID      int                `json:"role_id" gorm:"column:role_id"`
	Level       int                `json:"level"`
	ClassID     int                `json:"class_id" gorm:"column:class_id"`
	Permissions []model.Permission `json:"permissions" gorm:"-"`
}

type UserResponseNotStudent struct {
	base.ModelBase
	NameKH   string `json:"name_kh" gorm:"column:name_kh"`
	NameEN   string `json:"name_en" gorm:"column:name_en"`
	Gender   int    `json:"gender" gorm:"column:gender"`
	Code     string `json:"code" gorm:"column:code"`
	RoleID   int    `json:"role_id" gorm:"column:role_id"`
	RoleName string `json:"role_name" gorm:"column:role_name"`
}
