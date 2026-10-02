package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TranscriptController struct {
	service service.TranscriptService
}

func NewTranscriptController() TranscriptController {
	return TranscriptController{
		service: service.NewTranscriptService(),
	}
}

func (cr *TranscriptController) Transcript(c *gin.Context) {
	idparam := c.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	data, err := cr.service.Transcript(c, id)
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
	}
	share.RespondDate(c, http.StatusOK, data)
}
