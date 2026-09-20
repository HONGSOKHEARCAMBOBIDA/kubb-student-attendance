package service

import (
	"context"
	"errors"
	"fmt"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/response"
	"mysql/utils"

	"gorm.io/gorm"
)

type CompanyService interface {
	//GetCompanyColor(userID int) (response.CompanyColor, error)
	GetClass(id int, ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.ClassResponse, *model.PaginationMetadata, error)
	GetClassScan(ctx context.Context, id int) ([]response.ClassScanResponse, error)
	CreateClass(ctx context.Context, input request.ClassRequestCreate) error
	UpdateClass(ctx context.Context, id int, input request.ClassRequestUpdate) error
	ChangeStatusClass(ctx context.Context, id int) error
	UpdateTelegram(ctx context.Context, id int, input request.CompanyRequestUpdateTelegram) error
	ShowManageCompany(ctx context.Context, id int) ([]helper.ManageCompany, error)
	GetMajor(ctx context.Context) ([]model.Major, error)
	GetShift(ctx context.Context) ([]model.Shift, error)
	GetGeneration(ctx context.Context) ([]model.Generation, error)
	GetProgramme(ctx context.Context) ([]model.Programme, error)
}

type companyservice struct {
	db *gorm.DB
}

func NewCompanyService() CompanyService {
	return &companyservice{
		db: config.DB,
	}
}

func (s *companyservice) GetMajor(ctx context.Context) ([]model.Major, error) {
	var data []model.Major

	if err := s.db.WithContext(ctx).
		Where("is_active = ?", 1).
		Order("id ASC").
		Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (s *companyservice) GetShift(ctx context.Context) ([]model.Shift, error) {
	var data []model.Shift

	if err := s.db.WithContext(ctx).
		Order("id ASC").
		Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (s *companyservice) GetGeneration(ctx context.Context) ([]model.Generation, error) {
	var data []model.Generation

	if err := s.db.WithContext(ctx).
		Where("is_active = ?", 1).
		Order("id ASC").
		Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (s *companyservice) GetProgramme(ctx context.Context) ([]model.Programme, error) {
	var data []model.Programme

	if err := s.db.WithContext(ctx).
		Order("id ASC").
		Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

// func (s *companyservice) GetCompanyColor(userID int) (response.CompanyColor, error) {
// 	var color response.CompanyColor

// 	err := s.db.Table("company AS c").
// 		Select(`c.color AS color`).
// 		Joins("LEFT JOIN user u ON u.company_id = c.id").
// 		Where("u.id = ?", userID).
// 		Scan(&color).Error

// 	if err != nil {
// 		return color, err
// 	}

// 	return color, nil
// }

func (s *companyservice) GetClass(id int, ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.ClassResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)
	var data []response.ClassResponse
	var total int64
	var user model.User
	if err := s.db.WithContext(ctx).Preload("Role").First(&user, id).Error; err != nil {
		return nil, nil, err
	}
	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("class c").
			Joins("LEFT JOIN major AS m ON m.id = c.major_id").
			Joins("LEFT JOIN shift AS sh ON sh.id = c.shift_id").
			Joins("LEFT JOIN generation AS g ON g.id = c.generation_id").
			Joins("LEFT JOIN programmes AS p ON p.id = c.programme_id")
	}
	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["name"]; ok && v != "" {
			tx = tx.Where("c.name LIKE ?", "%"+v+"%")
		}
		if v, ok := filter["major_id"]; ok && v != "" {
			tx = tx.Where("c.major_id = ?", v)
		}
		if v, ok := filter["shift_id"]; ok && v != "" {
			tx = tx.Where("c.shift_id = ?", v)
		}
		if v, ok := filter["generation_id"]; ok && v != "" {
			tx = tx.Where("c.generation_id = ?", v)
		}
		if v, ok := filter["programme_id"]; ok && v != "" {
			tx = tx.Where("c.programme_id = ?", v)
		}
		return tx
	}
	if err := applyFilters(base()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count product: %w", err)
	}

	if total == 0 {
		return []response.ClassResponse{}, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize
	dataQuery := applyFilters(base()).Select(`
		c.id AS id,
		c.name AS name,
		c.is_active AS is_active,
		c.latitude AS latitude,
		c.longitude AS longitude,
		c.radius AS radius,
		c.bot_token AS bot_token,
		c.group_chatID AS group_chatID,
		c.can_scan_outsize AS can_scan_outsize,
		c.major_id AS major_id,
		m.name_kh AS major_name,
		c.shift_id AS shift_id,
		sh.name AS shift_name,
		c.generation_id AS generation_id,
		g.name_kh AS generation_name,
		c.year AS year,
		c.semester AS semester,
		c.` + "`group`" + ` AS ` + "`group`" + `,
		c.term AS term,
		c.programme_id AS programme_id,
		p.name AS programme_name
	`)
	if err := dataQuery.Offset(offset).Limit(pf.PageSize).Order("id DESC").Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch class: %w", err)
	}

	classIDs := make([]int, len(data))
	for i, c := range data {
		classIDs[i] = c.ID
	}

	var students []response.StudentWithClass
	if err := s.db.WithContext(ctx).Table("user u").
		Joins("INNER JOIN user_class uc ON uc.user_id = u.id").
		Where("uc.class_id IN ?", classIDs).
		Select(`
		u.id AS id,
		u.name_kh AS name_kh,
		u.name_en AS name_en,
		u.gender AS gender,
		u.code AS code,
		uc.class_id AS class_id
	`).Scan(&students).Error; err != nil {
		return nil, nil, err
	}

	studentByClass := make(map[int][]response.UserResponse, len(data))
	for _, st := range students {
		studentByClass[st.ClassID] = append(studentByClass[st.ClassID], st.UserResponse)
	}

	for i := range data {
		data[i].UserResponse = studentByClass[data[i].ID]
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *companyservice) GetClassScan(ctx context.Context, id int) ([]response.ClassScanResponse, error) {
	var user model.User
	if err := s.db.WithContext(ctx).
		Select("id", "role_id").
		Preload("Role").
		First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to load user: %w", err)
	}

	var class []response.ClassScanResponse
	query := s.db.WithContext(ctx).Table("class AS c").
		Select(`
			c.id AS id,
			c.name AS name
		`)
	query = helper.ManageClassFilter(query, s.db, user)

	if err := query.Scan(&class).Error; err != nil {
		return nil, fmt.Errorf("failed to scan company data: %w", err)
	}
	return class, nil
}

func (s *companyservice) CreateClass(ctx context.Context, input request.ClassRequestCreate) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	lat, lng, err := utils.ExtractLatLngFromGoogleMapsURL(input.MapLink)
	if err != nil {
		return fmt.Errorf("invalid map_link: %w", err)
	}
	newClass := model.Class{
		Name:           &input.Name,
		IsActive:       true,
		Latitude:       &lat,
		Longitude:      &lng,
		Radius:         input.Radius,
		BotToken:       nil,
		GroupChatID:    nil,
		CanScanOutsize: input.CanScanOutsize,
		MajorID:        input.MajorID,
		ShiftID:        input.ShiftID,
		GenerationID:   input.GenerationID,
		Year:           input.Year,
		Semester:       input.Semester,
		Group:          input.Group,
		Term:           input.Term,
		ProgrammeID:    input.ProgrammeID,
	}

	if err := tx.WithContext(ctx).
		Create(&newClass).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *companyservice) UpdateClass(ctx context.Context, id int, input request.ClassRequestUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var data model.Class
		if err := tx.Where("id = ?", id).First(&data).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "classcurriculumn not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch classcurriculumn", nil)
		}
		data.Name = &input.Name
		data.MajorID = input.MajorID
		data.ShiftID = input.ShiftID
		data.GenerationID = input.GenerationID
		data.Year = input.Year
		data.Semester = input.Semester
		data.Group = input.Group
		data.Term = input.Term
		data.ProgrammeID = input.ProgrammeID
		data.CanScanOutsize = input.CanScanOutsize
		if err := tx.Save(&data).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update product", nil)
		}
		return nil
	})
	return err
}

func (s *companyservice) ChangeStatusClass(ctx context.Context, id int) error {
	result := s.db.WithContext(ctx).Model(&model.Class{}).Where("id =?", id).Update("is_active", gorm.Expr("NOT is_active"))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *companyservice) UpdateTelegram(ctx context.Context, id int, input request.CompanyRequestUpdateTelegram) error {
	updates := map[string]interface{}{}
	if input.BotToken != nil {
		encryptedBottoken, err := utils.EncryptBotToken(*input.BotToken)
		if err != nil {
			return err
		}
		updates["bot_token"] = encryptedBottoken
	}
	if input.GroupLink != nil {
		chatID, err := utils.ResolveTelegramChatID(*input.BotToken, *input.GroupLink)
		if err != nil {
			return fmt.Errorf("could not resolve group link")
		}
		chatIDStr := fmt.Sprintf("%d", chatID)
		encryptedChatID, err := utils.EncryptChatID(chatIDStr)
		if err != nil {
			return err
		}
		updates["group_chatID"] = encryptedChatID
	}
	if len(updates) == 0 {
		return errors.New(" no field to update")
	}
	result := s.db.WithContext(ctx).Model(&model.Company{}).Where("id =?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *companyservice) ShowManageCompany(ctx context.Context, id int) ([]helper.ManageCompany, error) {
	var user model.User

	if err := s.db.WithContext(ctx).
		Preload("Role").
		Where("id = ?", id).
		First(&user).Error; err != nil {
		return nil, err
	}

	managecompany := helper.ShowManageCompany(user)

	return managecompany, nil
}
