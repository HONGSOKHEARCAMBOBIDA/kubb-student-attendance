package controller

import (
	"mysql/constant/share"
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
