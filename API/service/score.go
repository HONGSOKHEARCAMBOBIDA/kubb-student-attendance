package service

import (
	"context"
	"mysql/config"
	"mysql/model"
	"mysql/request"
	"mysql/utils"

	"gorm.io/gorm"
)

type ScoreService interface {
	GetGradeComponent(ctx context.Context) ([]model.GradeComponent, error)
	CreateScore(ctx context.Context, input request.CreateScoreRequest) error
}
type scoreservice struct {
	db *gorm.DB
}

func NewScoreService() ScoreService {
	return &scoreservice{
		db: config.DB,
	}
}

func (s *scoreservice) GetGradeComponent(ctx context.Context) ([]model.GradeComponent, error) {
	var data []model.GradeComponent

	if err := s.db.WithContext(ctx).
		Order("id ASC").
		Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (s *scoreservice) CreateScore(ctx context.Context, input request.CreateScoreRequest) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, student := range input.Students {
			score := model.Score{
				ClassID:      input.ClassID,
				UserID:       student.UserID,
				MajorID:      input.MajorID,
				GenerationID: input.GenerationID,
				ProgrammeID:  input.ProgrammeID,
				SubjectID:    input.SubjectID,
				Year:         input.Year,
				Semester:     input.Semester,
				IsActive:     true,
			}

			if err := tx.Create(&score).Error; err != nil {
				return err
			}

			for _, detail := range student.Details {
				scoreDetail := model.ScoreDetail{
					ScoreID:          int64(score.ID),
					GradeComponentID: detail.GradeComponentID,
					Score:            detail.Score,
				}

				if err := tx.Create(&scoreDetail).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}
