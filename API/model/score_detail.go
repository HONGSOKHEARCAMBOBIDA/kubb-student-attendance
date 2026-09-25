package model

import (
	"mysql/model/base"
)

type ScoreDetail struct {
	base.ModelBase
	ScoreID          int64   `gorm:"column:score_id;not null"`
	GradeComponentID uint64  `gorm:"column:grade_component_id;not null"`
	Score            float64 `gorm:"column:score;type:decimal(5,2);not null;default:0"`
}

func (ScoreDetail) TableName() string {
	return "score_detail"
}
