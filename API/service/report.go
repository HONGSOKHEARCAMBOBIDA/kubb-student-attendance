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
	GetAcademicConfirmation(ctx context.Context, userclassID int) (response.AcademicConfirmation, error)
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

	if err := Query.Offset(offset).Limit(pf.PageSize).Order("id").Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch data: %w", err)
	}
	for i := range data {
		data[i].ScheduleDate = helper.FormatDate(data[i].ScheduleDate)
	}
	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *reportservice) GetAcademicConfirmation(ctx context.Context, userclassID int) (response.AcademicConfirmation, error) {
	var data response.AcademicConfirmation
	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("user_class uc").
			Joins("LEFT JOIN user u ON u.id = uc.user_id").
			Joins("LEFT JOIN class c ON c.id = uc.class_id").
			Joins("LEFT JOIN student_profile sp ON sp.user_id = uc.user_id").
			Joins("LEFT JOIN villages v ON v.id = sp.village_id").
			Joins("LEFT JOIN communes cm ON cm.id = v.commune_id").
			Joins("LEFT JOIN districts d ON d.id = cm.district_id").
			Joins("LEFT JOIN provinces p ON p.id = d.province_id").
			Joins("LEFT JOIN generation g ON g.id = c.generation_id").
			Joins("LEFT JOIN major m ON m.id = c.major_id").
			Joins("LEFT JOIN programmes pm ON pm.id = c.programme_id").
			Joins("LEFT JOIN faculties f ON f.id = m.faculty_id").
			Where("uc.id =?", userclassID)
	}
	query := base().Select(`
		u.name_kh AS name_kh,
		u.name_en AS name_en,
		u.gender AS gender,
		sp.nationality AS national_kh,
		u.code AS code_en,
		sp.date_of_birth AS dob_en,
		v.name_kh AS village_kh,
		v.name_en AS village_en,
		cm.name_kh AS communce_kh,
		cm.name_en AS communce_en,
		d.name_kh AS district_kh,
		d.name_en AS district_en,
		p.name_kh AS province_kh,
		p.name_en AS province_en,
		pm.name AS programme_kh,
		pm.degree_title_en AS programme_en,
		c.year AS year_en,
		c.semester AS semester_en,
		f.name AS faculty_kh,
		f.name_en AS faculty_en,
		m.name_kh AS major_kh,
		m.name_en AS major_en,
		g.start_year AS start_en,
		g.end_year AS end_en
	`)

	if err := query.Scan(&data).Error; err != nil {
		return response.AcademicConfirmation{}, err
	}
	data.DobEn = helper.FormatDate(data.DobEn)
	gender_kh, gender_en := helper.GenderKhftEn(data.Gender)
	data.GenderEn = gender_en
	data.GenderKh = gender_kh
	data.NationEn = helper.National(data.Nationalkh)
	data.Codekh = helper.ConvertToKhmerNumber(data.CodeEN)
	data.Dobkh = helper.ConvertDateToKhmer(data.DobEn)
	data.Yearkh = helper.ConvertToKhmerNumber(data.YearEn)
	data.Semesterkh = helper.ConvertToKhmerNumber(data.SemesterEn)
	data.StartKh = helper.ConvertToKhmerNumber(data.StartEn)
	data.Endkh = helper.ConvertToKhmerNumber(data.EndEn)

	return data, nil

}
