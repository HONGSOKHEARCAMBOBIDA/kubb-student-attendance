package response

type AttendanceRecordResponse struct {
	ID            int    `json:"id"`
	AttendanceID  int    `json:"attendance_id"`
	CheckTime     string `json:"check_time"`
	Type          string `json:"type"`
	TypeString    string `json:"type_string"`
	Inzone        int    `json:"inzone"`
	Latitude      string `json:"latitdude" gorm:"column:latitdude"`
	Longitude     string `json:"longitude" gorm:"column:longitude"`
	ScheduledTime string `json:"scheduled_time"`
	Session1      string `gorm:"column:session1;not null;default:false" json:"session1"`
	Session2      string `gorm:"column:session2;not null;default:false" json:"session2"`
	Session3      string `gorm:"column:session3;not null;default:false" json:"session3"`
	Session4      string `gorm:"column:session4;not null;default:false" json:"session4"`
	Session5      string `gorm:"column:session5;not null;default:false" json:"session5"`
	Session6      string `gorm:"column:session6;not null;default:false" json:"session6"`
}
