package service

import (
	"context"
	"errors"
	"fmt"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/model"
	"mysql/request"
	"mysql/response"
	"mysql/utils"

	"gorm.io/gorm"
)

type ClassScheduleService interface {
	// GetByClass returns every schedule row for a class (active and
	// inactive), joined with subject for display.
	GetByClass(ctx context.Context, classID int) ([]response.ClassScheduleResponse, error)

	// GetAvailableSubjects returns the subjects that are allowed to be
	// scheduled for this class: subjects attached (active) to the class's
	// major for the class's own year/semester.
	GetAvailableSubjects(ctx context.Context, classID int) ([]response.SubjectOption, error)

	Create(ctx context.Context, classID int, input request.ClassScheduleRequestCreate) error
	Toggle(ctx context.Context, id int) error
}

type classScheduleService struct {
	db *gorm.DB
}

func NewClassScheduleService() ClassScheduleService {
	return &classScheduleService{
		db: config.DB,
	}
}

func (s *classScheduleService) GetByClass(ctx context.Context, classID int) ([]response.ClassScheduleResponse, error) {
	var data []response.ClassScheduleResponse

	err := s.db.WithContext(ctx).
		Table("class_schedule cs").
		Joins("JOIN subject sub ON sub.id = cs.subject_id").
		Where("cs.class_id = ?", classID).
		Select(`
			cs.id AS id,
			cs.class_id AS class_id,
			cs.subject_id AS subject_id,
			sub.code AS subject_code,
			sub.name_kh AS subject_name,
			cs.day_of_week AS day_of_week,
			cs.is_active AS is_active
		`).
		Order("cs.id DESC").
		Scan(&data).Error
	if err != nil {
		return nil, fmt.Errorf("fetch class schedule: %w", err)
	}

	return data, nil
}

func (s *classScheduleService) GetAvailableSubjects(ctx context.Context, classID int) ([]response.SubjectOption, error) {
	var class model.Class
	if err := s.db.WithContext(ctx).Where("id = ?", classID).First(&class).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.CodeNotFound, "class not found", nil)
		}
		return nil, apperror.New(apperror.CodeInternal, "failed to fetch class", nil)
	}

	var data []response.SubjectOption
	err := s.db.WithContext(ctx).
		Table("major_subject ms").
		Joins("JOIN subject sub ON sub.id = ms.subject_id").
		Where("ms.major_id = ? AND ms.year = ? AND ms.semester = ? AND ms.is_active = ?",
			class.MajorID, class.Year, class.Semester, true).
		Select(`
			sub.id AS id,
			sub.code AS code,
			sub.name_kh AS name_kh,
			sub.name_en AS name_en
		`).
		Order("sub.id ASC").
		Scan(&data).Error
	if err != nil {
		return nil, fmt.Errorf("fetch available subjects: %w", err)
	}

	return data, nil
}

func (s *classScheduleService) Create(ctx context.Context, classID int, input request.ClassScheduleRequestCreate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var class model.Class
		if err := tx.Where("id = ?", classID).First(&class).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "class not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch class", nil)
		}

		// The subject must actually belong to this class's major for this
		// class's year/semester (mirrors GetAvailableSubjects, enforced
		// again server-side so the frontend list can't be bypassed).
		var count int64
		if err := tx.Model(&model.MajorSubject{}).
			Where("major_id = ? AND subject_id = ? AND year = ? AND semester = ? AND is_active = ?",
				class.MajorID, input.SubjectID, class.Year, class.Semester, true).
			Count(&count).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to validate subject", nil)
		}
		if count == 0 {
			return apperror.New(apperror.CodeInvalidInput, "subject is not part of this class's curriculum for its year/semester", nil)
		}

		// Prevent scheduling the same subject twice on the same day for
		// this class.
		var dup int64
		if err := tx.Model(&model.ClassSchedule{}).
			Where("class_id = ? AND subject_id = ? AND day_of_week = ?",
				classID, input.SubjectID, input.DayOfWeek).
			Count(&dup).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to check existing schedule", nil)
		}
		if dup > 0 {
			return apperror.New(apperror.CodeConflict, "this subject is already scheduled on that day for this class", nil)
		}

		newdata := model.ClassSchedule{
			ClassID:   int64(classID),
			SubjectID: input.SubjectID,
			DayOfWeek: input.DayOfWeek,
			IsActive:  true,
		}
		if err := tx.Create(&newdata).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create class schedule", nil)
		}
		return nil
	})
}

func (s *classScheduleService) Toggle(ctx context.Context, id int) error {
	return utils.ToggleStatus[model.ClassSchedule](ctx, s.db, id)
}
