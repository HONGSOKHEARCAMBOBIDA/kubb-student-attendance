package response

type LeaveRequestResponse struct {
	ID             int     `json:"id"`
	UserID         int     `json:"user_id" gorm:"column:user_id"`
	UserGender     int     `json:"gender" gorm:"column:gender"`
	UserNamekh     string  `json:"user_name_kh"`
	UserNameEn     string  `json:"user_name_en"`
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
