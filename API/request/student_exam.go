package request

type StudentExam struct {
	ExamPlatformDetailID int64     `json:"exam_platform_detail_id" binding:"required"`
	Student              []Student `json:"student" binding:"required,min=1,dive"`
	Room                 string    `json:"room" gorm:"column:room"`
}

type Student struct {
	ClassID     int64 `json:"class_id" binding:"required"`
	StudentID   int64 `json:"student_id" binding:"required"`
	TableNumber *int  `json:"table_number"` // nil = unassigned
}
