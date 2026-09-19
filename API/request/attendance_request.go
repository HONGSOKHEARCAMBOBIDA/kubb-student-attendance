package request

type AttendanceRequestCreate struct {
	CompanyID       int    `json:"company_id"`
	Latitude        string `json:"latitude"`
	Longitude       string `json:"longitude"`
	ClassScheduleID int    `json:"class_schedule"`
}
