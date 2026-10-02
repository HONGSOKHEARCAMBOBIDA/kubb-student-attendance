package response

type Student struct {
	ID             int64
	Code           string
	NameKH         string
	NameEN         string
	Gender         int
	RegistrationNo string
	DateOfBirth    string
	PlaceOfBirth   string
	Nationality    string
	Campus         string
}

type Header struct {
	FacultyEN   string `json:"faculty_name" gorm:"column:faculty_name"`
	MajorEN     string `json:"major_name" gorm:"column:major_name"`
	DegreeTitle string `json:"degree_title"`
	StartYear   int    `json:"start_year"`
	EndYear     int    `json:"end_year"`
}

type RawResult struct {
	ScoreID    int64
	Year       int
	Semester   int
	Code       string
	NameEN     string
	NameKH     string
	Credits    int
	PointsSum  float64
	PercentSum float64
}

type Transcript struct {
	Student Student
	Header  Header
	Results []RawResult
}
