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

type SubjectService interface {
	Create(ctx context.Context, input request.SubjectRequestCreate) error
	GetWithPagination(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.Subject, *model.PaginationMetadata, error)
	Update(ctx context.Context, id int, input request.SubjectRequestUpdate) error
	Toggle(ctx context.Context, id int) error
}

type subjectService struct {
	db *gorm.DB
}

func NewSubjectService() SubjectService {
	return &subjectService{
		db: config.DB,
	}
}

func (s *subjectService) Create(ctx context.Context, input request.SubjectRequestCreate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newdata := model.Subject{
			Code:       input.Code,
			NameKh:     input.NameKh,
			NameEn:     input.NameEn,
			CreditHour: input.CreditHour,
			IsActive:   true,
		}
		if err := tx.Create(&newdata).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *subjectService) GetWithPagination(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.Subject, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)
	var data []response.Subject
	var total int64
	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("subject s")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["name"]; ok && v != "" {
			tx = tx.Where(
				"s.name_kh LIKE ? OR s.name_en LIKE ?",
				"%"+v+"%",
				"%"+v+"%",
			)
		}
		return tx
	}

	if err := applyFilters(base()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count customer: %w", err)
	}

	if total == 0 {
		return []response.Subject{}, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize

	dataQuery := applyFilters(base()).Select(`
		s.id AS id,
		s.code AS code,
		s.name_kh AS name_kh,
		s.name_en AS name_en,
		s.credit_hour AS credit_hour,
		s.is_active AS is_active
	`)
	if err := dataQuery.Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch companies: %w", err)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *subjectService) Update(ctx context.Context, id int, input request.SubjectRequestUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var data model.Subject
		if err := tx.Where("id = ?", id).First(&data).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "classcurriculumn not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch classcurriculumn", nil)
		}
		data.Code = input.Code
		data.NameKh = input.NameKh
		data.NameEn = input.NameEn
		data.CreditHour = input.CreditHour
		if err := tx.Save(&data).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update student", nil)
		}
		return nil
	})
	return err
}

func (s *subjectService) Toggle(ctx context.Context, id int) error {
	return utils.ToggleStatus[model.Subject](ctx, s.db, id)
}
