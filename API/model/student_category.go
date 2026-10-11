package model

type StudentCategory struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(250);not null" json:"name"`
}

func (StudentCategory) TableName() string {
	return "student_category"
}
