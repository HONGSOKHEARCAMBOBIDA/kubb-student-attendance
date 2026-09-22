package request

// ClassScheduleRequestCreate is used when a class admin attaches a subject
// to a specific day of the week for a class. SubjectID must belong to the
// class's major for the class's year/semester (validated in the service),
// so the frontend should only offer subjects returned by
// ClassScheduleService.GetAvailableSubjects.
type ClassScheduleRequestCreate struct {
	SubjectID    int64  `json:"subject_id" binding:"required"`
	DayOfWeek    int8   `json:"day_of_week" binding:"required,min=1,max=7"` // 1=Monday ... 7=Sunday
	ScheduleDate string `gorm:"column:schedule_date" json:"schedule_date"`
	TotalSession int    `gorm:"column:total_session" json:"total_session"`
}
