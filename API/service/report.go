package service

import (
	"context"
	"fmt"
	"mysql/config"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/response"

	"gorm.io/gorm"
)

type ReportService interface {
	GetWeekSchedule(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.WeekSchedule, *model.PaginationMetadata, error)
}

type reportservice struct {
	db *gorm.DB
}

func NewReportService() ReportService {
	return &reportservice{
		db: config.DB,
	}
}

func (s *reportservice) GetWeekSchedule(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.WeekSchedule, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)
	var data []response.WeekSchedule
	var total int64
	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("class_schedule cs").
			Joins("LEFT JOIN class c ON c.id = cs.class_id").
			Joins("LEFT JOIN subject s ON s.id = cs.subject_id").
			Joins("LEFT JOIN programmes p ON p.id = c.programme_id").
			Joins("LEFT JOIN shift sf ON sf.id = c.shift_id").
			Joins("LEFT JOIN generation g ON g.id = c.generation_id").
			Joins("LEFT JOIN major m ON m.id = c.major_id").
			Joins("LEFT JOIN teacher_rates tr ON tr.class_schedule_id = cs.id").
			Joins("LEFT JOIN user u ON u.id = tr.user_id").
			Joins("LEFT JOIN teacher_profile tp ON tp.user_id = u.id")
	}

	Filters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["programme_id"]; ok && v != "" {
			tx = tx.Where("c.programme_id = ?", v)
		}
		if v, ok := filter["generation_id"]; ok && v != "" {
			tx = tx.Where("c.generation_id = ?", v)
		}
		if v, ok := filter["start_date"]; ok && v != "" {
			tx = tx.Where("cs.schedule_date >= ?", v)
		}
		if v, ok := filter["end_date"]; ok && v != "" {
			tx = tx.Where("cs.schedule_date <= ?", v)
		}
		return tx
	}
	if err := Filters(base()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count : %w", err)
	}

	if total == 0 {
		return []response.WeekSchedule{}, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize
	Query := Filters(base()).Select(`
		p.name AS programme_name,
		sf.name AS shift_name,
		g.name_kh AS generation_name,
		m.name_kh AS major_name,
		s.name_kh AS subject_name_kh,
		u.name_kh AS name_kh,
		u.name_en AS name_en,
		u.gender AS gender,
		tp.phone AS phone,
		cs.room AS room,
		cs.schedule_date AS schedule_date,
		cs.id AS id,
		cs.total_student AS total_student,
		cs.total_session AS total_session,
		cs.is_mission AS is_mission,
		cs.day_of_week AS day_of_week
	`)

	if err := Query.Offset(offset).Limit(pf.PageSize).Order("id DESC").Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch data: %w", err)
	}
	for i := range data {
		data[i].ScheduleDate = helper.FormatDate(data[i].ScheduleDate)
	}
	return data, helper.BuildPaginationMeta(pf, total), nil
}
