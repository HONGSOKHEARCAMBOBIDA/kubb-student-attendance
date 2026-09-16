package request

type CompanyRequestCreate struct {
	Name             string `json:"name" bind:"required"`
	MapLink          string `json:"map_link"`
	Radius           string `json:"radius" bind:"required"`
	GroupLink        string `json:"group_link"`
	BotToken         string `json:"bot_token"`
	Currency         string `json:"currency"`
	LatePenalty      string `json:"late_penalty" gorm:"column:late_penalty"`
	LeftEarlyPenalty string `json:"left_early_penalty" gorm:"column:left_early_penalty"`
	CanScanOutsize   bool   `json:"can_scan_outsize" gorm:"column:can_scan_outsize"`
	Color            string `json:"color" gorm:"column:color"`
	TotalWorkDay     int    `json:"total_work_day" gorm:"column:total_work_day"`
}

type ClassRequestCreate struct {
	Name           string `json:"name" bind:"required"`
	MapLink        string `json:"map_link"`
	Radius         string `json:"radius" bind:"required"`
	CanScanOutsize bool   `json:"can_scan_outsize" gorm:"column:can_scan_outsize"`
	MajorID        *int64 `gorm:"column:major_id" json:"major_id"`
	ShiftID        *int64 `gorm:"column:shift_id" json:"shift_id"`
	GenerationID   *int64 `gorm:"column:generation_id" json:"generation_id"`
	Year           int    `gorm:"column:year" json:"year"`
	Semester       int    `gorm:"column:semester" json:"semester"`
	Group          int    `gorm:"column:group" json:"group"`
	Term           int    `gorm:"column:term" json:"term"`
	ProgrammeID    *int64 `gorm:"column:programme_id" json:"programme_id"`
}

type ClassRequestUpdate struct {
	Name           string `json:"name" bind:"required"`
	MapLink        string `json:"map_link"`
	Radius         string `json:"radius" bind:"required"`
	CanScanOutsize bool   `json:"can_scan_outsize" gorm:"column:can_scan_outsize"`
	MajorID        *int64 `gorm:"column:major_id" json:"major_id"`
	ShiftID        *int64 `gorm:"column:shift_id" json:"shift_id"`
	GenerationID   *int64 `gorm:"column:generation_id" json:"generation_id"`
	Year           int    `gorm:"column:year" json:"year"`
	Semester       int    `gorm:"column:semester" json:"semester"`
	Group          int    `gorm:"column:group" json:"group"`
	Term           int    `gorm:"column:term" json:"term"`
	ProgrammeID    *int64 `gorm:"column:programme_id" json:"programme_id"`
}

type CompanyRequesUpdate struct {
	Name             *string `json:"name"`
	MapLink          *string `json:"map_link"`
	Latitude         *string `json:"latitude"`
	Longitude        *string `json:"longitude"`
	Radius           *string `json:"radius"`
	Currency         *string `json:"currency"`
	LatePenalty      *string `json:"late_penalty" gorm:"column:late_penalty"`
	LeftEarlyPenalty *string `json:"left_early_penalty" gorm:"column:left_early_penalty"`
	CanScanOutsize   *int    `json:"can_scan_outsize" gorm:"column:can_scan_outsize"`
	Color            *string `json:"color" gorm:"column:color"`
	TotalWorkDay     *int    `json:"total_work_day" gorm:"column:total_work_day"`
}

type CompanyRequestUpdateTelegram struct {
	BotToken  *string `json:"bot_token"`
	GroupLink *string `json:"group_link"`
}
