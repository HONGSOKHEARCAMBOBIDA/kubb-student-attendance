package controller

import (
	"context"
	"errors"
	"mysql/constant/share"
	"mysql/helper"
	"mysql/request"
	"mysql/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LeaveController struct {
	service service.LeaveRequestService
}

func NewLeaveController() LeaveController {
	return LeaveController{
		service: service.NewLeaveRequestService(),
	}
}

func (cr *LeaveController) GetLeaveRequest(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}
	filter := map[string]string{
		"name":       c.Query("name"),
		"class_id":   c.Query("class_id"),
		"status":     c.Query("status"),
		"subject_id": c.Query("subject_id"),
	}

	data, meta, err := cr.service.GetLeaveRequest(c.Request.Context(), userID, request.Pagination{
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

func (cr *LeaveController) CreateLeaveRequest(c *gin.Context) {
	var input request.LeaveRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}
	if err := cr.service.CreateLeaveRequest(c, userID, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "company created")
}

func (cr *LeaveController) AddNotPermission(c *gin.Context) {
	var input request.NotPermissionLeaveRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.AddNotPermission(c, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "company created")
}

func (cr *LeaveController) UpdateLeaveRequest(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}
	var input request.LeaveRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateLeaveRequest(c, id, userID, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "company Updated")
}

func (cr *LeaveController) VerifyLeaveRequest(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}
	if err := cr.service.VerifyLeaveRequest(c, id, userID); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "company Updated")
}

func (cr *LeaveController) DeleteLeaveRequest(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.DeleteLeaveRequest(c, id); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "Delted")
}

func (cr *LeaveController) GetNotPermissionLeave(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}
	filter := map[string]string{
		"name":          c.Query("name"),
		"class_id":      c.Query("class_id"),
		"major_id":      c.Query("major_id"),
		"shift_id":      c.Query("shift_id"),
		"generation_id": c.Query("generation_id"),
		"check_date":    c.Query("check_date"),
	}

	data, meta, err := cr.service.GetNotPermissionLeave(c.Request.Context(), userID, request.Pagination{
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
