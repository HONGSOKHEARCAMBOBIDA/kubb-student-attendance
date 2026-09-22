package response

type Major struct {
	ID       int64  `json:"id"`
	Code     string `json:"code"`
	NameKh   string `json:"name_kh"`
	NameEn   string `json:"name_en"`
	IsActive bool   `json:"is_active"`
}

// MajorSubject is the "flattened" view of one major_subject row joined
// with subject, used to list which subjects belong to a major.
type MajorSubject struct {
	ID             int64  `json:"id"`
	MajorID        int64  `json:"major_id"`
	GenerationName string `json:"generation_name"`
	ProgrammeName  string `json:"programme_name"`
	SubjectID      int64  `json:"subject_id"`
	SubjectCode    string `json:"subject_code"`
	SubjectNameKh  string `json:"subject_name_kh"`
	SubjectNameEn  string `json:"subject_name_en"`
	CreditHour     int    `json:"credit_hour"`
	Year           uint8  `json:"year"`
	Semester       uint8  `json:"semester"`
	IsActive       bool   `json:"is_active"`
}
