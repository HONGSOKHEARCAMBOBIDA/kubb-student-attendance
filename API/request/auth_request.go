package request

type AuthRequest struct {
	Code     string `json:"code"`
	Password string `json:"password"`
}

type LoginQrRequest struct {
	QrToken string `json:"qr_token"`
}

type RegisterRequest struct {
	ClassID   int         `json:"class_id"`
	UserInput []UserInput `json:"users"`
}

type UserInput struct {
	NameKH string `json:"name_kh" gorm:"column:name_kh"`
	NameEN string `json:"name_en" gorm:"column:name_en"`
	Gender int    `json:"gender" gorm:"column:gender"`
	Code   string `json:"code" gorm:"column:code"`
	RoleID int    `json:"role_id"`
}

type RefreshTokenRequest struct {
	RefreshToken string `cookie:"refresh_token" binding:"required"`
}

type NewPasswordRequest struct {
	NewPassword string `json:"new_password"`
}

type UserRequestUpdate struct {
	NameKH *string `json:"name_kh" gorm:"column:name_kh"`
	NameEN *string `json:"name_en" gorm:"column:name_en"`
	Gender *int    `json:"gender" gorm:"column:gender"`
	Code   *string `json:"code" gorm:"column:code"`
	RoleID *int    `json:"role_id"`
}
