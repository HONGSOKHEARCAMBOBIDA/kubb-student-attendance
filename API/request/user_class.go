package request

type UserClass struct {
	ClassID int64     `gorm:"column:class_id;not null" json:"class_id"`
	UserIDs []UserIDs `json:"user_id"`
}

type UserIDs struct {
	UserID int64 `gorm:"column:user_id;not null" json:"user_id"`
}
