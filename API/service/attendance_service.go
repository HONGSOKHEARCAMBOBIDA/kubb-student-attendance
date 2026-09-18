package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	// "html"
	"mysql/config"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/utils"

	//"mysql/request"
	"mysql/response"
	// "mysql/utils"
	// "strconv"
	// "strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AttendanceService interface {
	CreateAttendance(ctx context.Context, id int, input request.AttendanceRequestCreate) error
	//GetAttendance(ctx context.Context, id int, pf request.Pagination, filter map[string]string) ([]response.AttendanceResponse, *model.PaginationMetadata, error)
	GetAttendanceDraft(ctx context.Context, id int) (response.AttendanceResponseDraft, error)
	GetAttendancePDF(ctx context.Context, id int, pf request.Pagination, filter map[string]string) ([]response.AttendanceResponseGenerate, *model.PaginationMetadata, error)
	//DeleteAttendance(ctx context.Context, id int) error
	GetAttendanceReport(ctx context.Context, id int, filter map[string]string) (*response.AttendanceReportResponse, error)
}

type attendanceservice struct {
	db *gorm.DB
}

func NewAttendanceService() AttendanceService {
	return &attendanceservice{
		db: config.DB,
	}
}

type LeaveSession int

const (
	LeaveNone LeaveSession = iota
	LeaveFull
	LeaveMorning
	LeaveEvening
)

var checkTypeKubbLabel = map[string]string{
	model.AttendanceSession1: "ចូលរៀនSessionទី១",
	model.AttendanceSession2: "ចូលរៀនSessionទី២",
	model.AttendanceSession3: "ចូលរៀនSessionទី៣",
	model.AttendanceSession4: "ចូលរៀនSessionទី៤",
	model.AttendanceSession5: "ចូលរៀនSessionទី៥",
}

var recordTypeLabel = map[string]string{
	model.AttendanceSession1: "ចូលរៀនSessionទី១",
	model.AttendanceSession2: "ចូលរៀនSessionទី២",
	model.AttendanceSession3: "ចូលរៀនSessionទី៣",
	model.AttendanceSession4: "ចូលរៀនSessionទី៤",
	model.AttendanceSession5: "ចូលរៀនSessionទី៥",
}

type sessionConfig struct {
	scheduledTime string
	recordType    string
}

var ErrAllSessionsRecorded = errors.New("all attendance sessions for today have already been recorded")

func buildSessionV2(shift model.Shift, leave LeaveSession) ([]sessionConfig, error) {
	if leave == LeaveFull {
		return nil, errors.New("today is a full-day approved leave")
	}
	switch leave {
	case LeaveMorning:
		return []sessionConfig{
			{scheduledTime: shift.Session4, recordType: model.AttendanceSession4},
			{scheduledTime: shift.Session5, recordType: model.AttendanceSession5},
		}, nil
	case LeaveEvening:
		return []sessionConfig{
			{scheduledTime: shift.Session1, recordType: model.AttendanceSession1},
			{scheduledTime: shift.Session2, recordType: model.AttendanceSession2},
			{scheduledTime: shift.Session3, recordType: model.AttendanceSession3},
		}, nil
	default:
		return []sessionConfig{
			{scheduledTime: shift.Session1, recordType: model.AttendanceSession1},
			{scheduledTime: shift.Session2, recordType: model.AttendanceSession2},
			{scheduledTime: shift.Session3, recordType: model.AttendanceSession3},
			{scheduledTime: shift.Session4, recordType: model.AttendanceSession4},
			{scheduledTime: shift.Session5, recordType: model.AttendanceSession5},
		}, nil
	}

}

func (s *attendanceservice) getApprovedLeaveSession(ctx context.Context, userID int) (LeaveSession, error) {
	currentDate := helper.CurrentDate()

	var leaveRequest model.LeaveRequest
	err := s.db.WithContext(ctx).
		Preload("LeaveDeductType").
		Where("user_id = ? AND status = ? AND start_date <= ? AND end_date >= ?",
			userID, model.LeaveStatusApprove, currentDate, currentDate).Order("id ASC").
		First(&leaveRequest).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// no approved leave today -> normal full-day attendance flow
			return LeaveNone, nil
		}
		return LeaveNone, fmt.Errorf("failed to load leave request: %w", err)
	}

	deductCode := leaveRequest.LeaveDeductType.Code
	switch deductCode {
	case "FULL":
		return LeaveFull, nil
	case "HALF_AM":
		return LeaveMorning, nil
	case "HALF_PM":
		return LeaveEvening, nil
	default:
		return LeaveNone, fmt.Errorf("unknown leave deduct type code: %s", deductCode)
	}
}

func (s *attendanceservice) CreateAttendance(ctx context.Context, id int, input request.AttendanceRequestCreate) error {
	currentDate := helper.CurrentDate()
	currentTime := helper.CurrentTime()

	var class model.Class

	if err := s.db.WithContext(ctx).First(&class, input.CompanyID).Error; err != nil {
		return fmt.Errorf("failed to load company: %w", err)
	}

	var user model.User
	if err := s.db.WithContext(ctx).Select("id").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return err
	}

	var userclass model.UserClass
	if err := s.db.WithContext(ctx).Select("id,class_id,user_id").
		Where("user_id = ? AND is_active = 1", user.ID).First(&userclass).Error; err != nil {
		return err
	}

	var shift model.Shift
	if err := s.db.WithContext(ctx).Where("id = ?", class.ShiftID).First(&shift).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return err
	}

	leave, err := s.getApprovedLeaveSession(ctx, user.ID)
	if err != nil {
		return err
	}

	sessions, err := buildSessionV2(shift, leave)
	if err != nil {
		return err
	}

	companyLat, err := strconv.ParseFloat(*class.Latitude, 64)
	if err != nil {
		return fmt.Errorf("invalid company latitude: %w", err)
	}

	companyLng, err := strconv.ParseFloat(*class.Longitude, 64)
	if err != nil {
		return fmt.Errorf("invalid company longitude: %w", err)
	}

	userLat, err := strconv.ParseFloat(input.Latitude, 64)
	if err != nil {
		return fmt.Errorf("invalid user latitude: %w", err)
	}

	userLng, err := strconv.ParseFloat(input.Longitude, 64)
	if err != nil {
		return fmt.Errorf("invalid user longitude :%w", err)
	}

	radius, err := strconv.ParseFloat(class.Radius, 64)
	if err != nil {
		return fmt.Errorf("invalid company redius: %w", err)
	}

	distance := utils.CalculateDistance(companyLat, companyLng, userLat, userLng)
	inzone := distance <= radius

	if !inzone && class.CanScanOutsize == false {
		return errors.New("អ្នកមិនអាចស្កែនក្រៅតំបន់សាលាបានទេ")
	}

	var current sessionConfig
	var record model.AttendanceRecord
	// var attendanceID uint
	// var justCompleted bool

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var attendance model.Attendance
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND check_date = ?", user.ID, currentDate).First(&attendance).Error
		switch {
		case err == nil:

		case errors.Is(err, gorm.ErrRecordNotFound):
			attendance = model.Attendance{
				UserID:         user.ID,
				ClassID:        input.CompanyID,
				CheckDate:      currentDate,
				Status:         "WORKING",
				LeaveRequestID: nil,
				VerifyBy:       nil,
			}
			if err := tx.Create(&attendance).Error; err != nil {
				return fmt.Errorf("failed to create attendance: %w", err)
			}
		default:
			return fmt.Errorf("failed to load attendance: %w", err)
		}
		var existingRecords []model.AttendanceRecord
		if err := tx.Where("attendance_id = ? AND status = ?", attendance.ID, model.StatusPresent).Order("id ASC").Find(&existingRecords).Error; err != nil {
			return fmt.Errorf("failed to load attendance :%w", err)
		}

		recordCount := len(existingRecords)
		maxRecords := len(sessions)

		if recordCount >= maxRecords {
			if err := tx.Model(&model.Attendance{}).Where("id = ?", attendance.ID).
				Update("status", "COMPLETE").Error; err != nil {
				return fmt.Errorf("failed to update attendance status:%w", err)
			}
			return errors.New("all check-in and check-out completed for today")
		}
		current = sessions[recordCount]
		//attendanceType := helper.DetermineAttendanceType(currentTime, current.scheduledTime, current.isCheckIn)
		record = model.AttendanceRecord{
			AttendanceID: attendance.ID,
			UserID:       user.ID,
			ClassID:      input.CompanyID,
			ShiftID:      shift.ID,
			CheckTime:    &currentTime,
			Type:         current.recordType,
			Inzone:       inzone,
			Latitude:     input.Latitude,
			Longitude:    input.Longitude,
			Status:       model.StatusPresent,
		}
		if err := tx.Create(&record).Error; err != nil {
			return fmt.Errorf("failed to created attendance record: %w", err)
		}
		// attendanceID = uint(attendance.ID)
		if recordCount+1 >= maxRecords {
			if err := tx.Model(&model.Attendance{}).Where("id = ?", attendance.ID).
				Update("status", "COMPLETE").Error; err != nil {
				return fmt.Errorf("faild to update attendance status %w", err)
			}
			// justCompleted = true
		}

		return nil
	})
	if txErr != nil {
		return txErr
	}

	// _ = attendanceID
	// _ = justCompleted
	return nil
}

func (s *attendanceservice) GetAttendanceDraft(ctx context.Context, id int) (response.AttendanceResponseDraft, error) {
	currentDate := helper.CurrentDate()

	var user model.User
	if err := s.db.WithContext(ctx).Select("id").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.AttendanceResponseDraft{}, fmt.Errorf("user not found: %w", err)
		}
		return response.AttendanceResponseDraft{}, fmt.Errorf("failed to load user: %w", err)
	}

	var userclass model.UserClass
	if err := s.db.WithContext(ctx).Preload("Class").Select("id,class_id,user_id").
		Where("user_id = ? AND is_active = 1", user.ID).First(&userclass).Error; err != nil {
		return response.AttendanceResponseDraft{}, fmt.Errorf("failed to load user class: %w", err)
	}

	var shift model.Shift
	if err := s.db.WithContext(ctx).Where("id = ?", userclass.Class.ShiftID).First(&shift).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.AttendanceResponseDraft{}, fmt.Errorf("shift not found: %w", err)
		}
		return response.AttendanceResponseDraft{}, fmt.Errorf("failed to load shift: %w", err)
	}

	leave, err := s.getApprovedLeaveSession(ctx, user.ID)
	if err != nil {
		return response.AttendanceResponseDraft{}, err
	}

	sessions, err := buildSessionV2(shift, leave)
	if err != nil {
		return response.AttendanceResponseDraft{}, err
	}

	var current sessionConfig
	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var attendance model.Attendance
		var existingRecords []model.AttendanceRecord

		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND check_date = ?", user.ID, currentDate).First(&attendance).Error

		switch {
		case err == nil:
			if err := tx.Where("attendance_id = ? AND status = ?", attendance.ID, model.StatusPresent).Order("id ASC").Find(&existingRecords).Error; err != nil {
				return fmt.Errorf("failed to load attendance records: %w", err)
			}
		case errors.Is(err, gorm.ErrRecordNotFound):
		default:
			return fmt.Errorf("failed to load attendance: %w", err)
		}

		recordCount := len(existingRecords)
		if recordCount >= len(sessions) {
			return ErrAllSessionsRecorded
		}
		current = sessions[recordCount]
		return nil
	})
	if txErr != nil {
		return response.AttendanceResponseDraft{}, txErr
	}

	label, ok := recordTypeLabel[current.recordType]
	if !ok {
		return response.AttendanceResponseDraft{}, fmt.Errorf("unknown record type: %s", current.recordType)
	}
	return response.AttendanceResponseDraft{
		Type:          current.recordType, // now a string, e.g. "session1"
		TypeString:    label,
		ScheduledTime: current.scheduledTime,
	}, nil
}

func applyAccessFilterAttendance(query *gorm.DB, db *gorm.DB, role model.Role, user model.User) *gorm.DB {
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
		return query.Where("a.class_id IN ?", classIDs)

	case role.Level >= 7:
		return query

	default:
		return query.Where("1 = 0")
	}
}

func applyCommonFilterAttendance(query *gorm.DB, filter map[string]string) *gorm.DB {
	for key, value := range filter {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		switch key {
		case "name":
			query = query.Where("u.name_kh LIKE ?", "%"+value+"%")
		case "class_id":
			query = query.Where("a.class_id =?", value)
		case "check_date":
			query = query.Where("a.check_date >=?", value)
		}
	}
	return query
}

func (s *attendanceservice) GetAttendanceReport(ctx context.Context, id int, filter map[string]string) (*response.AttendanceReportResponse, error) {
	var user model.User
	if err := s.db.WithContext(ctx).Preload("Role").First(&user, id).Error; err != nil {
		return nil, err
	}

	type attendanceHeader struct {
		ID        int
		UserID    int
		NameKH    string
		NameEN    string
		Code      string
		Gender    int
		ClassName string
		CheckDate string
	}
	var headers []attendanceHeader

	q := s.db.WithContext(ctx).Table("attendance a").
		Select(`
			a.id AS id,
			u.id AS user_id,
			u.name_kh AS name_kh,
			u.name_en AS name_en,
			u.code AS code,
			u.gender AS gender,
			c.name AS class_name,
			a.check_date AS check_date
		`).
		Joins("LEFT JOIN user u ON u.id = a.user_id").
		Joins("LEFT JOIN class c ON c.id = a.class_id")

	q = applyAccessFilterAttendance(q, s.db, user.Role, user)
	q = applyCommonFilterAttendance(q, filter) // now also handles date_from/date_to
	q = q.Order("u.id ASC, a.check_date ASC")

	if err := q.Scan(&headers).Error; err != nil {
		return nil, err
	}
	if len(headers) == 0 {
		return &response.AttendanceReportResponse{}, nil
	}

	attendanceIDs := make([]int, 0, len(headers))
	headerByAttendanceID := make(map[int]attendanceHeader, len(headers))
	for _, h := range headers {
		attendanceIDs = append(attendanceIDs, h.ID)
		headerByAttendanceID[h.ID] = h
	}

	type recordRow struct {
		AttendanceID int
		Type         string
		CheckTime    string
		Status       string
	}
	var records []recordRow
	if err := s.db.WithContext(ctx).Table("attendance_record ar").
		Select("ar.attendance_id AS attendance_id, ar.type AS type, ar.check_time AS check_time,ar.status AS status").
		Where("ar.attendance_id IN ?", attendanceIDs).
		Scan(&records).Error; err != nil {
		return nil, err
	}

	type colKey struct {
		Date string
		Type string
	}
	sessionOrder := []string{
		model.AttendanceSession1, model.AttendanceSession2, model.AttendanceSession3,
		model.AttendanceSession4, model.AttendanceSession5,
	}

	colSet := map[colKey]bool{}
	dateSet := map[string]bool{}
	recordsByAttendanceID := map[int][]recordRow{}
	for _, r := range records {
		recordsByAttendanceID[r.AttendanceID] = append(recordsByAttendanceID[r.AttendanceID], r)
		h := headerByAttendanceID[r.AttendanceID]
		colSet[colKey{Date: h.CheckDate, Type: r.Type}] = true
		dateSet[h.CheckDate] = true
	}

	dates := make([]string, 0, len(dateSet))
	for d := range dateSet {
		dates = append(dates, d)
	}
	sort.Strings(dates) // works fine if CheckDate is yyyy-mm-dd

	var columns []response.AttendanceReportColumn
	colIndex := map[colKey]int{}
	colNo := 1
	for _, d := range dates {
		for _, t := range sessionOrder {
			if !colSet[colKey{Date: d, Type: t}] {
				continue
			}
			columns = append(columns, response.AttendanceReportColumn{
				ColumnNo:  colNo,
				CheckDate: d,
				Type:      t,
			})
			colIndex[colKey{Date: d, Type: t}] = len(columns) - 1
			colNo++
		}
	}

	type userAgg struct {
		row response.AttendanceReportRow
	}
	userRows := map[int]*userAgg{}
	var userOrder []int

	for _, h := range headers {
		agg, ok := userRows[h.UserID]
		if !ok {
			agg = &userAgg{row: response.AttendanceReportRow{
				UserID:    h.UserID,
				NameKH:    h.NameKH,
				NameEN:    h.NameEN,
				Code:      h.Code,
				Gender:    h.Gender,
				ClassName: h.ClassName,
				Cells:     make([]response.AttendanceReportCell, len(columns)),
			}}
			userRows[h.UserID] = agg
			userOrder = append(userOrder, h.UserID)
		}

		for _, r := range recordsByAttendanceID[h.ID] {
			idx, ok := colIndex[colKey{Date: h.CheckDate, Type: r.Type}]
			if !ok {
				continue
			}
			status := "A"

			if r.Status == model.StatusPermission {
				status = model.StatusPermission
				agg.row.PermissionCount++
			} else if r.Status == model.StatusAbsence {
				status = model.StatusAbsence
				agg.row.AbsentCount++
			} else if r.Status == model.StatusPresent {
				status = model.StatusPresent
			}
			agg.row.Cells[idx] = response.AttendanceReportCell{Status: status, CheckTime: r.CheckTime}
		}
	}

	rows := make([]response.AttendanceReportRow, 0, len(userOrder))
	for i, uid := range userOrder {
		agg := userRows[uid]
		agg.row.Index = i + 1
		rows = append(rows, agg.row)
	}

	return &response.AttendanceReportResponse{Columns: columns, Rows: rows}, nil
}

func (s *attendanceservice) GetAttendancePDF(ctx context.Context, id int, pf request.Pagination, filter map[string]string) ([]response.AttendanceResponseGenerate, *model.PaginationMetadata, error) {
	var attendance []response.AttendanceResponseGenerate
	var user model.User
	if err := s.db.WithContext(ctx).Preload("Role").First(&user, id).Error; err != nil {
		return nil, nil, err
	}
	offset := (pf.Page - 1) * pf.PageSize

	attendancequery := s.db.WithContext(ctx).Table("attendance a").
		Select(`
			a.class_id AS class_id,
			a.id AS id,
			u.id AS user_id,
			u.name_kh AS name_kh,
			u.name_en AS name_en,
			u.code AS code,
			u.gender AS gender,
			c.name AS class_name,
			a.check_date AS check_date,
			a.status AS status
		`).
		Joins("LEFT JOIN user u ON u.id = a.user_id").
		Joins("LEFT JOIN class c ON c.id = a.class_id")

	attendancequery = attendancequery.Order("a.id DESC")
	attendancequery = applyAccessFilterAttendance(attendancequery, s.db, user.Role, user)
	attendancequery = applyCommonFilterAttendance(attendancequery, filter)

	var totalCount int64
	countQuery := attendancequery.Session(&gorm.Session{})
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, nil, err
	}
	if err := attendancequery.Offset(offset).Limit(pf.PageSize).Scan(&attendance).Error; err != nil {
		return nil, nil, err
	}

	for i := range attendance {
		attendance[i].CheckDate = helper.FormatDate(attendance[i].CheckDate)
	}

	if len(attendance) == 0 {
		return attendance, helper.BuildPaginationMeta(pf, totalCount), nil
	}

	attendanceIDs := make([]int, len(attendance))
	attendanceIndexByID := make(map[int]int, len(attendance))
	for i, a := range attendance {
		attendanceIDs[i] = a.ID
		attendanceIndexByID[a.ID] = i
	}

	var attendancerecords []response.AttendanceRecordResponse

	attendancerecordquery := s.db.WithContext(ctx).Table("attendance_record ar").
		Select(`
			ar.id AS id,
			ar.attendance_id AS attendance_id,
			s.session1 AS session1,
			s.session2 AS session2,
			s.session3 AS session3,
			s.session4 AS session4,
			s.session5 AS session5,
			ar.check_time AS check_time,
			ar.type AS type,
			ar.inzone AS inzone,
			ar.latitude AS latitude,
			ar.longitude AS longitude
		`).
		Joins("LEFT JOIN shift s ON s.id = ar.shift_id").
		Where("ar.attendance_id IN ?", attendanceIDs)

	if err := attendancerecordquery.Scan(&attendancerecords).Error; err != nil {
		return nil, nil, err
	}

	for i := range attendancerecords {
		r := &attendancerecords[i]
		idx, ok := attendanceIndexByID[r.AttendanceID]
		if !ok {
			continue
		}

		switch r.Type {
		case model.AttendanceSession1:
			attendance[idx].Session1 = r.CheckTime
		case model.AttendanceSession2:
			attendance[idx].Session2 = r.CheckTime
		case model.AttendanceSession3:
			attendance[idx].Session3 = r.CheckTime
		case model.AttendanceSession4:
			attendance[idx].Session4 = r.CheckTime
		case model.AttendanceSession5:
			attendance[idx].Session5 = r.CheckTime
		case model.AttendanceSession6:
			attendance[idx].Session6 = r.CheckTime
		}
	}

	return attendance, helper.BuildPaginationMeta(pf, totalCount), nil
}

// func (s *attendanceservice) DeleteAttendance(ctx context.Context, id int) error {
// 	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

// 		if err := tx.
// 			Where("attendance_id = ?", id).
// 			Delete(&model.AttendanceRecord{}).Error; err != nil {

// 			return fmt.Errorf("failed to delete attendance record: %w", err)
// 		}

// 		if err := tx.
// 			Where("id = ?", id).
// 			Delete(&model.Attendance{}).Error; err != nil {

// 			return fmt.Errorf("failed to delete attendance: %w", err)
// 		}

// 		return nil
// 	})
// }
