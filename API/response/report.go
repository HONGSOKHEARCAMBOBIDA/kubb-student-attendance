package response

type AcademicConfirmation struct {
	Namekh      string `json:"name_kh" gorm:"column:name_kh"`
	NameEn      string `json:"name_en"`
	Gender      int    `json:"gender"`
	GenderKh    string `json:"gender_kh"`
	GenderEn    string `json:"gender_en"`
	Nationalkh  string `json:"national_kh" gorm:"column:national_kh"`
	NationEn    string `json:"national_en" gorm:"column:national_en"`
	Codekh      string `json:"code_kh"`
	CodeEN      string `json:"code_en"`
	Dobkh       string `json:"dob_kh"`
	DobEn       string `json:"dob_en"`
	Villagekh   string `json:"village_kh" gorm:"column:village_kh"`
	VillageEn   string `json:"village_en"`
	Communcekh  string `json:"communce_kh" gorm:"column:communce_kh"`
	CommunceEn  string `json:"communce_en"`
	Districtkh  string `json:"district_kh" gorm:"column:district_kh"`
	DistrictEn  string `json:"district_en"`
	Provicekh   string `json:"province_kh" gorm:"column:province_kh"`
	ProviceEn   string `json:"province_en" gorm:"column:province_en"`
	IsCapital   bool   `json:"is_capital" gorm:"column:is_capital"`
	ProgrammeKh string `json:"programme_kh"`
	ProgrammeEn string `json:"programme_en"`
	Yearkh      string `json:"year_kh"`
	YearEn      string `json:"year_en"`
	Semesterkh  string `json:"semester_kh"`
	SemesterEn  string `json:"semester_en"`
	Facultykh   string `json:"faculty_kh" gorm:"column:faculty_kh"`
	FacultyEn   string `json:"faculty_en"`
	Majorkh     string `json:"major_kh" gorm:"column:major_kh"`
	MajorEn     string `json:"major_en"`
	StartKh     string `json:"start_kh"`
	StartEn     string `json:"start_en"`
	Endkh       string `json:"end_kh"`
	EndEn       string `json:"end_en"`
}

type TemporaryCertificate struct {
	ProgrammeKh string  `json:"programme_kh"`
	Namekh      string  `json:"name_kh" gorm:"column:name_kh"`
	NameEn      string  `json:"name_en"`
	Gender      int     `json:"gender"`
	GenderKh    string  `json:"gender_kh"`
	GenderEN    string  `json:"gender_en"`
	Nationalkh  string  `json:"national_kh" gorm:"column:national_kh"`
	Dobkh       string  `json:"dob_kh"`
	DobEn       string  `json:"dob_en" gorm:"column:dob_en"`
	FatherName  string  `gorm:"type:varchar(191);default:null" json:"father_name"`
	MotherName  *string `gorm:"type:varchar(191);default:null" json:"mother_name"`
	Majorkh     string  `json:"major_kh" gorm:"column:major_kh"`
	ExamOutDate string  `gorm:"type:tinyint;default:null" json:"exam_out_date"`
}
