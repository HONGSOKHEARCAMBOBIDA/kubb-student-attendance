package controller

import (
	"context"
	"errors"
	"mysql/constant/share"
	"mysql/helper"
	"mysql/request"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IncomeController struct {
	service service.IncomeService
}

func NewIncomeController(s service.IncomeService) *IncomeController {
	return &IncomeController{
		service: s,
	}
}

func (cr *IncomeController) GetIncomeCategory(c *gin.Context) {
	data, err := cr.service.GetIncomeCategory(c.Request.Context())
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *IncomeController) AddIncome(c *gin.Context) {
	var input request.IncomeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.AddIncome(c.Request.Context(), input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "Income Create")
}

func (cr *IncomeController) GetIncome(c *gin.Context) {

	page, pageSize := helper.GetPagination(c)
	filter := map[string]string{
		"name":        c.Query("name"),
		"income_date": c.Query("income_date"),
	}
	data, meta, err := cr.service.GetIncome(c.Request.Context(), request.Pagination{
		Page:     page,
		PageSize: pageSize,
	}, filter)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			share.ResponseError(c, http.StatusGatewayTimeout, err.Error())
			return
		}
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponsePagination(c, 200, data, meta)
}

func (cr *IncomeController) DeleteIncome(c *gin.Context) {
	id, ok := helper.GetParamID(c)
	if !ok {
		return
	}
	if err := cr.service.DeleteIncome(c, id); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "Delete")
}
