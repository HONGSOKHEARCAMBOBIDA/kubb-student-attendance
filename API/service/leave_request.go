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
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LeaveRequestService interface {
	CreateLeaveRequest(ctx context.Context, id int, input request.LeaveRequestCreate) error
	UpdateLeaveRequest(ctx context.Context, id int, userID int, input request.LeaveRequestUpdate) error
	DeleteLeaveRequest(ctx context.Context, id int) error
	GetLeaveRequest(ctx context.Context, id int, pf request.Pagination, filter map[string]string) ([]response.LeaveRequestResponse, *model.PaginationMetadata, error)
	VerifyLeaveRequest(ctx context.Context, id int, verifyBy int) error
}

type leaveRequestService struct {
	db *gorm.DB
}

func NewLeaveRequestService() LeaveRequestService {
	return &leaveRequestService{
		db: config.DB,
	}
}

func applyAccessFilterLeaveRequest(query *gorm.DB, db *gorm.DB, role model.Role, user model.User) *gorm.DB {
	if role.Level > 1 && role.Level < 7 {

		var classIDs []int
		db.Model(&model.UserClass{}).Where("user_id =?", user.ID).Pluck("class_id", &classIDs)
		if len(classIDs) == 0 {
			return query.Where("1 = 0")
		}
		return query.Where("l.class_id IN ?", classIDs)

	} else if role.Level <= 1 {
		return query.Where("l.user_id =?", user.ID)
	} else if role.Level >= 7 {
		return query
	}

	return query
}

func (s *leaveRequestService) GetLeaveRequest(ctx context.Context, id int, pf request.Pagination, filter map[string]string) ([]response.LeaveRequestResponse, *model.PaginationMetadata, error) {
	var data []response.LeaveRequestResponse
	var user model.User
	if err := s.db.WithContext(ctx).Preload("Role").First(&user, id).Error; err != nil {
		return nil, nil, err
	}
	helper.NormalizePagination(&pf)
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("leave_request l").
			Joins("LEFT JOIN user u ON u.id = l.user_id").
			Joins("LEFT JOIN class c ON c.id = l.class_id").
			Joins("LEFT JOIN leave_deduct_type ld ON ld.id = l.deduct_type_id").
			Joins("LEFT JOIN user ua ON ua.id = l.approve_by")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["name"]; ok && v != "" {
			tx = tx.Where("u.name LIKE ?", "%"+v+"%")
		}
		if v, ok := filter["class_id"]; ok && v != "" {
			tx = tx.Where("l.class_id = ?", v)
		}
		if v, ok := filter["status"]; ok && v != "" {
			tx = tx.Where("l.status = ?", v)
		}
		return tx
	}
	if err := applyFilters(base()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count product: %w", err)
	}

	if total == 0 {
		return []response.LeaveRequestResponse{}, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize
	query := applyFilters(base()).Table("leave_request l").
		Select(`
		l.id AS id,
		u.id AS user_id,
		u.gender AS gender,
		u.name_kh AS user_name_kh,
		u.name_en AS user_name_en,
		u.code AS user_code,
		c.name AS class_name,
		l.start_date AS start_date,
		l.end_date AS end_date,
		l.back_to_work_date AS back_to_work_date,
		l.total_day AS total_day,
		ld.id AS deduct_type_id,
		ld.code AS deduct_type_code,
		ld.name AS deduct_type_name,
		l.reason AS reason,
		l.status AS status,
		ua.id AS approve_by,
		ua.name_kh AS approve_by_name,
		l.approved_at AS approved_at
	`)

	query = applyAccessFilterLeaveRequest(query, s.db, user.Role, user)
	if err := query.Offset(offset).Limit(pf.PageSize).Order("id DESC").Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch class: %w", err)
	}

	for i := range data {
		data[i].StartDate = helper.FormatDate(data[i].StartDate)
		data[i].EndDate = helper.FormatDate(data[i].EndDate)
		data[i].BackToWorkDate = helper.FormatDate(data[i].BackToWorkDate)
		data[i].ApproveAt = helper.FormatDate(data[i].ApproveAt)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *leaveRequestService) CreateLeaveRequest(ctx context.Context, id int, input request.LeaveRequestCreate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var user model.User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return err
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newleave := model.LeaveRequest{
			UserID:         id,
			ClassID:        input.ClassID,
			StartDate:      input.StartDate,
			EndDate:        input.EndDate,
			BackToWorkDate: input.BackToWorkDate,
			TotalDay:       input.TotalDay,
			DeductTypeID:   input.DeductTypeID,
			Reason:         input.Reason,
			Status:         model.LeaveStatusPending,
		}
		if err := tx.Create(&newleave).Error; err != nil {
			return err
		}
		return nil
	})
	return err

}

func (s *leaveRequestService) VerifyLeaveRequest(ctx context.Context, id int, verifyBy int) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	var leaveforupdte model.LeaveRequest
	if err := s.db.WithContext(ctx).Preload("LeaveDeductType").First(&leaveforupdte, id).Error; err != nil {
		return err
	}

	var user model.User
	if err := s.db.WithContext(ctx).Select("id").First(&user, verifyBy).Error; err != nil {
		return err
	}

	leaveSession, err := leaveSessionFromCode(leaveforupdte.LeaveDeductType.Code)
	if err != nil {
		return err
	}

	var class model.Class
	if err := s.db.WithContext(ctx).First(&class, leaveforupdte.ClassID).Error; err != nil {
		return err
	}

	var shift model.Shift
	if err := s.db.WithContext(ctx).Where("id = ?", class.ShiftID).First(&shift).Error; err != nil {
		return err
	}

	sessions, err := buildSessionForLeave(shift, leaveSession)
	if err != nil {
		return err
	}

	const leaveScore = -0.5

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		updates := map[string]interface{}{
			"status":      model.LeaveStatusApprove,
			"approve_by":  verifyBy,
			"approved_at": now,
		}
		if err := tx.Model(&model.LeaveRequest{}).
			Where("id = ?", leaveforupdte.ID).
			Updates(updates).Error; err != nil {
			return fmt.Errorf("failed to approve leave request: %w", err)
		}

		startDate, err := helper.ParseLeaveDate(leaveforupdte.StartDate)
		if err != nil {
			return err
		}

		endDate, err := helper.ParseLeaveDate(leaveforupdte.EndDate)
		if err != nil {
			return err
		}

		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
			checkDate := d.Format("2006-01-02")

			var attendance model.Attendance
			err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("user_id = ? AND check_date = ?", leaveforupdte.UserID, checkDate).
				First(&attendance).Error

			switch {
			case err == nil:
				if err := tx.Model(&model.Attendance{}).Where("id = ?", attendance.ID).
					Update("leave_request_id", leaveforupdte.ID).Error; err != nil {
					return fmt.Errorf("failed to link attendance to leave: %w", err)
				}
			case errors.Is(err, gorm.ErrRecordNotFound):
				leaveID := leaveforupdte.ID
				attendance = model.Attendance{
					UserID:         leaveforupdte.UserID,
					ClassID:        leaveforupdte.ClassID,
					CheckDate:      checkDate,
					Status:         "LEAVE",
					LeaveRequestID: &leaveID,
				}
				if err := tx.Create(&attendance).Error; err != nil {
					return fmt.Errorf("failed to create attendance: %w", err)
				}
			default:
				return fmt.Errorf("failed to load attendance: %w", err)
			}

			for _, sess := range sessions {
				var existing model.AttendanceRecord
				err := tx.Where("attendance_id = ? AND type = ?", attendance.ID, sess.recordType).
					First(&existing).Error
				if err == nil {
					continue // real record already covers this session
				}
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("failed to check attendance record: %w", err)
				}

				record := model.AttendanceRecord{
					AttendanceID: attendance.ID,
					UserID:       leaveforupdte.UserID,
					ClassID:      leaveforupdte.ClassID,
					ShiftID:      shift.ID,
					CheckTime:    nil,
					Type:         sess.recordType,
					Inzone:       false,
					Score:        leaveScore,
					Status:       model.StatusPermission,
				}
				if err := tx.Create(&record).Error; err != nil {
					return fmt.Errorf("failed to create leave attendance record: %w", err)
				}

				if err := tx.Model(&model.Attendance{}).
					Where("id = ?", attendance.ID).
					UpdateColumn("score", gorm.Expr("score + ?", leaveScore)).Error; err != nil {
					return fmt.Errorf("failed to update attendance score: %w", err)
				}
			}
		}
		return nil
	})
}

// derives which sessions are covered, straight from the leave request's own deduct type
func leaveSessionFromCode(code string) (helper.LeaveSession, error) {
	switch code {
	case "FULL":
		return helper.LeaveFull, nil
	case "HALF_AM":
		return helper.LeaveMorning, nil
	case "HALF_PM":
		return helper.LeaveEvening, nil
	default:
		return helper.LeaveNone, fmt.Errorf("unknown leave deduct type code: %s", code)
	}
}

func buildSessionForLeave(shift model.Shift, leave helper.LeaveSession) ([]sessionConfig, error) {
	switch leave {
	case helper.LeaveFull:
		return []sessionConfig{
			{scheduledTime: shift.Session1, recordType: model.AttendanceSession1},
			{scheduledTime: shift.Session2, recordType: model.AttendanceSession2},
			{scheduledTime: shift.Session3, recordType: model.AttendanceSession3},
			{scheduledTime: shift.Session4, recordType: model.AttendanceSession4},
			{scheduledTime: shift.Session5, recordType: model.AttendanceSession5},
		}, nil
	case helper.LeaveMorning:
		return []sessionConfig{
			{scheduledTime: shift.Session1, recordType: model.AttendanceSession1},
			{scheduledTime: shift.Session2, recordType: model.AttendanceSession2},
			{scheduledTime: shift.Session3, recordType: model.AttendanceSession3},
		}, nil
	case helper.LeaveEvening:
		return []sessionConfig{
			{scheduledTime: shift.Session4, recordType: model.AttendanceSession4},
			{scheduledTime: shift.Session5, recordType: model.AttendanceSession5},
		}, nil
	default:
		return nil, errors.New("no approved leave sessions to record")
	}
}

func (s *leaveRequestService) UpdateLeaveRequest(ctx context.Context, id int, userID int, input request.LeaveRequestUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var data model.LeaveRequest
		if err := tx.Where("id = ? AND status = ?", id, model.LeaveStatusPending).First(&data).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "classcurriculumn not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch classcurriculumn", nil)
		}
		data.ClassID = input.ClassID
		data.StartDate = input.StartDate
		data.EndDate = input.EndDate
		data.BackToWorkDate = input.BackToWorkDate
		data.TotalDay = input.TotalDay
		data.DeductTypeID = input.DeductTypeID
		data.Reason = input.Reason
		data.ApprovedAt = nil
		if err := tx.Save(&data).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update product", nil)
		}
		return nil
	})
	return err
}

func (s *leaveRequestService) DeleteLeaveRequest(ctx context.Context, id int) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.
			Where("id = ? AND payroll_id IS NULL", id).
			Delete(&model.LeaveRequest{})

		if result.Error != nil {
			return fmt.Errorf("failed to delete leave request: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("leave request not found or has already been processed in payroll")
		}

		return nil
	})
}
