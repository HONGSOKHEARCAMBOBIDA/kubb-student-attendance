package request

type FeeRequestCreate struct {
	StudentID      int    `json:"student_id"`
	ClassID        int    `json:"class_id"`
	SchoolarshipID int    `json:"scholarship_id" gorm:"column:scholarship_id"`
	MajorPriceID   int    `json:"major_price_id"`
	FeeScheduleID  int    `json:"fee_schedule_id"`
	Date           string `json:"date"`
}
