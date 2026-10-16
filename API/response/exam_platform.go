package response

import "mysql/model/base"

type ExamPlatform struct {
	base.ModelBase
	ClassID            int                  `gorm:"not null" json:"class_id"`
	Type               string               `gorm:"type:enum('MIDTERM','FINAL');not null" json:"type"`
	StartUse           string               `gorm:"type:date;not null" json:"start_use"`
	ExpireDate         string               `gorm:"type:date;not null" json:"expire_date"`
	ShowResult         *string              `gorm:"type:date" json:"show_result"`
	CreateBy           int                  `json:"create_by"`
	CreateByName       string               `json:"create_by_name"`
	CheckBy            int                  `json:"check_by"`
	CheckByName        string               `json:"check_by_name"`
	ExamPlatformDetail []ExamPlatformDetail `json:"details" gorm:"-"`
}

type ExamPlatformDetail struct {
	base.ModelBase
	ExamPlatformID int64  `gorm:"not null" json:"exam_platform_id"`
	SubjectID      int64  `gorm:"not null" json:"subject_id"`
	SubjectCode    string `json:"subject_code"`
	SubjectName    string `json:"subject_name"`
	ExamDate       string `gorm:"type:date;not null" json:"exam_date"`
	StartTime      string `gorm:"type:time;not null" json:"start_time"`
	EndTime        string `gorm:"type:time;not null" json:"end_time"`
	TeacherID      int    `json:"teacher_id"`
	TeacherName    string `json:"teacher_name"`
	DocumentURL    string `gorm:"type:text" json:"document_url"`
}
