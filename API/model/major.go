package model

type Major struct {
	ID       int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	NameKh   string `gorm:"column:name_kh;size:255" json:"name_kh"`
	NameEn   string `gorm:"column:name_en;size:255" json:"name_en"`
	Code     string `gorm:"column:code;size:100;uniqueIndex:uq_major_code" json:"code"`
	IsActive bool   `gorm:"column:is_active;not null;default:true" json:"is_active"`
}

func (Major) TableName() string {
	return "major"
}
