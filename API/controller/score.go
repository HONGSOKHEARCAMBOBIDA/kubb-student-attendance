package controller

import (
	"mysql/constant/share"
	"mysql/request"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScoreController struct {
	service service.ScoreService
}

func NewScoreController() ScoreController {
	return ScoreController{
		service: service.NewScoreService(),
	}
}

func (cr *ScoreController) CreateScore(c *gin.Context) {
	var input request.CreateScoreRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.CreateScore(c, input); err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponseSuccess(c, http.StatusOK, "score created")
}

func (cr *ScoreController) GetGradeComponent(c *gin.Context) {
	data, err := cr.service.GetGradeComponent(c)
	if err != nil {
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
