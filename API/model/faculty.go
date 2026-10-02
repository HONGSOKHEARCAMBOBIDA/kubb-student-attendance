package model

type Faculty struct {
	ID          uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        *string `gorm:"type:varchar(191);uniqueIndex:faculties_code_unique" json:"code"`
	Name        string  `gorm:"type:mediumtext;not null" json:"name"`
	Description *string `gorm:"type:longtext" json:"description"`
	Active      int8    `gorm:"type:tinyint;not null;default:1" json:"active"`
}

func (Faculty) TableName() string {
	return "faculties"
}
