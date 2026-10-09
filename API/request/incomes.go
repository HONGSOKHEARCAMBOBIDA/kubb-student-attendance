package request

type IncomeRequest struct {
	CustomerID  int                   `json:"customer_id"`
	DueDate     string                `json:"due_date"`
	Tax         float64               `json:"tax"`
	Discount    float64               `json:"discount"`
	Description *string               `json:"description"`
	IncomeItems []IncomeItemRequest   `json:"income_items"`
	Payment     *IncomePaymentRequest `json:"payment"`
}

type IncomeItemRequest struct {
	IncomeCategoryID int     `json:"income_category_id"`
	Description      string  `json:"description"`
	Qty              float64 `json:"qty"`
	UnitPrice        float64 `json:"unit_price"`
}

type IncomePaymentRequest struct {
	Amount      float64 `json:"amount"`
	Method      string  `json:"method"`
	Reference   *string `json:"reference"`
	Description *string `json:"description"`
}
