package service

import (
	"context"
	"mysql/config"
	"mysql/response"

	"gorm.io/gorm"
)

type TranscriptService interface {
	Student(ctx context.Context, userID int) (*response.Student, error)
	LatestProgramme(ctx context.Context, userID int) (int, error)
	Header(ctx context.Context, userID, programmeID int) (*response.Header, error)
	Transcript(
		ctx context.Context,
		userID int,
	) (*response.Transcript, error)
}

type transcriptservice struct {
	db *gorm.DB
}

func NewTranscriptService() TranscriptService {
	return &transcriptservice{
		db: config.DB,
	}
}

func (s *transcriptservice) Student(ctx context.Context, userID int) (*response.Student, error) {
	var st response.Student

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("user u").
			Joins("LEFT JOIN student_profile p ON p.user_id = u.id").
			Where("u.id = ?", userID)
	}

	query := base().Select(`
		u.id AS id,
		u.code AS code,
		COALESCE(u.name_kh, '') AS name_kh,
		COALESCE(u.name_en, '') AS name_en,
		COALESCE(u.gender, 0) AS gender,
		p.date_of_birth AS date_of_birth,
		COALESCE(p.nationality, '') AS nationality
	`)

	if err := query.Scan(&st).Error; err != nil {
		return nil, err
	}

	return &st, nil
}

func (s *transcriptservice) LatestProgramme(ctx context.Context, userID int) (int, error) {
	var programmeID int

	err := s.db.WithContext(ctx).
		Table("score s").
		Select("s.programme_id").
		Where("s.user_id = ?", userID).
		Order("s.id DESC").
		Limit(1).
		Scan(&programmeID).Error

	return programmeID, err
}

func (s *transcriptservice) Header(ctx context.Context, userID, programmeID int) (*response.Header, error) {
	var hd response.Header

	err := s.db.WithContext(ctx).
		Table("score s").
		Joins("JOIN major m ON m.id = s.major_id").
		Joins("LEFT JOIN faculties f ON f.id = m.faculty_id").
		Joins("JOIN programmes p ON p.id = s.programme_id").
		Joins("JOIN generation g ON g.id = s.generation_id").
		Where(
			"s.user_id = ? AND s.programme_id = ? AND s.is_active = 1",
			userID,
			programmeID,
		).
		Select(`
			COALESCE(f.name, '') AS faculty_name,
			COALESCE(m.name_en, '') AS major_name,
			COALESCE(p.degree_title_en, '') AS degree_title,
			COALESCE(g.start_year, 0) AS start_year,
			COALESCE(g.end_year, 0) AS end_year
		`).
		Order("s.id DESC").Limit(1).
		Scan(&hd).Error

	if err != nil {
		return nil, err
	}

	return &hd, nil
}

func (s *transcriptservice) Results(
	ctx context.Context,
	userID, programmeID int,
) ([]response.RawResult, error) {

	var out []response.RawResult

	err := s.db.WithContext(ctx).
		Table("score s").
		Joins("JOIN subject sub ON sub.id = s.subject_id").
		Joins("LEFT JOIN score_detail sd ON sd.score_id = s.id").
		Joins("LEFT JOIN grade_components gc ON gc.id = sd.grade_component_id AND gc.active = 1").
		Where(
			"s.user_id = ? AND s.programme_id = ? AND s.is_active = 1",
			userID,
			programmeID,
		).
		Select(`
			s.id AS score_id,
			s.year AS year,
			s.semester AS semester,
			sub.code AS code,
			COALESCE(sub.name_en, '') AS name_en,
			COALESCE(sub.name_kh, '') AS name_kh,
			COALESCE(sub.credit_hour, 0) AS credits,
			COALESCE(
				SUM(CASE WHEN gc.id IS NULL THEN 0 ELSE sd.score END),
				0
			) AS points_sum,
			COALESCE(
				SUM(sd.score * gc.weight_percentage / 100),
				0
			) AS percent_sum
		`).
		Group(`
			s.id,
			s.year,
			s.semester,
			sub.code,
			sub.name_en,
			sub.name_kh,
			sub.credit_hour
		`).
		Order("s.year, s.semester, sub.code").
		Scan(&out).Error

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *transcriptservice) Transcript(
	ctx context.Context,
	userID int,
) (*response.Transcript, error) {

	programmeID, err := s.LatestProgramme(ctx, userID)
	if err != nil {
		return nil, err
	}

	st, err := s.Student(ctx, userID)
	if err != nil {
		return nil, err
	}

	h, err := s.Header(ctx, userID, programmeID)
	if err != nil {
		return nil, err
	}

	raw, err := s.Results(ctx, userID, programmeID)
	if err != nil {
		return nil, err
	}

	return &response.Transcript{
		Student: *st,
		Header:  *h,
		Results: raw,
	}, nil
}
