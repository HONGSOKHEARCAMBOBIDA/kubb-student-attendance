package response

import "mysql/model/base"

type Student struct {
	base.ModelBase
	Code           string  `json:"code" gorm:"column:code"`
	NameKH         string  `json:"name_kh" gorm:"column:name_kh"`
	NameEN         string  `json:"name_en" gorm:"column:name_en"`
	Gender         int     `json:"gender" gorm:"column:gender"`
	RegistrationNo *string `gorm:"column:registration_no;size:50" json:"registration_no"`
	DateOfBirth    string  `gorm:"column:date_of_birth;type:date" json:"date_of_birth"`
	PlaceOfBirth   string  `gorm:"column:place_of_birth;size:255" json:"place_of_birth"`
	Nationality    string  `gorm:"column:nationality;size:100" json:"nationality"`
	Campus         string  `gorm:"column:campus;size:100" json:"campus"`
}

type Header struct {
	FacultyEN   string `json:"faculty_name" gorm:"column:faculty_name"`
	MajorEN     string `json:"major_name" gorm:"column:major_name"`
	DegreeTitle string `json:"degree_title"`
	StartYear   int    `json:"start_year"`
	EndYear     int    `json:"end_year"`
}

type RawResult struct {
	ScoreID    int64   `json:"score_id"`
	Year       int     `json:"year"`
	Semester   int     `json:"semester"`
	Code       string  `json:"code"`
	NameEN     string  `json:"name_en"`
	NameKH     string  `json:"name_kh"`
	Credits    int     `json:"credits"`
	PointsSum  float64 `json:"point_sum"`
	PercentSum float64 `json:"percent_sum"`
}

type Transcript struct {
	Student Student
	Header  Header
	Results []RawResult
}
