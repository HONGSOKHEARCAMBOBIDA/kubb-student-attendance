package request

type RoleRequestUpdate struct {
	Name        *string `json:"name"`
	DisPlayName *string `json:"display_name" gorm:"column:display_name"`
	Level       int     `json:"level" gorm:"column:level"`
}

type RoleRequestCreate struct {
	Name        *string `json:"name"`
	DisPlayName *string `json:"display_name" gorm:"column:display_name"`
	Level       int     `json:"level" gorm:"column:level"`
}
