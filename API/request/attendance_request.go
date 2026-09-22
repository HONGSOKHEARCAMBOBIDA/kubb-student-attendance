package request

type AttendanceRequestCreate struct {
	CompanyID       int    `json:"company_id"`
	Latitude        string `json:"latitude"`
	Longitude       string `json:"longitude"`
	ClassScheduleID int    `json:"class_schedule"`
}

type UpdateAttendanceRecordRequest struct {
	Status string `json:"status" validate:"required,oneof=PR P A"`
}
