package controller

import (
	"context"
	"errors"
	"mysql/constant/share"
	"mysql/helper"
	"mysql/pattern/res"
	"mysql/request"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	service service.ReportService
}

func NewReportController() ReportController {
	return ReportController{
		service: service.NewReportService(),
	}
}

func (cr *ReportController) GetWeekSchedule(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)
	filter := map[string]string{
		"generation_id": c.Query("generation_id"),
		"programme_id":  c.Query("programme_id"),
		"start_date":    c.Query("start_date"),
		"end_date":      c.Query("end_date"),
	}

	data, meta, err := cr.service.GetWeekSchedule(c.Request.Context(), request.Pagination{
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

func (cr *ReportController) AcademicConfirmation(c *gin.Context) {
	id, ok := helper.GetParamID(c)
	if !ok {
		return
	}
	data, err := cr.service.GetAcademicConfirmation(c, id)
	if err != nil {
		res.Fail(c, 500, "NO DATA", err.Error())
		return
	}
	share.RespondDate(c, 200, data)
}

func (cr *ReportController) TemporaryCertificate(c *gin.Context) {
	id, ok := helper.GetParamID(c)
	if !ok {
		return
	}
	data, err := cr.service.GetTemporaryCertificate(c, id)
	if err != nil {
		res.Fail(c, 500, "NO DATA", err.Error())
		return
	}
	share.RespondDate(c, 200, data)
}
