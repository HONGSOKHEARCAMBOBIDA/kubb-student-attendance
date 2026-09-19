package request

type MajorRequestCreate struct {
	Code   string `json:"code" binding:"required"`
	NameKh string `json:"name_kh" binding:"required"`
	NameEn string `json:"name_en" binding:"required"`
}

type MajorRequestUpdate struct {
	Code   string `json:"code" binding:"required"`
	NameKh string `json:"name_kh" binding:"required"`
	NameEn string `json:"name_en" binding:"required"`
}

// MajorSubjectRequestCreate is used to attach a subject to a major
// for a given academic year/semester (row in major_subject).
type MajorSubjectRequestCreate struct {
	SubjectID int   `json:"subject_id" binding:"required"`
	Year      uint8 `json:"year" binding:"required"`
	Semester  uint8 `json:"semester" binding:"required"`
}

// MajorSubjectRequestUpdate lets you move a subject to a different
// year/semester within the same major without removing/re-adding it.
type MajorSubjectRequestUpdate struct {
	Year     uint8 `json:"year" binding:"required"`
	Semester uint8 `json:"semester" binding:"required"`
}
