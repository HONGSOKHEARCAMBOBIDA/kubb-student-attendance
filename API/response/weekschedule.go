package response

type WeekSchedule struct {
	ProgrammeName  string  `json:"programme_name"`
	ShiftName      string  `json:"shift_name"`
	GenerationName string  `json:"generation_name"`
	MajorName      string  `json:"major_name"`
	SubjectName    string  `json:"subject_name_kh" gorm:"column:subject_name_kh"`
	NameKH         string  `json:"name_kh" gorm:"column:name_kh"`
	NameEN         string  `json:"name_en" gorm:"column:name_en"`
	Gender         int     `json:"gender" gorm:"column:gender"`
	Phone          *string `gorm:"type:varchar(191)" json:"phone"`
	Room           string  `json:"room" gorm:"column:room"`
	ScheduleDate   string  `json:"schedule_date"`
	TotalStudent   int     `json:"total_student" gorm:"column:total_student"`
	TotalSession   int     `gorm:"column:total_session" json:"total_session"`
	IsMission      bool    `json:"is_mission" gorm:"column:is_mission"`
	DayOfWeek      int8    `gorm:"not null" json:"day_of_week"`
}
