package response

type AttendanceResponse struct {
	ID                       int                        `json:"id"`
	UserID                   int                        `json:"user_id"`
	Name                     string                     `json:"name"`
	Gender                   int                        `json:"gender"`
	GenderString             string                     `json:"gender_string"`
	CompanyID                int                        `json:"company_id"`
	CompanyName              string                     `json:"company_name"`
	RoleID                   int                        `json:"role_id"`
	RoleName                 string                     `json:"role_name" gorm:"column:role_name"`
	CheckDate                string                     `json:"check_date"`
	Status                   string                     `json:"status"`
	AttendanceRecordResponse []AttendanceRecordResponse `json:"attendance_record" gorm:"-"`
}

type AttendanceResponseDraft struct {
	Type          string `json:"type"`
	TypeString    string `json:"type_string"`
	ScheduledTime string `json:"scheduled_time"`
}

type AttendanceResponseGenerate struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	NameKH    string `json:"name_kh" gorm:"column:name_kh"`
	NameEN    string `json:"name_en" gorm:"column:name_en"`
	Code      string `json:"code" gorm:"column:code"`
	Gender    int    `json:"gender"`
	ClassID   int    `json:"class_id" gorm:"column:class_id"`
	ClassName string `json:"class_name"`

	CheckDate string `json:"check_date"`
	Status    string `json:"status"`

	Session1 string `gorm:"column:session1;not null;default:false" json:"session1"`
	Session2 string `gorm:"column:session2;not null;default:false" json:"session2"`
	Session3 string `gorm:"column:session3;not null;default:false" json:"session3"`
	Session4 string `gorm:"column:session4;not null;default:false" json:"session4"`
	Session5 string `gorm:"column:session5;not null;default:false" json:"session5"`
	Session6 string `gorm:"column:session6;not null;default:false" json:"session6"`

	Reason string `json:"reason"`
}

type AttendanceReportCell struct {
	Status    string `json:"status"` // "P", "A", or "" (not scheduled)
	CheckTime string `json:"check_time,omitempty"`
}

type AttendanceReportColumn struct {
	ColumnNo  int    `json:"column_no"` // global running number -> "Session N" header
	CheckDate string `json:"check_date"`
	Type      string `json:"type"` // session1..session6
}

type AttendanceReportRow struct {
	Index        int                    `json:"index"`
	UserID       int                    `json:"user_id"`
	NameKH       string                 `json:"name_kh"`
	NameEN       string                 `json:"name_en"`
	Code         string                 `json:"code"`
	Gender       int                    `json:"gender"`
	ClassName    string                 `json:"class_name"`
	Cells        []AttendanceReportCell `json:"cells"` // same order/length as Columns
	AbsentCount  int                    `json:"absent_count"`
	PresentCount int                    `json:"present_count"`
	Permission   int                    `json:"permission_count"`
}

type AttendanceReportResponse struct {
	Columns []AttendanceReportColumn `json:"columns"`
	Rows    []AttendanceReportRow    `json:"rows"`
}
