package response

type LeaveRequestResponse struct {
	ID             int     `json:"id"`
	UserID         int     `json:"user_id" gorm:"column:user_id"`
	UserGender     int     `json:"gender" gorm:"column:gender"`
	UserNamekh     string  `json:"user_name_kh" gorm:"column:user_name_kh"`
	UserNameEn     string  `json:"user_name_en" gorm:"column:user_name_en"`
	UserCode       string  `json:"user_code"`
	ClassName      string  `json:"class_name"`
	StartDate      string  `json:"start_date" gorm:"column:start_date"`
	EndDate        string  `json:"end_date" gorm:"column:end_date"`
	BackToWorkDate string  `json:"back_to_work_date" gorm:"column:back_to_work_date"`
	TotalDay       float64 `json:"total_day" gorm:"column:total_day"`
	DeductTypeID   int     `json:"deduct_type_id" gorm:"column:deduct_type_id"`
	DeductTypeCode string  `json:"deduct_type_code"`
	DeductTypeName string  `json:"deduct_type_name"`
	Reason         string  `json:"reason" gorm:"column:reason"`
	Status         string  `json:"status" gorm:"column:status"`
	ApproveBy      int     `json:"approve_by" gorm:"column:approve_by"`
	ApproveByName  string  `json:"approve_by_name"`
	ApproveAt      string  `json:"approved_at" gorm:"column:approved_at"`
}

type NotPermissionLeave struct {
	UserID         int    `json:"user_id" gorm:"column:user_id"`
	UserGender     int    `json:"gender" gorm:"column:gender"`
	UserNamekh     string `json:"user_name_kh"`
	UserNameEn     string `json:"user_name_en"`
	UserCode       string `json:"user_code"`
	ClassID        int    `json:"class_id"`
	ClassName      string `json:"class_name"`
	MajorID        *int64 `gorm:"column:major_id" json:"major_id"`
	MajorName      string `json:"major_name"`
	ShiftID        *int64 `gorm:"column:shift_id" json:"shift_id"`
	ShiftName      string `json:"shift_name"`
	GenerationID   *int64 `gorm:"column:generation_id" json:"generation_id"`
	GenerationName string `json:"generation_name"`
	Year           *int8  `gorm:"column:year" json:"year"`
	Semester       *int8  `gorm:"column:semester" json:"semester"`
	Group          *int8  `gorm:"column:group" json:"group"`
	Term           *int8  `gorm:"column:term" json:"term"`
	ProgrammeID    *int64 `gorm:"column:programme_id" json:"programme_id"`
	ProgrammeName  string `json:"programme_name"`
}
