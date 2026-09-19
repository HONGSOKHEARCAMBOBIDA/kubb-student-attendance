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

type MajorService interface {
	Create(ctx context.Context, input request.MajorRequestCreate) error
	GetWithPagination(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.Major, *model.PaginationMetadata, error)
	Update(ctx context.Context, id int, input request.MajorRequestUpdate) error
	Toggle(ctx context.Context, id int) error

	// Subjects belonging to a major (major_subject junction)
	AddSubject(ctx context.Context, majorID int, input request.MajorSubjectRequestCreate) error
	GetSubjects(ctx context.Context, majorID int, pf request.Pagination) ([]response.MajorSubject, *model.PaginationMetadata, error)
	UpdateSubject(ctx context.Context, majorSubjectID int, input request.MajorSubjectRequestUpdate) error
	ToggleSubject(ctx context.Context, majorSubjectID int) error
	RemoveSubject(ctx context.Context, majorSubjectID int) error
}

type majorService struct {
	db *gorm.DB
}

func NewMajorService() MajorService {
	return &majorService{
		db: config.DB,
	}
}

func (s *majorService) Create(ctx context.Context, input request.MajorRequestCreate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newdata := model.Major{
			Code:     input.Code,
			NameKh:   input.NameKh,
			NameEn:   input.NameEn,
			IsActive: true,
		}
		if err := tx.Create(&newdata).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *majorService) GetWithPagination(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.Major, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)
	var data []response.Major
	var total int64
	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("major m")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["name"]; ok && v != "" {
			tx = tx.Where(
				"m.name_kh LIKE ? OR m.name_en LIKE ?",
				"%"+v+"%",
				"%"+v+"%",
			)
		}
		return tx
	}

	if err := applyFilters(base()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count major: %w", err)
	}

	if total == 0 {
		return []response.Major{}, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize

	dataQuery := applyFilters(base()).Select(`
		m.id AS id,
		m.code AS code,
		m.name_kh AS name_kh,
		m.name_en AS name_en,
		m.is_active AS is_active
	`)
	if err := dataQuery.Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch major: %w", err)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *majorService) Update(ctx context.Context, id int, input request.MajorRequestUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var data model.Major
		if err := tx.Where("id = ?", id).First(&data).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "major not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch major", nil)
		}
		data.Code = input.Code
		data.NameKh = input.NameKh
		data.NameEn = input.NameEn
		if err := tx.Save(&data).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update major", nil)
		}
		return nil
	})
	return err
}

func (s *majorService) Toggle(ctx context.Context, id int) error {
	return utils.ToggleStatus[model.Major](ctx, s.db, id)
}

// AddSubject attaches a subject to a major for a given year/semester.
// It validates that both the major and the subject exist, and blocks
// adding the exact same subject/year/semester combo twice.
func (s *majorService) AddSubject(ctx context.Context, majorID int, input request.MajorSubjectRequestCreate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var major model.Major
		if err := tx.Where("id = ?", majorID).First(&major).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "major not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch major", nil)
		}

		var subject model.Subject
		if err := tx.Where("id = ?", input.SubjectID).First(&subject).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "subject not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch subject", nil)
		}

		var count int64
		if err := tx.Model(&model.MajorSubject{}).
			Where("major_id = ? AND subject_id = ? AND year = ? AND semester = ?",
				majorID, input.SubjectID, input.Year, input.Semester).
			Count(&count).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to check existing subject", nil)
		}
		if count > 0 {
			// If your apperror package doesn't have CodeConflict yet, swap
			// this for CodeInternal/CodeBadRequest, whichever fits your codes.
			return apperror.New(apperror.CodeConflict, "subject already added to this major for that year/semester", nil)
		}

		newdata := model.MajorSubject{
			MajorID:   majorID,
			SubjectID: input.SubjectID,
			Year:      input.Year,
			Semester:  input.Semester,
			IsActive:  true,
		}
		if err := tx.Create(&newdata).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to add subject to major", nil)
		}
		return nil
	})
	return err
}

// GetSubjects lists the subjects attached to a major, joined with
// subject so the UI can show code/name/credit hour directly.
func (s *majorService) GetSubjects(ctx context.Context, majorID int, pf request.Pagination) ([]response.MajorSubject, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)
	var data []response.MajorSubject
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("major_subject ms").
			Joins("JOIN subject s ON s.id = ms.subject_id").
			Where("ms.major_id = ?", majorID)
	}

	if err := base().Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count major subject: %w", err)
	}

	if total == 0 {
		return []response.MajorSubject{}, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize

	dataQuery := base().Select(`
		ms.id AS id,
		ms.major_id AS major_id,
		ms.subject_id AS subject_id,
		s.code AS subject_code,
		s.name_kh AS subject_name_kh,
		s.name_en AS subject_name_en,
		s.credit_hour AS credit_hour,
		ms.year AS year,
		ms.semester AS semester,
		ms.is_active AS is_active
	`)
	if err := dataQuery.Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch major subject: %w", err)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *majorService) UpdateSubject(ctx context.Context, majorSubjectID int, input request.MajorSubjectRequestUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var data model.MajorSubject
		if err := tx.Where("id = ?", majorSubjectID).First(&data).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "major subject not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch major subject", nil)
		}
		data.Year = input.Year
		data.Semester = input.Semester
		if err := tx.Save(&data).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update major subject", nil)
		}
		return nil
	})
	return err
}

func (s *majorService) ToggleSubject(ctx context.Context, majorSubjectID int) error {
	return utils.ToggleStatus[model.MajorSubject](ctx, s.db, majorSubjectID)
}

func (s *majorService) RemoveSubject(ctx context.Context, majorSubjectID int) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Where("id = ?", majorSubjectID).Delete(&model.MajorSubject{})
		if res.Error != nil {
			return apperror.New(apperror.CodeInternal, "failed to remove subject from major", nil)
		}
		if res.RowsAffected == 0 {
			return apperror.New(apperror.CodeNotFound, "major subject not found", nil)
		}
		return nil
	})
}
