package controller

import (
	"context"
	"errors"
	"log"

	//"log"
	"mysql/constant/share"
	"mysql/helper"
	"mysql/request"

	//"mysql/request"
	"mysql/service"
	"net/http"

	//"strconv"

	"github.com/gin-gonic/gin"
)

type AttendanceController struct {
	service service.AttendanceService
}

func NewAttendanceController() AttendanceController {
	return AttendanceController{
		service: service.NewAttendanceService(),
	}
}

func (cr *AttendanceController) CreateAttendance(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "please login")
		return
	}
	var input request.AttendanceRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	//log.Printf("CreateAttendance - UserID: %d, Input: %+v", userID, input)
	if err := cr.service.CreateAttendance(c, userID, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "attendance created")
}

// func (cr *AttendanceController) GetAttendance(c *gin.Context) {
// 	userID, ok := helper.GetUserID(c)
// 	if !ok {
// 		share.ResponseError(c, http.StatusUnauthorized, "please login")
// 		return
// 	}
// 	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
// 	if page < 1 {
// 		page = 1
// 	}
// 	if pageSize < 1 || pageSize > 100 {
// 		pageSize = 10
// 	}

// 	filter := map[string]string{
// 		"name":       c.Query("name"),
// 		"company_id": c.Query("company_id"),
// 		"role_id":    c.Query("role_id"),
// 		"check_date": c.Query("check_date"),
// 	}

// 	attendances, metadata, err := cr.service.GetAttendance(c, userID, request.Pagination{
// 		Page:     page,
// 		PageSize: pageSize,
// 	}, filter)

// 	if err != nil {
// 		share.ResponseError(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success":    true,
// 		"data":       attendances,
// 		"pagination": metadata,
// 	})
// }

func (cr *AttendanceController) GetAttendanceDraft(c *gin.Context) {

	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "please login")
		return
	}
	data, err := cr.service.GetAttendanceDraft(c, userID)
	if err != nil {
		log.Printf("Errr %w", err)
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *AttendanceController) GetAttendancePDF(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}
	filter := map[string]string{
		"name":       c.Query("name"),
		"class_id":   c.Query("class_id"),
		"check_date": c.Query("check_date"),
	}

	data, meta, err := cr.service.GetAttendancePDF(c.Request.Context(), userID, request.Pagination{
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

func (h *AttendanceController) GetAttendanceReport(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}
	filter := map[string]string{
		"name":      c.Query("name"),
		"class_id":  c.Query("class_id"),
		"date_from": c.Query("date_from"),
		"date_to":   c.Query("date_to"),
	}
	report, err := h.service.GetAttendanceReport(c.Request.Context(), userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": report})
}

// func (cr *AttendanceController) DeleteAttendance(c *gin.Context) {
// 	idparam := c.Param("id")
// 	id, err := strconv.Atoi(idparam)
// 	if err != nil {
// 		share.ResponseError(c, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	if err := cr.service.DeleteAttendance(c, id); err != nil {
// 		share.ResponseError(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	share.ResponseSuccess(c, http.StatusOK, "Deleted")
// }
