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
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	// "gorm.io/gorm/clause"
)

type AuthService interface {
	Login(input request.AuthRequest, c *gin.Context) (*response.AuthResponse, error)
	RefreshToken(refreshToken string, c *gin.Context) (*response.AuthResponse, error)
	Register(ctx context.Context, input request.RegisterRequest, c *gin.Context, userID int) error
	RegisterMain(ctx context.Context, input request.UserInput) error
	ToggleUserStatus(ctx context.Context, id int, userID int) error
	UpdateUser(ctx context.Context, input request.UserRequestUpdate, id int) error
	GetRole(ctx context.Context, id int) ([]model.Role, error)
	GetUserData(ctx context.Context, id int) (response.UserDataResponse, error)
	CreateUserClass(ctx context.Context, input request.UserClass, c *gin.Context, userID int) error
	UpdateUserClass(ctx context.Context, id int, input request.UserClassUpdateStatus) error
	GetUserIncludeStudent(ctx context.Context, id int, pf request.Pagination, filter map[string]string) ([]response.UserResponseNotStudent, *model.PaginationMetadata, error)
	GetStudentCategory(ctx context.Context) ([]model.StudentCategory, error)
	AddStudentProfile(ctx context.Context, input request.StudentProfileRequestCreate) error
	GetStudentProfile(ctx context.Context, id int) (*response.StudentProfile, error)
	UpdateStudentProfile(ctx context.Context, input request.StudentProfileRequestUpdate) error
	AddStudentFamily(ctx context.Context, input request.StudentFamily) error
	GetStudentFamily(ctx context.Context, id int) (*model.StudentFamily, error)
	UpdateStudentFamily(ctx context.Context, input request.StudentFamily) error
	AddStudentEducation(ctx context.Context, input request.StudentEducation) error
	GetStudentEducation(ctx context.Context, id int) (*[]response.StudentEducation, error)
	UpdateStudentEducation(ctx context.Context, id int, input request.StudentEducation) error
	GetDocumentType(ctx context.Context) ([]model.DocumentType, error)
	AddStudentDocument(ctx context.Context, input request.StudentDocumentRequestCreate) error
	UpdateStudentDocument(ctx context.Context, id int, input request.StudentDocumentRequestUpdate) error
	GetStudentDocuments(ctx context.Context, id int) ([]response.StudentDocument, error)
	AddTeacherProfile(ctx context.Context, input request.TeacherProfile) error
	UpdateTeacherProfile(ctx context.Context, input request.TeacherProfile) error
	GetTeacherProfile(ctx context.Context, id int) (*response.TeacherProfile, error)
	UpdateTeacherRate(ctx context.Context, classscheduleID int, input request.TeacherRate) error
	GetStaft(ctx context.Context) ([]response.UserResponseNotStudent, error)
	AddStudentExam(ctx context.Context, input request.StudentExam) error
	GetStudentExam(ctx context.Context, detailID int) ([]model.StudentExam, error)
	DeleteStudentExam(ctx context.Context, id int) error
	GetStudentStudyHistory(ctx context.Context, id int) (response.StudentStudyHistory, error)
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
	"add.backup", "view.backup",
	"view.download.backup", "delete.backup", "add.company",
	"edit.company", "edit.user", "add.user", "edit.leave.type",
	"add.leave.type", "edit.leave.request", "edit.status.leave.request",
	"add.role.has.permission", "add.leave.request", "edit.leave.request",
	"edit.status.leave.request", "delete.leave.request", "add.subject", "edit.subject",
	"add.Major", "edit.Major", "add.score",
}

func (s *authservice) RegisterMain(ctx context.Context, input request.UserInput) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	password := utils.HasPassword("kubb")

	roleID := input.RoleID
	if roleID == 0 {
		roleID = 5
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newdata := model.User{
			NameKH:   input.NameKH,
			NameEN:   input.NameEN,
			Gender:   input.Gender,
			Code:     input.Code,
			Password: password,
			RoleID:   roleID,
		}
		if err := tx.Create(&newdata).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create user", nil)
		}
		return nil
	})
	return err
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
					Status:  model.UserClassStatusSTUDY,
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
					UserID:    n.UserID,
					ClassID:   input.ClassID,
					Status:    model.UserClassStatusSTUDY,
					FeeStatus: model.FeeStatusUserClassACTIVE,
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
	if input.RoleID != nil {
		updates["role_id"] = *input.RoleID
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

func applyAccessFilterUserIncludeStudent(query *gorm.DB, db *gorm.DB, user model.User) *gorm.DB {
	return query.Where("r.level <= ?", user.Role.Level)
}

func applyCommonFilterUserIncludeStudent(query *gorm.DB, filter map[string]string) *gorm.DB {
	for key, value := range filter {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		switch key {
		case "name":
			query = query.Where("u.name_kh LIKE ? OR u.name_en LIKE ?", "%"+value+"%", "%"+value+"%")
		case "role_id":
			query = query.Where("u.role_id =?", value)
		}
	}
	return query
}

func (s *authservice) GetUserIncludeStudent(ctx context.Context, id int, pf request.Pagination, filter map[string]string) ([]response.UserResponseNotStudent, *model.PaginationMetadata, error) {
	var data []response.UserResponseNotStudent
	var user model.User
	if err := s.db.WithContext(ctx).Preload("Role").First(&user, id).Error; err != nil {
		return nil, nil, err
	}
	offset := (pf.Page - 1) * pf.PageSize

	q := s.db.WithContext(ctx).Table("user u").
		Select(`
		u.id AS id,
		u.name_kh AS name_kh,
		u.name_en AS name_en,
		u.gender AS gender,
		u.code AS code,
		r.id AS role_id,
		r.display_name AS role_name,
		r.name AS role_for_check
	`).
		Joins("LEFT JOIN role r ON r.id = u.role_id")
	q = applyAccessFilterUserIncludeStudent(q, s.db, user)
	q = applyCommonFilterUserIncludeStudent(q, filter)
	var totalCount int64
	if err := q.Count(&totalCount).Error; err != nil {
		return nil, nil, err
	}
	if err := q.Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, err
	}
	return data, helper.BuildPaginationMeta(pf, totalCount), nil
}

func (s *authservice) GetUserData(ctx context.Context, id int) (response.UserDataResponse, error) {
	var userdata response.UserDataResponse

	err := s.db.WithContext(ctx).
		Table("user u").
		Select(`
			u.id AS id,
			u.name_kh AS name,
			u.role_id AS role_id,
			r.level AS level
		`).
		Joins("LEFT JOIN role r ON r.id = u.role_id").
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
	if result.RowsAffected > 0 {
		userdata.ClassID = int(classID)
	}

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

func (s *authservice) UpdateUserClass(ctx context.Context, id int, input request.UserClassUpdateStatus) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var data model.UserClass
		if err := tx.Where("id = ?", id).First(&data).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "classcurriculumn not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch classcurriculumn", nil)
		}
		data.Status = input.Status
		if err := tx.Save(&data).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update student", nil)
		}
		return nil
	})
	return err
}

func (s *authservice) GetStudentCategory(ctx context.Context) ([]model.StudentCategory, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	var data []model.StudentCategory

	err := s.db.WithContext(ctx).Find(&data).Error
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, apperror.New(
			apperror.CodeNotFound,
			"student category not found",
			nil,
		)
	}

	return data, nil
}

func (s *authservice) AddStudentProfile(ctx context.Context, input request.StudentProfileRequestCreate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newdata := model.StudentProfile{
			UserID:            input.UserID,
			StudentCategoryID: input.StudentCategoryID,
			DateOfBirth:       input.DateOfBirth,
			Nationality:       input.Nationality,
			Phone:             input.Phone,
			VillageID:         input.VillageID,
			Occupation:        input.Occupation,
			AcademicStream:    input.AcademicStream,
			TelegramUsername:  input.TelegramUsername,
			ExamIn:            false,
			ExamOut:           false,
		}
		if err := tx.Create(&newdata).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create student profile", nil)
		}
		return nil
	})
	return err
}

func (s *authservice) UpdateStudentProfile(
	ctx context.Context,
	input request.StudentProfileRequestUpdate,
) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		data := model.StudentProfile{
			StudentCategoryID: input.StudentCategoryID,
			DateOfBirth:       input.DateOfBirth,
			Nationality:       input.Nationality,
			Phone:             input.Phone,
			VillageID:         input.VillageID,
			Occupation:        input.Occupation,
			AcademicStream:    input.AcademicStream,
			TelegramUsername:  input.TelegramUsername,
		}

		result := tx.Model(&model.StudentProfile{}).
			Where("user_id = ?", input.UserID).
			Updates(&data)

		if result.Error != nil {
			return apperror.New(
				apperror.CodeInternal,
				"failed to update student profile",
				nil,
			)
		}

		if result.RowsAffected == 0 {
			return apperror.New(
				apperror.CodeNotFound,
				"student profile not found",
				nil,
			)
		}

		return nil
	})

	return err
}

func (s *authservice) GetStudentProfile(
	ctx context.Context,
	id int,
) (*response.StudentProfile, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	var data response.StudentProfile

	err := s.db.WithContext(ctx).
		Table("student_profile sp").
		Where("user_id = ?", id).
		Select(`
			sp.id AS id,
			sp.user_id AS user_id,
			sp.student_category_id AS student_category_id,
			sp.date_of_birth AS date_of_birth,
			sp.nationality AS nationality,
			sp.phone AS phone,
			sp.village_id,
			c.id AS commune_id,
			d.id AS district_id,
			p.id AS province_id,
			sp.occupation AS occupation,
			sp.academic_stream AS academic_stream,
			sp.telegram_username AS telegram_username
		`).
		Joins("LEFT JOIN villages v ON v.id = sp.village_id").
		Joins("LEFT JOIN communes c ON c.id = v.commune_id").
		Joins("LEFT JOIN districts d ON d.id = c.district_id").
		Joins("LEFT JOIN provinces p ON p.id = d.province_id").Scan(&data).Error
	if err != nil {
		return nil, err
	}

	data.DateOfBirth = helper.FormatDate(data.DateOfBirth)

	return &data, nil
}

func (s *authservice) AddStudentFamily(ctx context.Context, input request.StudentFamily) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newdata := model.StudentFamily{
			StudentID:         input.StudentID,
			FatherName:        input.FatherName,
			FatherEnglishName: input.FatherEnglishName,
			FatherAge:         input.FatherAge,
			FatherIsAlive:     input.FatherIsAlive,
			FatherPhoneNumber: input.FatherPhoneNumber,
			FatherOccupation:  input.FatherOccupation,
			FatherWorkplace:   input.FatherWorkplace,
			MotherName:        input.MotherName,
			MotherEnglishName: input.MotherEnglishName,
			MotherAge:         input.MotherAge,
			MotherIsAlive:     input.MotherIsAlive,
			MotherPhoneNumber: input.MotherPhoneNumber,
			MotherOccupation:  input.MotherOccupation,
			MotherWorkplace:   input.MotherWorkplace,
		}
		if err := tx.Create(&newdata).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create student profile", nil)
		}
		return nil
	})
	return err
}

func (s *authservice) UpdateStudentFamily(ctx context.Context, input request.StudentFamily) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		data := model.StudentFamily{
			FatherName:        input.FatherName,
			FatherEnglishName: input.FatherEnglishName,
			FatherAge:         input.FatherAge,
			FatherIsAlive:     input.FatherIsAlive,
			FatherPhoneNumber: input.FatherPhoneNumber,
			FatherOccupation:  input.FatherOccupation,
			FatherWorkplace:   input.FatherWorkplace,
			MotherName:        input.MotherName,
			MotherEnglishName: input.MotherEnglishName,
			MotherAge:         input.MotherAge,
			MotherIsAlive:     input.MotherIsAlive,
			MotherPhoneNumber: input.MotherPhoneNumber,
			MotherOccupation:  input.MotherOccupation,
			MotherWorkplace:   input.MotherWorkplace,
		}
		result := tx.Model(&model.StudentFamily{}).
			Where("student_id = ?", input.StudentID).
			Updates(&data)

		if result.Error != nil {
			return apperror.New(
				apperror.CodeInternal,
				"failed to update student family",
				nil,
			)
		}

		if result.RowsAffected == 0 {
			return apperror.New(
				apperror.CodeNotFound,
				"student profile not found",
				nil,
			)
		}

		return nil
	})

	return err
}

func (s *authservice) GetStudentFamily(ctx context.Context, id int) (*model.StudentFamily, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	var data model.StudentFamily
	err := s.db.WithContext(ctx).Table("student_families").
		Where("student_id =?", id).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *authservice) AddStudentEducation(ctx context.Context, input request.StudentEducation) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newdata := model.StudentEducation{
			StudentID:       input.StudentID,
			Level:           input.Level,
			SchoolName:      input.SchoolName,
			StartDate:       input.StartDate,
			EndDate:         input.EndDate,
			CertificateDate: input.CertificateDate,
			Score:           input.Score,
			GPA:             input.GPA,
			Grade:           input.Grade,
		}
		if err := tx.Create(&newdata).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create student profile", nil)
		}
		return nil
	})
	return err
}

func (s *authservice) GetStudentEducation(
	ctx context.Context,
	id int,
) (*[]response.StudentEducation, error) {
	ctx, cancel := context.WithTimeout(
		ctx,
		utils.DefaultQueryTimeout,
	)
	defer cancel()

	var data []response.StudentEducation

	err := s.db.WithContext(ctx).
		Table("student_education").
		Where("student_id = ?", id).
		Find(&data).Error

	if err != nil {
		return nil, err
	}

	for i := range data {
		if data[i].StartDate != nil {
			formatted := helper.FormatDate(*data[i].StartDate)
			data[i].StartDate = &formatted
		}

		if data[i].EndDate != nil {
			formatted := helper.FormatDate(*data[i].EndDate)
			data[i].EndDate = &formatted
		}

		if data[i].CertificateDate != nil {
			formatted := helper.FormatDate(*data[i].CertificateDate)
			data[i].CertificateDate = &formatted
		}
	}

	return &data, nil
}

func (s *authservice) UpdateStudentEducation(ctx context.Context, id int, input request.StudentEducation) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		data := model.StudentEducation{
			Level:           input.Level,
			SchoolName:      input.SchoolName,
			StartDate:       input.StartDate,
			EndDate:         input.EndDate,
			CertificateDate: input.CertificateDate,
			Score:           input.Score,
			GPA:             input.GPA,
			Grade:           input.Grade,
		}
		result := tx.Model(&model.StudentEducation{}).
			Where("id = ?", id).
			Updates(&data)

		if result.Error != nil {
			return apperror.New(
				apperror.CodeInternal,
				"failed to update student education",
				nil,
			)
		}

		if result.RowsAffected == 0 {
			return apperror.New(
				apperror.CodeNotFound,
				"student educdation not found",
				nil,
			)
		}

		return nil
	})

	return err
}

func (s *authservice) GetDocumentType(ctx context.Context) ([]model.DocumentType, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	var data []model.DocumentType

	err := s.db.WithContext(ctx).Find(&data).Error
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, apperror.New(
			apperror.CodeNotFound,
			"document type not found",
			nil,
		)
	}

	return data, nil
}

func (s *authservice) AddStudentDocument(ctx context.Context, input request.StudentDocumentRequestCreate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newdata := model.StudentDocument{
			StudentID:      input.StudentID,
			DocumentTypeID: input.DocumentTypeID,
			RequiredQty:    input.RequiredQty,
			ReceivedQty:    input.ReceivedQty,
			Remark:         input.Remark,
		}
		if err := tx.Create(&newdata).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create student document", nil)
		}
		return nil
	})
	return err
}

func (s *authservice) UpdateStudentDocument(ctx context.Context, id int, input request.StudentDocumentRequestUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		data := model.StudentDocument{
			DocumentTypeID: input.DocumentTypeID,
			RequiredQty:    input.RequiredQty,
			ReceivedQty:    input.ReceivedQty,
			Remark:         input.Remark,
		}
		update := tx.Model(&model.StudentDocument{}).Where("id =?", id).Updates(&data)
		if update.Error != nil {
			return apperror.New(
				apperror.CodeInternal,
				"failed to update student document",
				nil,
			)
		}

		if update.RowsAffected == 0 {
			return apperror.New(
				apperror.CodeNotFound,
				"student document not found",
				nil,
			)
		}
		return nil
	})
	return err
}

func (s *authservice) GetStudentDocuments(ctx context.Context, id int) ([]response.StudentDocument, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	var docs []response.StudentDocument

	err := s.db.WithContext(ctx).
		Table("student_documents AS sd").
		Select(`
			sd.id            AS id,
			d.id             AS document_type_id,
			d.name_kh        AS document_type_name,
			sd.required_qty  AS required_qty,
			sd.received_qty  AS received_qty,
			sd.remark        AS remark
		`).
		Where("sd.student_id =?", id).
		Joins("LEFT JOIN document_types AS d ON d.id = sd.document_type_id").
		Scan(&docs).Error
	if err != nil {
		return nil, fmt.Errorf("get student documents: %w", err)
	}

	return docs, nil
}

func (s *authservice) AddTeacherProfile(ctx context.Context, input request.TeacherProfile) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newdata := model.TeacherProfile{
			UserID:       input.UserID,
			DateOfBirth:  input.DateOfBirth,
			PlaceOfBirth: input.PlaceOfBirth,
			Nationality:  input.Nationality,
			Address:      input.Address,
			Phone:        input.Phone,
		}
		if err := tx.Create(&newdata).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create student document", nil)
		}
		return nil
	})
	return err
}

func (s *authservice) UpdateTeacherProfile(ctx context.Context, input request.TeacherProfile) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		data := model.TeacherProfile{
			DateOfBirth:  input.DateOfBirth,
			PlaceOfBirth: input.PlaceOfBirth,
			Nationality:  input.Nationality,
			Address:      input.Address,
			Phone:        input.Phone,
		}
		update := tx.Model(&model.TeacherProfile{}).Where("user_id =?", input.UserID).Updates(&data)
		if update.Error != nil {
			return apperror.New(
				apperror.CodeInternal,
				"failed to update teacher profile",
				nil,
			)
		}

		if update.RowsAffected == 0 {
			return apperror.New(
				apperror.CodeNotFound,
				"teacher profile not found",
				nil,
			)
		}
		return nil
	})
	return err
}

func (s *authservice) GetTeacherProfile(ctx context.Context, id int) (*response.TeacherProfile, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var data response.TeacherProfile
	err := s.db.WithContext(ctx).Table("teacher_profile").
		Where("user_id =?", id).First(&data).Error
	if err != nil {
		return nil, err
	}
	*data.DateOfBirth = helper.FormatDate(*data.DateOfBirth)
	return &data, nil
}

func (s *authservice) UpdateTeacherRate(
	ctx context.Context,
	classscheduleID int,
	input request.TeacherRate,
) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		teacherRateUpdates := map[string]interface{}{
			"user_id":     input.UserID,
			"hourly_rate": input.HourlyRate,
		}
		result := tx.Model(&model.TeacherRate{}).
			Where("class_schedule_id = ?", classscheduleID).
			Updates(teacherRateUpdates)

		if result.Error != nil {
			return result.Error
		}

		classScheduleUpdates := map[string]interface{}{
			"subject_id":    input.SubjectID,
			"day_of_week":   input.DayOfWeek,
			"schedule_date": input.ScheduleDate,
			"total_session": input.TotalSession,
			"total_student": input.TotalStudent,
			"room":          input.Room,
			"is_mission":    input.IsMission,
		}

		scheduleResult := tx.Model(&model.ClassSchedule{}).
			Where("id = ?", classscheduleID).
			Updates(classScheduleUpdates)

		if scheduleResult.Error != nil {
			return scheduleResult.Error
		}

		return nil
	})
}

func (s *authservice) GetStaft(ctx context.Context) ([]response.UserResponseNotStudent, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var data []response.UserResponseNotStudent
	err := s.db.WithContext(ctx).Table("user u").
		Joins("LEFT JOIN role r ON r.id = u.role_id").
		Where("r.level >= ?", 2).
		Select(`
		u.id AS id,
		u.name_kh AS name_kh,
		u.name_en AS name_en,
		u.gender AS gender,
		r.display_name AS role_name
	`).Scan(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *authservice) AddStudentExam(ctx context.Context, input request.StudentExam) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	data := make([]model.StudentExam, 0, len(input.Student))
	for _, st := range input.Student {
		data = append(data, model.StudentExam{
			ExamPlatformDetailID: input.ExamPlatformDetailID,
			ClassID:              st.ClassID,
			StudentID:            st.StudentID,
			TableNumber:          st.TableNumber,
			Room:                 input.Room,
		})
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "exam_platform_detail_id"}, {Name: "student_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"table_number", "room", "class_id"}),
		}).Create(&data).Error
	})
}

func (s *authservice) GetStudentExam(ctx context.Context, detailID int) ([]model.StudentExam, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var rows []model.StudentExam
	err := s.db.WithContext(ctx).Where("exam_platform_detail_id = ?", detailID).Find(&rows).Error
	return rows, err
}

func (s *authservice) DeleteStudentExam(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	if err := s.db.WithContext(ctx).Where("id = ?", id).Delete(&model.StudentExam{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *authservice) GetStudentStudyHistory(ctx context.Context, id int) (response.StudentStudyHistory, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var data response.StudentStudyHistory
	res := s.db.WithContext(ctx).Table("user u").
		Joins("LEFT JOIN role r ON r.id = u.role_id").
		Joins("LEFT JOIN student_profile sp ON sp.user_id = u.id").
		Joins("LEFT JOIN student_category sc ON sc.id = sp.student_category_id").
		Joins("LEFT JOIN villages v ON v.id = sp.village_id").
		Joins("LEFT JOIN communes c ON c.id = v.commune_id").
		Joins("LEFT JOIN districts d ON d.id = c.district_id").
		Joins("LEFT JOIN provinces p ON p.id = d.province_id").
		Where("u.id = ?", id).
		Select(`
			u.id AS id,
			u.name_kh AS name_kh,
			u.name_en AS name_en,
			u.gender AS gender,
			u.code AS code,
			r.display_name AS role_name,
			sc.name AS student_category_name,
			sp.date_of_birth AS date_of_birth,
			sp.nationality AS nationality,
			sp.phone AS phone,
			CONCAT_WS(' - ', p.name_kh, p.name_en) AS province,
			CONCAT_WS(' - ', d.name_kh, d.name_en) AS district,
			CONCAT_WS(' - ', c.name_kh, c.name_en) AS communce,
			CONCAT_WS(' - ', v.name_kh, v.name_en) AS village,
			sp.occupation AS occupation,
			sp.academic_stream AS academic_stream,
			sp.telegram_username AS telegram_username,
			sp.exam_in AS exam_in,
			sp.exam_out AS exam_out
		`).Scan(&data)
	if res.Error != nil {
		return response.StudentStudyHistory{}, res.Error
	}
	if res.RowsAffected == 0 {
		return response.StudentStudyHistory{}, gorm.ErrRecordNotFound
	}

	data.DateOfBirth = helper.FormatDate(data.DateOfBirth)

	userClasses := make([]response.UserClassRes, 0)
	if err := s.db.WithContext(ctx).
		Table("user_class uc").
		Joins("LEFT JOIN class c ON c.id = uc.class_id").
		Joins("LEFT JOIN major m ON m.id = c.major_id").
		Joins("LEFT JOIN shift s ON s.id = c.shift_id").
		Joins("LEFT JOIN generation g ON g.id = c.generation_id").
		Joins("LEFT JOIN programmes p ON p.id = c.programme_id").
		Where("uc.user_id = ?", id).
		Order("c.year ASC, c.semester ASC").
		Select(`
			uc.id AS user_class_id,
			uc.user_id AS user_id,
			c.name AS class_name,
			c.type AS type,
			c.is_active AS active,
			m.id AS major_id,
			CONCAT_WS(' ', m.name_kh, m.name_en) AS major_name,
			s.id AS shift_id,
			s.name AS shift_name,
			g.id AS generation_id,
			g.code AS generation_name,
			g.start_year AS generation_start,
			g.end_year AS generation_end,
			c.year AS year,
			c.semester AS semester,
			c.` + "`group`" + ` AS ` + "`group_name`" + `,
			c.term AS term,
			p.id AS programme_id,
			p.name AS programme_name
		`).Scan(&userClasses).Error; err != nil {
		return response.StudentStudyHistory{}, err
	}

	for i := range userClasses {
		uc := &userClasses[i]
		uc.SubjectRes = make([]response.SubjectRes, 0)
		if uc.MajorID == nil || uc.GenerationID == nil || uc.ProgrammeID == nil || uc.Year == nil {
			continue
		}

		var subjects []response.SubjectRes
		if err := s.db.WithContext(ctx).
			Table("major_subject ms").
			Joins("LEFT JOIN subject s ON s.id = ms.subject_id").
			Where(`ms.major_id = ? AND ms.generation_id = ? AND ms.programme_id = ?
			       AND ms.year = ? AND ms.semester = ?`,
				*uc.MajorID, *uc.GenerationID, *uc.ProgrammeID, *uc.Year, uc.Semester).
			Select(`
				s.id AS subject_id,
				s.code AS subject_code,
				s.name_kh AS subject_name_kh,
				s.name_en AS subject_name_en,
				s.credit_hour AS credit_hour
			`).Scan(&subjects).Error; err != nil {
			return response.StudentStudyHistory{}, err
		}

		for j := range subjects {
			subjects[j].UserClassID = uc.UserClassID
		}
		if subjects != nil {
			uc.SubjectRes = subjects
		}
	}

	data.UserClassRes = userClasses
	return data, nil
}
