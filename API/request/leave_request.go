package request

type LeaveRequestCreate struct {
	ClassID        int     `gorm:"not null;index:idx_leave_request_class_id" json:"class_id"`
	StartDate      string  `gorm:"type:date;not null;index:idx_leave_request_date" json:"start_date"`
	EndDate        string  `gorm:"type:date;not null;index:idx_leave_request_date" json:"end_date"`
	BackToWorkDate string  `gorm:"type:date;not null" json:"back_to_work_date"`
	TotalDay       float64 `gorm:"type:decimal(5,2);not null;default:0.00" json:"total_day"`
	DeductTypeID   int     `gorm:"column:deduct_type_id" json:"deduct_type_id"`
	Reason         string  `gorm:"type:text" json:"reason"`
}

type LeaveRequestUpdate struct {
	ClassID        int     `gorm:"not null;index:idx_leave_request_class_id" json:"class_id"`
	StartDate      string  `gorm:"type:date;not null;index:idx_leave_request_date" json:"start_date"`
	EndDate        string  `gorm:"type:date;not null;index:idx_leave_request_date" json:"end_date"`
	BackToWorkDate string  `gorm:"type:date;not null" json:"back_to_work_date"`
	TotalDay       float64 `gorm:"type:decimal(5,2);not null;default:0.00" json:"total_day"`
	DeductTypeID   int     `gorm:"column:deduct_type_id" json:"deduct_type_id"`
	Reason         string  `gorm:"type:text" json:"reason"`
}

type LeaveRequestUpdateStatus struct {
	Status *int `json:"status" gorm:"column:status"`
}

type LeaveRequestDetailRequest struct {
	Type string `gorm:"type:enum('session1','session2','session3','session4','session5');not null;index:idx_leave_request_details_type" json:"type"`
}

type NotPermissionLeaveRequest struct {
	CheckDate               string                    `json:"check_date"`
	NotPermissionLeaveInput []NotPermissionLeaveInput `json:"data"`
}

type NotPermissionLeaveInput struct {
	UserID  int `json:"user_id"`
	ClassID int `json:"class_id"`
	ShiftID int `json:"shift_id"`
}
