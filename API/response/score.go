package response

import "mysql/model/base"

type ScoreResponse struct {
	base.ModelBase
	UserID         int    `json:"user_id"`
	NameKH         string `json:"name_kh" gorm:"column:name_kh"`
	NameEN         string `json:"name_en" gorm:"column:name_en"`
	Code           string `json:"code" gorm:"column:code"`
	Gender         int    `json:"gender"`
	ClassID        int    `json:"class_id" gorm:"column:class_id"`
	ClassName      string `json:"class_name"`
	ProgrammeID    int64  `gorm:"column:programme_id;not null"`
	ProgrammeName  string `gorm:"column:programme_name;not null"`
	GenerationID   int64  `gorm:"column:generation_id;not null"`
	GenerationName string `gorm:"column:generation_name;not null"`
	MajorID        int    `gorm:"column:major_id;not null"`
	MajorName      string `gorm:"column:major_name;not null"`
	SubjectID      int64  `gorm:"column:subject_id;not null"`
	SubjectName    string `gorm:"column:subject_name;not null"`

	Attendance float64 `json:"attendance"`
	Research   float64 `json:"research"`
	Midterm    float64 `json:"midterm"`
	Final      float64 `json:"final"`
}

type ScoreDetailResponse struct {
	ID                 int     `json:"id"`
	ScoreID            int     `json:"score_id"`
	GradeComponentID   int     `json:"grade_component_id"`
	GradeComponentName string  `json:"grade_component_Name"`
	Score              float64 `json:"score" gorm:"column:score"`
}
