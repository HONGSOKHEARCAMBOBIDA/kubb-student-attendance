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

// handler/score.go
func (cr *ScoreController) ImportScoreExcel(c *gin.Context) {
	var req request.ImportScoreExcelRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "file is required"})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(400, gin.H{"error": "cannot open file"})
		return
	}
	defer file.Close()

	result, err := cr.service.ImportScoreFromExcell(c.Request.Context(), req, file)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": result})
}

func (cr *ScoreController) GetScore(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}
	filter := map[string]string{
		"name":          c.Query("name"),
		"class_id":      c.Query("class_id"),
		"generation_id": c.Query("generation_id"),
		"major_id":      c.Query("major_id"),
		"programme_id":  c.Query("programme_id"),
		"subject_id":    c.Query("subject_id"),
		"code":          c.Query("code"),
		"year":          c.Query("year"),
		"semester":      c.Query("semester"),
	}

	data, meta, err := cr.service.GetScore(c.Request.Context(), userID, request.Pagination{
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
