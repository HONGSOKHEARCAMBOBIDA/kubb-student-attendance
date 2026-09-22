package request

type GenerationRequestCreate struct {
	NameKh    string `gorm:"column:name_kh;size:255" json:"name_kh"`
	NameEn    string `gorm:"column:name_en;size:255" json:"name_en"`
	Code      string `gorm:"column:code;size:100;uniqueIndex:uq_generation_code" json:"code"`
	StartYear int    `gorm:"column:start_year;type:year" json:"start_year"`
	EndYear   int    `gorm:"column:end_year;type:year" json:"end_year"`
}

type GenerationRequestUpdate struct {
	NameKh    string `gorm:"column:name_kh;size:255" json:"name_kh"`
	NameEn    string `gorm:"column:name_en;size:255" json:"name_en"`
	Code      string `gorm:"column:code;size:100;uniqueIndex:uq_generation_code" json:"code"`
	StartYear int    `gorm:"column:start_year;type:year" json:"start_year"`
	EndYear   int    `gorm:"column:end_year;type:year" json:"end_year"`
}
