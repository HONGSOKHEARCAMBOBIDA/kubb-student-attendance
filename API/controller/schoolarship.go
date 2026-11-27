package controller

import (
	"mysql/constant/share"
	"mysql/helper"
	"mysql/request"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SchoolarshipController struct {
	service service.SchoolarshipService
}

func NewSchoolarshipController() SchoolarshipController {
	return SchoolarshipController{
		service: service.NewSchoolarshipService(),
	}
}

func (cr *SchoolarshipController) CreateSchoolarship(c *gin.Context) {
	var input request.SchoolarshipRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.CreateSchoolarship(c, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "Craeted")
}

func (cr *SchoolarshipController) UpdateSchoolarship(c *gin.Context) {
	id, ok := helper.GetParamID(c)
	if !ok {
		return
	}
	var input request.SchoolarshipRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateSchoolarship(c, id, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "Updated")
}
