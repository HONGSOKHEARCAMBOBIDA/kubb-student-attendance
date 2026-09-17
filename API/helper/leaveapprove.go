package helper

import (
	"context"
	"errors"
	"fmt"
	"mysql/model"

	"gorm.io/gorm"
)

type LeaveSession int

const (
	LeaveNone LeaveSession = iota
	LeaveFull
	LeaveMorning
	LeaveEvening
)

func GetApprovedLeaveSession(ctx context.Context, db *gorm.DB, userID int) (LeaveSession, error) {
	currentDate := CurrentDate()

	var leaveRequest model.LeaveRequest
	err := db.WithContext(ctx).
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
	case "FULL_DAY":
		return LeaveFull, nil
	case "HALF_AM":
		return LeaveMorning, nil
	case "HALF_PM":
		return LeaveEvening, nil
	default:
		return LeaveNone, fmt.Errorf("unknown leave deduct type code: %s", deductCode)
	}
}
