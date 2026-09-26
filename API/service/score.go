package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"mysql/config"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/response"
	"mysql/utils"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ScoreService interface {
	GetGradeComponent(ctx context.Context) ([]model.GradeComponent, error)
	CreateScore(ctx context.Context, input request.CreateScoreRequest) error
	ImportScoreFromExcell(ctx context.Context, req request.ImportScoreExcelRequest, file multipart.File) (*request.ImportScoreResult, error)
	GetScore(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.ScoreResponse, *model.PaginationMetadata, error)
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

const excelCodeHeaderKH = "អត្តលេខ"

func (s *scoreservice) ImportScoreFromExcell(ctx context.Context, req request.ImportScoreExcelRequest, file multipart.File) (*request.ImportScoreResult, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("invalid excel file: %w", err)
	}
	defer f.Close()

	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil || len(rows) < 2 {
		return nil, fmt.Errorf("sheet is empty")
	}

	var components []model.GradeComponent
	if err := s.db.WithContext(ctx).Find(&components).Error; err != nil {
		return nil, err
	}

	header := rows[0]
	codeCol := -1
	componentCols := make(map[int]uint64)

	for i, cell := range header {
		h := strings.TrimSpace(cell)
		if h == excelCodeHeaderKH || strings.EqualFold(h, "code") {
			codeCol = i
			continue
		}
		for _, gc := range components {
			if strings.EqualFold(strings.TrimSpace(gc.Name), h) {
				componentCols[i] = gc.ID
				break
			}
		}
	}

	if codeCol == -1 {
		return nil, fmt.Errorf("column '%s' not found in header", excelCodeHeaderKH)
	}

	if len(componentCols) == 0 {
		return nil, fmt.Errorf("no grade component columns matched the header")
	}

	result := &request.ImportScoreResult{}

	for rIdx, row := range rows[1:] {
		lineNo := rIdx + 2
		if codeCol >= len(row) {
			continue
		}
		code := strings.TrimSpace(row[codeCol])
		if code == "" {
			continue
		}
		var userID int64
		err := s.db.WithContext(ctx).
			Table("user").
			Select("user.id").
			Joins("JOIN user_class ON user_class.user_id = user.id").
			Where("user_class.class_id = ? AND user.code = ?", req.ClassID, code).
			Scan(&userID).Error
		if err != nil || userID == 0 {
			result.Skipped++
			result.Errors = append(result.Errors,
				fmt.Sprintf("line %d: student code '%s' not found in this class", lineNo, code))
			continue
		}
		details := make([]model.ScoreDetail, 0, len(componentCols))
		for col, gcID := range componentCols {
			var val float64
			if col < len(row) {
				raw := strings.TrimSpace(row[col])
				if raw != "" {
					v, perr := strconv.ParseFloat(raw, 64)
					if perr != nil {
						result.Errors = append(result.Errors,
							fmt.Sprintf("line %d: invalid score '%s'", lineNo, raw))
						continue
					}
					val = v
				}
			}
			details = append(details, model.ScoreDetail{
				GradeComponentID: gcID,
				Score:            val,
			})
		}
		txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			score := model.Score{
				ClassID:      req.ClassID,
				UserID:       userID,
				MajorID:      req.MajorID,
				GenerationID: req.GenerationID,
				ProgrammeID:  req.ProgrammeID,
				SubjectID:    req.SubjectID,
				Year:         req.Year,
				Semester:     req.Semester,
				IsActive:     true,
			}
			if err := tx.Create(&score).Error; err != nil {
				return err
			}
			for _, d := range details {
				d.ScoreID = int64(score.ID)
				if err := tx.Create(&d).Error; err != nil {
					return err
				}
			}
			return nil
		})
		if txErr != nil {
			result.Skipped++
			result.Errors = append(result.Errors,
				fmt.Sprintf("line %d (code %s): %v", lineNo, code, txErr))
			continue
		}
		result.Imported++
	}
	return result, nil
}

func applyAccessFilterScore(query *gorm.DB, db *gorm.DB, role model.Role, user model.User) *gorm.DB {
	switch {
	case role.Level <= 1:
		return query.Where("u.id = ?", user.ID)

	case role.Level > 1 && role.Level <= 6:
		var classIDs []int
		db.WithContext(query.Statement.Context).
			Model(&model.UserClass{}).
			Where("user_id = ?", user.ID).
			Pluck("class_id", &classIDs)

		if len(classIDs) == 0 {
			return query.Where("1 = 0")
		}
		return query.Where("s.class_id IN ?", classIDs)

	case role.Level >= 7:
		return query

	default:
		return query.Where("1 = 0")
	}
}

func applyCommonFilterScore(query *gorm.DB, filter map[string]string) *gorm.DB {
	for key, value := range filter {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		switch key {
		case "name":
			query = query.Where("u.name_kh LIKE ?", "%"+value+"%")
		case "class_id":
			query = query.Where("s.class_id =?", value)
		case "generation_id":
			query = query.Where("s.generation_id =?", value)
		case "major_id":
			query = query.Where("s.major_id =?", value)
		case "programme_id":
			query = query.Where("s.programme_id =?", value)
		case "subject_id":
			query = query.Where("s.subject_id =?", value)
		case "code":
			query = query.Where("u.code =?", value)
		case "year":
			query = query.Where("s.year =?", value)
		case "semester":
			query = query.Where("s.semester =?", value)
		}
	}
	return query
}

func (s *scoreservice) GetScore(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.ScoreResponse, *model.PaginationMetadata, error) {
	var data []response.ScoreResponse

	var user model.User
	if err := s.db.WithContext(ctx).Preload("Role").First(&user, userID).Error; err != nil {
		return nil, nil, err
	}
	offset := (pf.Page - 1) * pf.PageSize

	scorequery := s.db.WithContext(ctx).Table("score s").
		Select(`
		s.id AS id,
		u.id AS user_id,
		u.name_kh AS name_kh,
		u.name_en AS name_en,
		u.code AS code,
		u.gender AS gender,
		c.id AS class_id,
		c.name AS class_name,
		p.id AS programme_id,
		p.name AS programme_name,
		g.id AS generation_id,
		g.name_kh AS generation_name,
		m.id AS major_id,
		m.name_kh AS major_name,
		sb.id AS subject_id,
		sb.name_kh AS subject_name
	`).
		Joins("LEFT JOIN class c ON c.id = s.class_id").
		Joins("LEFT JOIN generation g ON g.id = s.generation_id").
		Joins("LEFT JOIN major m ON m.id = s.major_id").
		Joins("LEFT JOIN programmes p ON p.id = s.programme_id").
		Joins("LEFT JOIN subject sb ON sb.id = s.subject_id").
		Joins("LEFT JOIN user u ON u.id = s.user_id")

	scorequery = scorequery.Order("s.id DESC")
	scorequery = applyAccessFilterScore(scorequery, s.db, user.Role, user)
	scorequery = applyCommonFilterScore(scorequery, filter)

	var totalCount int64
	countQuery := scorequery.Session(&gorm.Session{})
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, nil, err
	}
	if err := scorequery.Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, err
	}
	if len(data) == 0 {
		return data, helper.BuildPaginationMeta(pf, totalCount), nil
	}
	scoreIDs := make([]int, len(data))
	scoreIndexByID := make(map[int]int, len(data))
	for i, a := range data {
		scoreIDs[i] = a.ID
		scoreIndexByID[a.ID] = i
	}

	var detail []response.ScoreDetailResponse
	detailquery := s.db.WithContext(ctx).Table("score_detail sd").
		Select(`
		sd.id AS id,
		sd.score_id AS score_id,
		g.id  AS grade_component_id,
		g.name AS grade_component_Name,
		sd.score AS score
	`).
		Joins("LEFT JOIN grade_components g ON g.id = sd.grade_component_id").
		Where("sd.score_id IN ?", scoreIDs)

	if err := detailquery.Scan(&detail).Error; err != nil {
		return nil, nil, err
	}

	for i := range detail {
		r := &detail[i]
		idx, ok := scoreIndexByID[r.ScoreID]
		if !ok {
			continue
		}
		switch r.GradeComponentName {
		case "វត្តមាននិស្សិត":
			data[idx].Attendance = r.Score
		case "កិច្ចការស្រាវជ្រាវ":
			data[idx].Research = r.Score
		case "ប្រឡងពាក់កណ្តាលឆមាស":
			data[idx].Midterm = r.Score
		case "ប្រឡងបញ្ចប់ឆមាស":
			data[idx].Final = r.Score
		}
	}
	return data, helper.BuildPaginationMeta(pf, totalCount), nil

}
