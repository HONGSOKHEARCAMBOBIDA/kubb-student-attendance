package request

type CreateScoreRequest struct {
	ClassID      int64 `json:"class_id" binding:"required"`
	MajorID      int64 `json:"major_id" binding:"required"`
	GenerationID int64 `json:"generation_id" binding:"required"`
	ProgrammeID  int64 `json:"programme_id" binding:"required"`
	SubjectID    int64 `json:"subject_id" binding:"required"`
	Year         uint8 `json:"year" binding:"required"`
	Semester     uint8 `json:"semester" binding:"required"`

	Students []CreateScoreStudentRequest `json:"students" binding:"required,min=1"`
}

type CreateScoreStudentRequest struct {
	UserID  int64                      `json:"user_id" binding:"required"`
	Details []CreateScoreDetailRequest `json:"details" binding:"required,min=1"`
}

type CreateScoreDetailRequest struct {
	GradeComponentID uint64  `json:"grade_component_id" binding:"required"`
	Score            float64 `json:"score" binding:"gte=0"`
}

type ImportScoreExcelRequest struct {
	ClassID      int64 `form:"class_id" binding:"required"`
	SubjectID    int64 `form:"subject_id" binding:"required"`
	MajorID      int64 `form:"major_id" binding:"required"`
	GenerationID int64 `form:"generation_id" binding:"required"`
	ProgrammeID  int64 `form:"programme_id" binding:"required"`
	Year         uint8 `form:"year" binding:"required"`
	Semester     uint8 `form:"semester" binding:"required"`
}

type ImportScoreResult struct {
	Imported int      `json:"imported"`
	Skipped  int      `json:"skipped"`
	Errors   []string `json:"errors"`
}
