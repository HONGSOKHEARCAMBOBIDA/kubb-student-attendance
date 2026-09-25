package model

type GradeComponent struct {
	ID               uint64  `gorm:"column:id;primaryKey;autoIncrement"`
	Name             string  `gorm:"column:name;type:varchar(191);not null"`
	WeightPercentage float64 `gorm:"column:weight_percentage;type:decimal(5,2);not null;default:0"`
	Active           bool    `gorm:"column:active;not null;default:true"`
}

func (GradeComponent) TableName() string {
	return "grade_components"
}
