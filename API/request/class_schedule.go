package request

// ClassScheduleRequestCreate is used when a class admin attaches a subject
// to a specific day of the week for a class. SubjectID must belong to the
// class's major for the class's year/semester (validated in the service),
// so the frontend should only offer subjects returned by
// ClassScheduleService.GetAvailableSubjects.
type ClassScheduleRequestCreate struct {
	ClassDetailID int     `json:"class_detail_id" gorm:"column:class_detail_id"`
	SubjectID     int64   `json:"subject_id" binding:"required"`
	DayOfWeek     int8    `json:"day_of_week" binding:"required,min=1,max=7"` // 1=Monday ... 7=Sunday
	ScheduleDate  string  `gorm:"column:schedule_date" json:"schedule_date"`
	TotalSession  int     `gorm:"column:total_session" json:"total_session"`
	TeacherID     int     `json:"teacher_id"`
	HourlyRate    float64 `json:"hourly_rate"`
	EffectiveDate string  `gorm:"type:date;not null;index:idx_teacher_rates_effective_date" json:"effective_date"`
	EndDate       *string `gorm:"type:date" json:"end_date"`
	TotalStudent  int     `json:"total_student" gorm:"column:total_student"`
}
