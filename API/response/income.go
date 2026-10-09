package response

import (
	"mysql/model"
	"mysql/model/base"
)

type IncomeResponse struct {
	base.ModelBase
	IncomeCode     string                  `json:"income_code"`
	CustomerID     uint64                  `json:"customer_id"`
	CustomerName   string                  `json:"customer_name"`
	IncomeDate     string                  `json:"income_date"`
	DueDate        *string                 `json:"due_date"`
	Subtotal       float64                 `json:"subtotal"`
	Tax            float64                 `json:"tax"`
	Discount       float64                 `json:"discount"`
	Total          float64                 `json:"total"`
	Paid           float64                 `json:"paid"`
	Status         model.IncomeStatus      `json:"status"`
	Description    *string                 `json:"description"`
	IncomeItems    []IncomeItemsResponse   `json:"income_items" gorm:"-"`
	IncomePayments []IncomePaymentResponse `json:"income_payments" gorm:"-"`
}

type IncomeItemsResponse struct {
	base.ModelBase
	IncomeID           uint64  `json:"income_id"`
	IncomeCategoryID   *int    `json:"income_category_id"`
	IncomeCategoryName *string `json:"income_category_name"`
	Description        *string `json:"description"`
	Quantity           float64 `json:"quantity"`
	UnitPrice          float64 `json:"unit_price"`
	SortOrder          int     `json:"sort_order"`
}

type IncomePaymentResponse struct {
	base.ModelBase
	IncomeID    uint64  `json:"income_id"`
	PaymentCode string  `json:"payment_code"`
	PaymentDate string  `json:"payment_date"`
	Amount      float64 `json:"amount"`
	Method      string  `json:"method"`
	Reference   *string `json:"reference"`
	Description *string `json:"description"`
}
