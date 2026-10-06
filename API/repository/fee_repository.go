package repository

import (
	"context"
	"math"
	"mysql/constant/apperror"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"time"

	"gorm.io/gorm"
)

type FeeRepository interface {
	GetFeeSchedule(ctx context.Context) ([]model.FeeSchedule, error)
	AddFee(ctx context.Context, input request.FeeRequestCreate) error
}

type feerepository struct {
	db *gorm.DB
}

func NewFeeRepository(db *gorm.DB) FeeRepository {
	return &feerepository{
		db: db,
	}
}

func (r *feerepository) GetFeeSchedule(ctx context.Context) ([]model.FeeSchedule, error) {
	var data []model.FeeSchedule
	err := r.db.WithContext(ctx).Find(&data).Error
	return data, err
}

func (r *feerepository) AddFee(ctx context.Context, input request.FeeRequestCreate) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var majorprice model.MajorPrice
		var class model.Class
		var feeschedule model.FeeSchedule
		var scholarship model.Schoolarship
		userClass := model.UserClass{
			UserID:   int64(input.StudentID),
			ClassID:  int64(input.ClassID),
			IsActive: true,
		}
		if err := tx.Create(&userClass).Error; err != nil {
			return err
		}
		if err := tx.First(&majorprice, input.MajorPriceID).Error; err != nil {
			return err
		}
		if err := tx.First(&class, input.ClassID).Error; err != nil {
			return err
		}
		if err := tx.First(&feeschedule, input.FeeScheduleID).Error; err != nil {
			return err
		}
		if err := tx.First(&scholarship, input.SchoolarshipID).Error; err != nil {
			return err
		}
		BaseAmount := helper.GetFeeAmountPerYear(majorprice, feeschedule)
		scholarshipdiscount := helper.CalculateDiscountBySchoolarship(BaseAmount, &scholarship)
		NetAmount := BaseAmount - scholarshipdiscount
		fee := model.Fee{
			UserClassID:   userClass.ID,
			UserID:        input.StudentID,
			ClassID:       input.ClassID,
			MajorID:       int(*class.MajorID),
			GenerationID:  int(*class.GenerationID),
			ProgrammeID:   int(*class.ProgrammeID),
			Year:          class.Year,
			Term:          &class.Term,
			ScholarshipID: &scholarship.ID,
			FeeScheduleID: input.FeeScheduleID,
			Date:          input.Date,
			Amount:        BaseAmount,
			Discount:      scholarshipdiscount,
			Total:         NetAmount,
			Active:        true,
		}
		if err := tx.Create(&fee).Error; err != nil {
			return err
		}
		scheduleCount, monthInterval := helper.GetNextDueDate(feeschedule)
		paymentAmount := NetAmount / float64(scheduleCount)
		paymentAmount = math.Round(paymentAmount*100) / 100
		installments := make([]model.Installment, 0, scheduleCount)
		dueDate := time.Now()
		for i := 1; i <= scheduleCount; i++ {
			amount := paymentAmount
			if i == scheduleCount {
				paidSoFar := paymentAmount * float64(scheduleCount-1)
				amount = NetAmount - paidSoFar
			}
			installments = append(installments, model.Installment{
				FeeID:      uint64(fee.ID),
				SequenceNo: i,
				DueDate:    dueDate.Format("2006-01-02"),
				Amount:     amount,
				Status:     string(model.InstallmentStatusPending),
			})
			dueDate = dueDate.AddDate(0, monthInterval, 0)
		}
		if err := tx.Create(&installments).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create installments", nil)
		}
		return nil
	})
	return err
}
