package repository

import (
	"context"
	"errors"
	"math"
	"mysql/constant/apperror"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/response"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FeeRepository interface {
	GetFeeSchedule(ctx context.Context) ([]model.FeeSchedule, error)
	AddFee(ctx context.Context, input request.FeeRequestCreate) error
	GetSchoolarship(ctx context.Context) ([]model.Schoolarship, error)
	GetUserClass(ctx context.Context, userID int) ([]response.UserClass, error)
	AddFeeTransaction(ctx context.Context, input request.FeeTransaction) error
	PrintInvoice(ctx context.Context, id int) (response.PrintInvoiceResponse, error)
	DeleteFeeTransaction(ctx context.Context, id int) error
}

type feerepository struct {
	db *gorm.DB
}

func NewFeeRepository(db *gorm.DB) FeeRepository {
	return &feerepository{
		db: db,
	}
}

func (r *feerepository) GetUserClass(ctx context.Context, userID int) ([]response.UserClass, error) {
	var data []response.UserClass

	err := r.db.WithContext(ctx).
		Table("user_class uc").
		Select(`
			uc.id AS id,
			u.id AS user_id,
			u.name_kh AS name_kh,
			u.name_en AS name_en,
			u.gender AS gender,
			u.code AS code,
			c.id AS class_id,
			c.name AS class_name,
			c.type AS type,
			c.is_active AS active,
			m.id AS major_id,
			m.name_kh AS major_name,
			s.id AS shift_id,
			s.name AS shift_name,
			g.id AS generation_id,
			g.name_kh AS generation_name,
			g.start_year AS generation_start,
			g.end_year AS generation_end,
			c.year AS year,
			c.`+"`group`"+` AS group_name,
			c.term AS term,
			p.id AS programme_id,
			p.name AS programme_name,
			sc.description AS schoolarship,
			f.description AS fee_schedule,
			COALESCE(fe.amount, 0) AS amount,
			COALESCE(fe.discount, 0) AS discount,
			COALESCE(fe.total, 0) AS total,
			COALESCE(fe.paid_amount,0) AS paid_amount,
			COALESCE(fe.id, 0) AS fee_id
		`).
		Joins("JOIN user u ON u.id = uc.user_id").
		Joins("JOIN class c ON c.id = uc.class_id").
		Joins("LEFT JOIN major m ON m.id = c.major_id").
		Joins("LEFT JOIN generation g ON g.id = c.generation_id").
		Joins("LEFT JOIN shift s ON s.id = c.shift_id").
		Joins("LEFT JOIN programmes p ON p.id = c.programme_id").
		Joins("INNER JOIN fees fe ON fe.user_class_id = uc.id").
		Joins("LEFT JOIN scholarships sc ON sc.id = fe.scholarship_id").
		Joins("LEFT JOIN fee_schedules f ON f.id = fe.fee_schedule_id").
		Where("uc.user_id = ?", userID).
		Scan(&data).Error
	if err != nil {
		return nil, err
	}

	feeIDs := make([]uint64, 0, len(data))
	seen := make(map[uint64]struct{}, len(data))
	for _, d := range data {
		if d.FeeID == 0 {
			continue
		}
		if _, ok := seen[d.FeeID]; ok {
			continue
		}
		seen[d.FeeID] = struct{}{}
		feeIDs = append(feeIDs, d.FeeID)
	}
	if len(feeIDs) == 0 {
		return data, nil
	}

	var installments []response.InstallmentRespone
	if err := r.db.WithContext(ctx).
		Table("installments").
		Joins("LEFT JOIN fee_transactions ft ON ft.installment_id = installments.id").
		Where("installments.fee_id IN ?", feeIDs).
		Order("installments.fee_id, installments.sequence_no").
		Select(`
		installments.id AS id,
		installments.fee_id AS fee_id,
		installments.sequence_no AS sequence_no,
		installments.due_date AS due_date,
		installments.amount AS amount,
		installments.status AS status,
		ft.id AS fee_transacntion_id,
		ft.total AS fee_transaction_total
		`).Scan(&installments).Error; err != nil {
		return nil, err
	}

	for i := range installments {
		installments[i].DueDate = helper.FormatDate(installments[i].DueDate)
	}

	byFee := make(map[uint64][]response.InstallmentRespone, len(feeIDs))
	for _, inst := range installments {
		byFee[inst.FeeID] = append(byFee[inst.FeeID], inst)
	}

	for i := range data {
		if list, ok := byFee[data[i].FeeID]; ok {
			data[i].InstallmentRespone = list
		} else {
			data[i].InstallmentRespone = []response.InstallmentRespone{}
		}
	}

	return data, nil
}

func (r *feerepository) GetFeeSchedule(ctx context.Context) ([]model.FeeSchedule, error) {
	var data []model.FeeSchedule
	err := r.db.WithContext(ctx).Find(&data).Error
	return data, err
}

func (r *feerepository) GetSchoolarship(ctx context.Context) ([]model.Schoolarship, error) {
	var data []model.Schoolarship
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
			Status:   model.UserClassStatusSTUDY,
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
			PaidAmount:    0,
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

func (r *feerepository) AddFeeTransaction(ctx context.Context, input request.FeeTransaction) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var installment model.Installment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&installment, "id = ?", input.InstallmentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "installment not found", err)
			}
			return err
		}
		if installment.Status == string(model.InstallmentStatusPaid) {
			return apperror.New(apperror.CodeConflict, "installment already paid", nil)
		}

		var fee model.Fee
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&fee, "id = ?", installment.FeeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "fee not found", err)
			}
			return err
		}

		// TODO: validate input.Total == input.Amount - input.Discount + input.Tax
		// and that it matches the installment amount.

		ft := model.FeeTransaction{
			FeeID:         fee.ID,
			InstallmentID: &installment.ID,
			Date:          input.Date,
			DueDate:       &installment.DueDate,
			Amount:        input.Amount,
			Discount:      input.Discount,
			Tax:           input.Tax,
			Total:         input.Total,
			Reference:     input.Reference,
			Method:        input.Method,
			Message:       input.Message,
			Description:   input.Description,
			Active:        true,
		}
		if err := tx.Create(&ft).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create fee transaction", err)
		}

		if err := tx.Model(&ft).Updates(map[string]interface{}{
			"code":      helper.GenerateCode("PAY", uint(ft.ID)),
			"reference": helper.GenerateCode("REF", uint(ft.ID)),
		}).Error; err != nil {
			return apperror.New(
				apperror.CodeInternal,
				"failed to update transaction",
				err,
			)
		}

		if err := tx.Model(&installment).
			Update("status", string(model.InstallmentStatusPaid)).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update installment status", err)
		}

		if err := tx.Model(&fee).
			UpdateColumn("paid_amount", gorm.Expr("paid_amount + ?", ft.Total)).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update fee total", err)
		}
		return nil
	})
}

func (r *feerepository) PrintInvoice(ctx context.Context, id int) (response.PrintInvoiceResponse, error) {
	var data response.PrintInvoiceResponse

	res := r.db.WithContext(ctx).
		Table("fee_transactions ft").
		Select(`
			u.name_kh        AS name_kh,
			u.name_en        AS name_en,
			u.gender         AS gender,
			u.code           AS code,
			p.name           AS programme_name,
			g.code        AS generation_name,
			m.name_kh        AS major_name,
			c.term           AS term,
			c.year           AS year,
			ft.code          AS transacntion_code,
			ft.date          AS date,
			ft.due_date      AS due_date,
			ft.amount        AS amount,
			ft.discount      AS discount,
			ft.tax           AS tax,
			ft.total         AS total,
			ft.reference     AS reference,
			ft.method        AS method,
			ft.message       AS message,
			ft.description   AS description,
			i.sequence_no AS sequence_no
		`).
		Joins("LEFT JOIN fees f ON f.id = ft.fee_id").
		Joins("LEFT JOIN installments i ON i.id = ft.installment_id").
		Joins("LEFT JOIN user u ON u.id = f.user_id").
		Joins("LEFT JOIN programmes p ON p.id = f.programme_id").
		Joins("LEFT JOIN generation g ON g.id = f.generation_id").
		Joins("LEFT JOIN major m ON m.id = f.major_id").
		Joins("LEFT JOIN class c ON c.id = f.class_id").
		Where("ft.id = ?", id).
		Scan(&data)
	data.Date = helper.FormatDate(data.Date)
	if res.Error != nil {
		return response.PrintInvoiceResponse{}, res.Error
	}
	if res.RowsAffected == 0 {
		return response.PrintInvoiceResponse{}, gorm.ErrRecordNotFound
	}
	return data, nil
}

func (r *feerepository) DeleteFeeTransaction(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var ft model.FeeTransaction
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&ft, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "fee transaction not found", err)
			}
			return apperror.New(apperror.CodeInternal, "failed to load fee transaction", err)
		}

		// Reverse the payment on the fee
		res := tx.Model(&model.Fee{}).
			Where("id = ?", ft.FeeID).
			UpdateColumn("paid_amount", gorm.Expr("paid_amount - ?", ft.Total))
		if res.Error != nil {
			return apperror.New(apperror.CodeInternal, "failed to update fee paid amount", res.Error)
		}
		if res.RowsAffected == 0 {
			return apperror.New(apperror.CodeNotFound, "fee not found", nil)
		}

		// Reopen the installment
		res = tx.Model(&model.Installment{}).
			Where("id = ?", ft.InstallmentID).
			Update("status", string(model.InstallmentStatusPending))
		if res.Error != nil {
			return apperror.New(apperror.CodeInternal, "failed to update installment status", res.Error)
		}
		if res.RowsAffected == 0 {
			return apperror.New(apperror.CodeNotFound, "installment not found", nil)
		}

		if err := tx.Delete(&ft).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to delete fee transaction", err)
		}
		return nil
	})
}
