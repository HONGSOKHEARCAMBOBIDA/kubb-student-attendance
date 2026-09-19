package response

// ClassScheduleResponse is a class_schedule row joined with subject, so the
// table can show the subject code/name without a second round trip.
type ClassScheduleResponse struct {
	ID          int64  `json:"id"`
	ClassID     int64  `json:"class_id"`
	SubjectID   int64  `json:"subject_id"`
	SubjectCode string `json:"subject_code"`
	SubjectName string `json:"subject_name_kh"`
	DayOfWeek   int8   `json:"day_of_week"`
	IsActive    bool   `json:"is_active"`
}

// SubjectOption is the trimmed-down shape used to populate the "add
// schedule" subject select — only subjects valid for a given
// class/major/year/semester combination.
type SubjectOption struct {
	ID     int64  `json:"id"`
	Code   string `json:"code"`
	NameKh string `json:"name_kh"`
	NameEn string `json:"name_en"`
}
