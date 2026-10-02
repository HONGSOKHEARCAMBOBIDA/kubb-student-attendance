package model

type StudentProfile struct {
	UserID         int64   `gorm:"column:user_id;primaryKey" json:"user_id"`
	RegistrationNo *string `gorm:"column:registration_no;size:50" json:"registration_no"`
	DateOfBirth    string  `gorm:"column:date_of_birth;type:date" json:"date_of_birth"`
	PlaceOfBirth   string  `gorm:"column:place_of_birth;size:255" json:"place_of_birth"`
	Nationality    string  `gorm:"column:nationality;size:100" json:"nationality"`
	Campus         string  `gorm:"column:campus;size:100" json:"campus"`
}

func (StudentProfile) TableName() string {
	return "student_profile"
}
