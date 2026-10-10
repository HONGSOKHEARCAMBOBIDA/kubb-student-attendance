package repository

import (
	"context"
	"fmt"
	"math"
	"mysql/constant/apperror"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/response"

	"gorm.io/gorm"
)

type IncomeRepository interface {
	GetIncomeCategory(ctx context.Context) ([]model.IncomeCategory, error)
	AddIncome(ctx context.Context, input request.IncomeRequest) error
	GetIncome(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.IncomeResponse, *model.PaginationMetadata, error)
	DeleteIncome(ctx context.Context, id int) error
}

type incomerepository struct {
	db *gorm.DB
}

func NewIncomeRepository(db *gorm.DB) IncomeRepository {
	return &incomerepository{
		db: db,
	}
}

func (r *incomerepository) GetIncomeCategory(ctx context.Context) ([]model.IncomeCategory, error) {
	var data []model.IncomeCategory
	err := r.db.WithContext(ctx).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, err
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func (r *incomerepository) AddIncome(ctx context.Context, input request.IncomeRequest) error {
	var income model.Income
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if input.CustomerID <= 0 {
			return apperror.New(apperror.CodeInternal, "Customer ID required", &apperror.AppError{})
		}
		if len(input.IncomeItems) == 0 {
			return apperror.New(apperror.CodeInternal, "At least one income item required", &apperror.AppError{})
		}
		var subtotal float64
		for _, it := range input.IncomeItems {
			if it.Qty <= 0 || it.UnitPrice < 0 {
				return apperror.New(apperror.CodeInternal, "Invalid item quantity or price", &apperror.AppError{})
			}
			subtotal += it.Qty * it.UnitPrice
		}
		subtotal = round2(subtotal)
		total := round2(subtotal + input.Tax - input.Discount)
		if total < 0 {
			return apperror.New(apperror.CodeInternal, "Total cannot be negative", &apperror.AppError{})
		}
		var paid float64
		if input.Payment != nil {
			paid = round2(input.Payment.Amount)
			if paid <= 0 {
				return apperror.New(apperror.CodeInternal, "Payment amount must be greater than 0", &apperror.AppError{})
			}
			if paid > total {
				return apperror.New(apperror.CodeInternal, "Payment exceeds income total", &apperror.AppError{})
			}
		}

		status := model.IncomeStatusUNPAID
		switch {
		case paid > 0 && paid < total:
			status = model.IncomeStatusPARTIAL
		case paid > 0 && paid >= total:
			status = model.IncomeStatusPAID
		}

		now := helper.CurrentDate()
		income = model.Income{
			CustomerID:  input.CustomerID,
			IncomeDate:  now,
			DueDate:     input.DueDate,
			Subtotal:    subtotal,
			Tax:         input.Tax,
			Discount:    input.Discount,
			Total:       total,
			Paid:        paid,
			Status:      status,
			Description: input.Description,
		}
		if err := tx.Create(&income).Error; err != nil {
			return err
		}

		income.IncomeCode = helper.GenerateCode("ICO", uint(income.ID))
		if err := tx.Model(&income).Update("income_code", helper.GenerateCode("ICO", uint(income.ID))).Error; err != nil {
			return err
		}
		items := make([]model.IncomeItem, 0, len(input.IncomeItems))
		for i, it := range input.IncomeItems {
			item := model.IncomeItem{
				IncomeID:  uint64(income.ID),
				Quantity:  it.Qty,
				UnitPrice: it.UnitPrice,
				SortOrder: i + 1,
			}
			if it.IncomeCategoryID > 0 {
				cid := it.IncomeCategoryID
				item.IncomeCategoryID = &cid
			}
			if it.Description != "" {
				d := it.Description
				item.Description = &d
			}
			items = append(items, item)
		}
		if err := tx.Create(&items).Error; err != nil {
			return err
		}
		if input.Payment != nil {
			method := input.Payment.Method
			if method == "" {
				method = "CASH"
			}
			payment := model.IncomePayment{
				IncomeID:    uint64(income.ID),
				PaymentDate: now,
				Amount:      paid,
				Method:      method,
				Reference:   input.Payment.Reference,
				Description: input.Payment.Description,
			}
			if err := tx.Create(&payment).Error; err != nil {
				return err
			}
			if err := tx.Model(&payment).Update("payment_code", helper.GenerateCode("ICOPA", uint(payment.ID))).Error; err != nil {
				return err
			}
		}

		return nil
	})
	return err
}

func (r *incomerepository) GetIncome(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.IncomeResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)
	base := func() *gorm.DB {
		return r.db.WithContext(ctx).
			Table("incomes i").
			Joins("LEFT JOIN user u ON u.id = i.customer_id")
	}
	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v := filter["name"]; v != "" {
			like := "%" + v + "%"
			tx = tx.Where("(u.name_kh LIKE ? OR u.name_en LIKE ?)", like, like)
		}
		if v := filter["income_date"]; v != "" {
			tx = tx.Where("i.income_date >= ?", v)
		}
		return tx
	}
	var total int64
	if err := applyFilters(base()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count incomes: %w", err)
	}
	var data []response.IncomeResponse
	offset := (pf.Page - 1) * pf.PageSize
	err := applyFilters(base()).
		Select(`
			i.id            AS id,
			i.income_code   AS income_code,
			u.id            AS customer_id,
			u.name_kh       AS customer_name,
			i.income_date   AS income_date,
			i.due_date      AS due_date,
			i.subtotal      AS subtotal,
			i.tax           AS tax,
			i.discount      AS discount,
			i.total         AS total,
			i.paid          AS paid,
			i.status        AS status,
			i.description   AS description
		`).
		Order("i.id DESC").
		Offset(offset).
		Limit(pf.PageSize).
		Scan(&data).Error
	for i := range data {
		data[i].IncomeDate = helper.FormatDate(data[i].IncomeDate)
		*data[i].DueDate = helper.FormatDate(*data[i].DueDate)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("fetch incomes: %w", err)
	}
	if len(data) == 0 {
		return data, helper.BuildPaginationMeta(pf, total), nil
	}
	ids := make([]uint64, len(data))
	for i, d := range data {
		ids[i] = uint64(d.ID)
	}
	var items []response.IncomeItemsResponse
	err = r.db.WithContext(ctx).
		Table("income_items it").
		Joins("LEFT JOIN income_categories ic ON ic.id = it.income_category_id").
		Where("it.income_id IN ?", ids).
		Select(`
			it.id                AS id,
			it.income_id         AS income_id,
			ic.id                AS income_category_id,
			ic.name              AS income_category_name,
			it.description       AS description,
			it.quantity          AS quantity,
			it.unit_price        AS unit_price,
			it.sort_order        AS sort_order
		`).
		Order("it.sort_order ASC, it.id ASC").
		Scan(&items).Error
	if err != nil {
		return nil, nil, fmt.Errorf("fetch income items: %w", err)
	}
	var payments []response.IncomePaymentResponse
	err = r.db.WithContext(ctx).
		Table("income_payments ip").
		Where("ip.income_id IN ?", ids).
		Select(`
			ip.id            AS id,
			ip.income_id     AS income_id,
			ip.payment_code  AS payment_code,
			ip.payment_date  AS payment_date,
			ip.amount        AS amount,
			ip.method        AS method,
			ip.reference     AS reference,
			ip.description   AS description
		`).
		Order("ip.payment_date ASC, ip.id ASC").
		Scan(&payments).Error
	if err != nil {
		return nil, nil, fmt.Errorf("fetch income payments: %w", err)
	}
	itemsByIncome := make(map[uint64][]response.IncomeItemsResponse, len(data))
	for _, it := range items {
		itemsByIncome[it.IncomeID] = append(itemsByIncome[it.IncomeID], it)
	}
	paymentsByIncome := make(map[uint64][]response.IncomePaymentResponse, len(data))
	for _, p := range payments {
		paymentsByIncome[p.IncomeID] = append(paymentsByIncome[p.IncomeID], p)
	}
	for i := range data {
		id := uint64(data[i].ID)
		data[i].IncomeItems = itemsByIncome[id]
		data[i].IncomePayments = paymentsByIncome[id]
	}
	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (r *incomerepository) DeleteIncome(ctx context.Context, id int) error {

}
