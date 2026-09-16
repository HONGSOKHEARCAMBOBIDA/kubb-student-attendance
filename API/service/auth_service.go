package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/response"
	"mysql/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	// "gorm.io/gorm/clause"
)

type AuthService interface {
	Login(input request.AuthRequest, c *gin.Context) (*response.AuthResponse, error)
	RefreshToken(refreshToken string, c *gin.Context) (*response.AuthResponse, error)
	Register(ctx context.Context, input request.RegisterRequest, c *gin.Context, userID int) error
	ToggleUserStatus(ctx context.Context, id int, userID int) error
	UpdateUser(ctx context.Context, input request.UserRequestUpdate, id int) error
	GetRole(ctx context.Context, id int) ([]model.Role, error)
	GetUserData(ctx context.Context, id int) (response.UserDataResponse, error)
	CreateUserClass(ctx context.Context, input request.UserClass, c *gin.Context, userID int) error
}

type authservice struct {
	db *gorm.DB
}

func NewAuthService() AuthService {
	return &authservice{
		db: config.DB,
	}
}

var requiredPermissions = []string{
	"add.payroll", "edit.payroll", "add.backup", "view.backup",
	"view.download.backup", "delete.backup", "add.company",
	"edit.company", "edit.user", "add.user", "edit.leave.type",
	"add.leave.type", "edit.leave.request", "edit.status.leave.request",
	"add.role.has.permission",
}

func (s *authservice) Login(input request.AuthRequest, c *gin.Context) (*response.AuthResponse, error) {
	key := "login_attempt:" + input.Code
	attempts, _ := utils.Redis.Get(utils.Ctx, key).Int()
	if attempts >= 100 {
		return nil, errors.New("អ្នកព្យាយាមចូលច្រើនពេក សូមព្យាយាមម្តងទៀតក្រោយ 10 នាទី")
	}
	var user model.User
	if err := s.db.Select("id,password, role_id, name_kh,code").
		Where("code = ? ", input.Code).
		First(&user).Error; err != nil {
		return nil, errors.New("ព័ត៌មានមិនត្រឹមត្រូវ ឬ អ្នកប្រើប្រាស់ត្រូវបានបិទគណនី")
	}

	var settings []model.Setting
	if err := s.db.Where("`key` IN ?", []string{
		"ACCESS_TOKEN_EXPIRE_HOURS",
		"REFRESH_TOKEN_EXPIRE_DAYS",
	}).Find(&settings).Error; err != nil {
		return nil, errors.New("Setting Not Found")
	}

	settingMap := make(map[string]string)
	for _, s := range settings {
		settingMap[s.Key] = s.Value
	}

	accesstoken, err := strconv.Atoi(settingMap["ACCESS_TOKEN_EXPIRE_HOURS"])
	if err != nil {
		return nil, errors.New("Bad request")
	}
	refreshtoken, err := strconv.Atoi(settingMap["REFRESH_TOKEN_EXPIRE_DAYS"])
	if err != nil {
		return nil, errors.New("Bad request")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Password),
	); err != nil {
		utils.Redis.Incr(utils.Ctx, key)
		utils.Redis.Expire(utils.Ctx, key, 10*time.Minute)
		return nil, errors.New("ព័ត៌មានមិនត្រឹមត្រូវ")
	}
	utils.Redis.Del(utils.Ctx, key)

	accessExpiry := time.Now().Add(time.Duration(accesstoken) * time.Hour)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"code":    user.Code,
		"role_id": user.RoleID,
		"exp":     accessExpiry.Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenStr, err := accessToken.SignedString(utils.Jwtkey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshTokenBytes := make([]byte, 32)
	if _, err := rand.Read(refreshTokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	refreshTokenStr := hex.EncodeToString(refreshTokenBytes)
	tokenPrefix := refreshTokenStr[:16]
	hashedRefresh := utils.HashToken(refreshTokenStr)
	if err := s.db.Where("user_id = ?", user.ID).Delete(&model.Session{}).Error; err != nil {
		return nil, fmt.Errorf("failed to delete session")
	}
	refreshExpiry := time.Now().Add(time.Duration(refreshtoken) * 24 * time.Hour)
	session := model.Session{
		UserID:       uint(user.ID),
		RefreshToken: string(hashedRefresh),
		TokenPrefix:  tokenPrefix,
		ExpiresAt:    refreshExpiry,
	}
	if err := s.db.Create(&session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	maxAge := int(time.Until(refreshExpiry).Seconds())
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		"refresh_token",
		refreshTokenStr,
		maxAge,
		"/",
		"",   // domain: leave empty for current host, or set "yourdomain.com" explicitly
		true, // Secure - HTTPS only
		true, // httpOnly - JS cannot read this
	)
	resp := &response.AuthResponse{
		AccessToken: accessTokenStr,
	}

	return resp, nil
}

func (s *authservice) RefreshToken(refreshToken string, c *gin.Context) (*response.AuthResponse, error) {
	if len(refreshToken) < 16 {
		return nil, errors.New("Invalid refresh token")
	}
	prefix := refreshToken[:16]

	var session model.Session
	err := s.db.Where("token_prefix = ?", prefix).First(&session).Error
	if err != nil {
		return nil, errors.New("Invalid or expired refresh token")
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("Invalid or expired refresh token")
	}

	if !utils.VerifyToken(session.RefreshToken, refreshToken) {
		return nil, errors.New("Invalid or expired refresh token")
	}

	var settings []model.Setting
	if err := s.db.Where("`key` IN ?", []string{
		"ACCESS_TOKEN_EXPIRE_HOURS",
		"REFRESH_TOKEN_EXPIRE_DAYS",
	}).Find(&settings).Error; err != nil {
		return nil, err
	}

	settingMap := make(map[string]string)
	for _, s := range settings {
		settingMap[s.Key] = s.Value
	}

	accesstoken, err := strconv.Atoi(settingMap["ACCESS_TOKEN_EXPIRE_HOURS"])
	if err != nil {
		return nil, err
	}

	refreshtoken, err := strconv.Atoi(settingMap["REFRESH_TOKEN_EXPIRE_DAYS"])
	if err != nil {
		return nil, errors.New("Bad request")
	}

	accessExpiry := time.Now().Add(time.Duration(accesstoken) * time.Minute)
	refreshExpiry := time.Now().Add(time.Duration(refreshtoken) * 24 * time.Hour)
	newRefreshBytes := make([]byte, 32)
	if _, err := rand.Read(newRefreshBytes); err != nil {
		return nil, errors.New("failed to generate refresh token")
	}
	newRefreshStr := hex.EncodeToString(newRefreshBytes)
	newHash := utils.HashToken(newRefreshStr)
	newPrefix := newRefreshStr[:16]

	if err := s.db.Model(&session).Updates(model.Session{
		RefreshToken: newHash,
		TokenPrefix:  newPrefix,
		ExpiresAt:    refreshExpiry,
	}).Error; err != nil {
		return nil, err
	}

	var user model.User
	if err := s.db.Select("id,role_id").Where("id = ?", session.UserID).First(&user).Error; err != nil {
		return nil, err
	}
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role_id": user.RoleID,
		"exp":     accessExpiry.Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(utils.Jwtkey)
	if err != nil {
		return nil, err
	}

	maxAge := int(time.Until(refreshExpiry).Seconds())
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		"refresh_token",
		newRefreshStr,
		maxAge,
		"/",
		"",   // domain: leave empty for current host, or set "yourdomain.com" explicitly
		true, // Secure - HTTPS only
		true, // httpOnly - JS cannot read this
	)

	return &response.AuthResponse{
		AccessToken: accessToken,
		//Permissions: permissions,
	}, nil
}

func (s *authservice) Register(ctx context.Context, input request.RegisterRequest, c *gin.Context, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	password := utils.HasPassword("kubb")

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(input.UserInput) > 0 {
			newdata := make([]model.User, 0, len(input.UserInput))
			for _, n := range input.UserInput {
				newdata = append(newdata, model.User{
					NameKH:   n.NameKH,
					NameEN:   n.NameEN,
					Gender:   n.Gender,
					Code:     n.Code,
					RoleID:   5,
					Password: password,
				})
			}

			if err := tx.Create(&newdata).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to create user", nil)
			}

			newUserClass := make([]model.UserClass, 0, len(newdata))
			for _, u := range newdata {
				newUserClass = append(newUserClass, model.UserClass{
					UserID:  int64(u.ID),
					ClassID: int64(input.ClassID),
				})
			}

			if err := tx.Create(&newUserClass).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to create user_class", nil)
			}
		}
		return nil
	})
	return err
}

func (s *authservice) CreateUserClass(ctx context.Context, input request.UserClass, c *gin.Context, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(input.UserIDs) > 0 {
			newdata := make([]model.UserClass, 0, len(input.UserIDs))
			for _, n := range input.UserIDs {
				newdata = append(newdata, model.UserClass{
					UserID:  n.UserID,
					ClassID: input.ClassID,
				})
			}

			if err := tx.Create(&newdata).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to create user", nil)
			}
		}
		return nil
	})
	return err
}

// func applyAccessFilter(query *gorm.DB, db *gorm.DB, role model.Role, user model.User) *gorm.DB {

// }

func applyCommonFilter(query *gorm.DB, filter map[string]string) *gorm.DB {
	for key, value := range filter {
		if value == "" {
			continue
		}
		switch key {
		case "name":
			query = query.Where("u.name_kh LIKE ?", "%"+value+"%")
		}
	}
	return query
}

var ErrCannotToggleOwnStatus = errors.New("cannot toggle your own status")

func (s *authservice) ToggleUserStatus(ctx context.Context, id int, userID int) error {
	if id == userID {
		return ErrCannotToggleOwnStatus
	}
	result := s.db.WithContext(ctx).Model(&model.User{}).Where("id =?", id).Update("is_active", gorm.Expr("NOT is_active"))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *authservice) ChangePassword(ctx context.Context, userID int, input request.NewPasswordRequest) error {
	var user model.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return err
	}

	hash := utils.HasPassword(input.NewPassword)
	if err := s.db.WithContext(ctx).Model(&user).Update("password_hash", hash).Error; err != nil {
		return err
	}
	return nil

}

func (s *authservice) UpdateUser(ctx context.Context, input request.UserRequestUpdate, id int) error {
	updates := map[string]interface{}{}

	if input.NameKH != nil {
		updates["name_kh"] = *&input.NameKH
	}
	if input.NameEN != nil {
		updates["name_en"] = *&input.NameEN
	}
	if input.Gender != nil {
		updates["gender"] = *input.Gender
	}
	if input.Code != nil {
		updates["code"] = *input.Code
	}
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	committed := false
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if !committed {
			tx.Rollback()
		}
	}()

	if len(updates) > 0 {
		if err := tx.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			slog.Error("failed to update user", "error", err)
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	committed = true
	return nil
}

func (s *authservice) CountUser(ctx context.Context, id int) (response.UserCount, error) {
	var countUser response.UserCount
	var user model.User
	if err := s.db.WithContext(ctx).Preload("Role").First(&user, id).Error; err != nil {
		return response.UserCount{}, err
	}

	userQuery := s.db.WithContext(ctx).
		Table("user u").
		Select(`
            COUNT(DISTINCT CASE WHEN u.is_active = '1' THEN u.id END) AS total
        `)

	userQuery = helper.ApplyAccessFilter(userQuery, s.db, user.Role, user)
	if err := userQuery.Scan(&countUser).Error; err != nil {
		return response.UserCount{}, err
	}
	return countUser, nil
}

func (s *authservice) GetRole(ctx context.Context, id int) ([]model.Role, error) {
	var role []model.Role
	var user model.User
	if err := s.db.WithContext(ctx).Preload("Role").First(&user, id).Error; err != nil {
		return nil, err
	}
	roleQuery := s.db.WithContext(ctx).Table("role r").
		Select(`
		 r.id AS id,
		 r.name AS name,
		 r.display_name AS display_name
	`)
	roleQuery = helper.ApplyAccessGetRole(roleQuery, s.db, user.Role, user)

	if err := roleQuery.Scan(&role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (s *authservice) GetUserData(ctx context.Context, id int) (response.UserDataResponse, error) {
	var userdata response.UserDataResponse

	err := s.db.WithContext(ctx).
		Table("user u").
		Select(`
			u.id AS id,
			u.name_kh AS name,
			u.role_id AS role_id
		`).
		Where("u.id = ?", id).
		First(&userdata).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return userdata, fmt.Errorf("user with id %d not found", id)
	}
	if err != nil {
		return userdata, fmt.Errorf("failed to get user data: %w", err)
	}

	var classID int64
	result := s.db.WithContext(ctx).
		Table("user_class uc").
		Where("uc.user_id = ? AND uc.is_active = 1", id).
		Select("uc.class_id").
		Scan(&classID)
	if result.Error != nil {
		return userdata, fmt.Errorf("failed to get user classes: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return userdata, fmt.Errorf("no active class found for user %d", id)
	}
	userdata.ClassID = int(classID)

	var permissions []model.Permission
	if err := s.db.WithContext(ctx).
		Table("permission p").
		Select("p.name AS name").
		Joins("JOIN role_permission rp ON rp.permission_id = p.id").
		Where("rp.role_id = ? AND p.name IN ?", userdata.RoleID, requiredPermissions).
		Scan(&permissions).Error; err != nil {
		return userdata, fmt.Errorf("failed to get user permissions: %w", err)
	}
	userdata.Permissions = permissions

	return userdata, nil
}

// func (s *authservice) GetUserApprove(ctx context.Context, id int) ([]response.UserApprove, error) {
// 	var user model.User
// 	if err := s.db.WithContext(ctx).
// 		Preload("Role").
// 		Preload("Company").
// 		First(&user, id).Error; err != nil {
// 		return nil, err
// 	}

// 	var userApproves []response.UserApprove
// 	err := s.db.WithContext(ctx).
// 		Table("user u").
// 		Select(`
// 			u.id AS id,
// 			u.name AS user_name,
// 			u.company_id
// 		`).
// 		Joins("JOIN role r ON r.id = u.role_id").
// 		Where("u.company_id = ? AND r.level >= ?", user.CompanyID, user.Role.Level).
// 		Scan(&userApproves).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return userApproves, nil
// }

// func (s *authservice) VerifyUser(ctx context.Context, id int) error {
// 	if id <= 0 {
// 		return fmt.Errorf("invalid id: %d", id)
// 	}

// 	tx := s.db.WithContext(ctx).Begin()
// 	if tx.Error != nil {
// 		return fmt.Errorf("failed to start transaction: %w", tx.Error)
// 	}

// 	committed := false
// 	defer func() {
// 		if !committed {
// 			tx.Rollback()
// 		}
// 	}()

// 	var user model.User
// 	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, id).Error; err != nil {
// 		return fmt.Errorf("failed to lock row: %w", err)
// 	}

// 	result := tx.Model(&model.User{}).
// 		Where("id = ?", id).
// 		Update("is_verify", !user.IsVerify)
// 	if result.Error != nil {
// 		return fmt.Errorf("failed to verify user: %w", result.Error)
// 	}
// 	if result.RowsAffected == 0 {
// 		return errors.New("failed to update verification status")
// 	}

// 	if err := tx.Commit().Error; err != nil {
// 		return fmt.Errorf("failed to commit transaction: %w", err)
// 	}
// 	committed = true
// 	return nil
// }

func (s *authservice) Logout(ctx context.Context, userID int) error {
	if userID <= 0 {
		return fmt.Errorf("invalid user id: %d", userID)
	}

	result := s.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.Session{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete sessions for user %d: %w", userID, result.Error)
	}
	return nil
}
