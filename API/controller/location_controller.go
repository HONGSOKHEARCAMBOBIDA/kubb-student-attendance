package controller

import (
	"mysql/constant/share"
	"mysql/helper"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LocationController struct {
	service service.LocationService
}

func NewLocationController(s service.LocationService) *LocationController {
	return &LocationController{
		service: s,
	}
}

func (cr *LocationController) GetProvince(c *gin.Context) {
	data, err := cr.service.GetProvince(c.Request.Context())
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	share.RespondDate(c, http.StatusOK, data)
}

func (cr *LocationController) GetDistrict(c *gin.Context) {
	id, ok := helper.GetParamID(c)
	if !ok {
		return
	}

	data, err := cr.service.GetDistrict(c.Request.Context(), id)
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	share.RespondDate(c, http.StatusOK, data)
}

func (cr *LocationController) GetCommune(c *gin.Context) {
	id, ok := helper.GetParamID(c)
	if !ok {
		return
	}

	data, err := cr.service.GetCommune(c.Request.Context(), id)
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	share.RespondDate(c, http.StatusOK, data)
}

func (cr *LocationController) GetVillage(c *gin.Context) {
	id, ok := helper.GetParamID(c)
	if !ok {
		return
	}

	data, err := cr.service.GetVillage(c.Request.Context(), id)
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	share.RespondDate(c, http.StatusOK, data)
}
