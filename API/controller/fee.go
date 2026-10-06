package controller

import (
	"mysql/constant/share"
	"mysql/request"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FeeController struct {
	service service.FeeService
}

func NewFeeController(s service.FeeService) *FeeController {
	return &FeeController{
		service: s,
	}
}

func (cr *FeeController) GetFeeSchedule(c *gin.Context) {
	data, err := cr.service.GetFeeSchedule(c.Request.Context())
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	share.RespondDate(c, http.StatusOK, data)
}

func (cr *FeeController) AddFee(c *gin.Context) {
	var input request.FeeRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.AddFee(c.Request.Context(), input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "Create Success")
}
