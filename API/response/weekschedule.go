package response

type WeekSchedule struct {
	ProgrammeName  string `json:"programme_name"`
	ShiftName      string `json:"shift_name"`
	GenerationName string `json:"generation_name"`
	MajorName      string `json:"major_name"`
	SubjectName    string `json:"subject_name_kh" gorm:"column:subject_name_kh"`
	NameKh         string `json:"name_kh"`
	Room           string `json:"room" gorm:"column:room"`
	ScheduleDate   string `json:"schedule_date"`
}
