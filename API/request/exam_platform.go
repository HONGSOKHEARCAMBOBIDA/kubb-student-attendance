package request

type ExamPlatform struct {
	ClassID            int                  `gorm:"not null" json:"class_id"`
	StartUse           string               `gorm:"type:date;not null" json:"start_use"`
	ExpireDate         string               `gorm:"type:date;not null" json:"expire_date"`
	ShowResult         *string              `gorm:"type:date" json:"show_result"`
	CheckBy            int                  `json:"check_by"`
	ExamPlatformDetail []ExamPlatformDetail `json:"details"`
}

type ExamPlatformDetail struct {
	SubjectID   int64  `gorm:"not null" json:"subject_id"`
	ExamDate    string `gorm:"type:date;not null" json:"exam_date"`
	StartTime   string `gorm:"type:time;not null" json:"start_time"`
	EndTime     string `gorm:"type:time;not null" json:"end_time"`
	TeacherID   int    `json:"teacher_id"`
	DocumentURL string `gorm:"type:text" json:"document_url"`
}

type ExamPlatformRequestUpdate struct {
	StartUse   string  `gorm:"type:date;not null" json:"start_use"`
	ExpireDate string  `gorm:"type:date;not null" json:"expire_date"`
	ShowResult *string `gorm:"type:date" json:"show_result"`
	CheckBy    int     `json:"check_by"`
}
