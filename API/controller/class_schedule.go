package controller

// NOTE: this file is a best-effort sketch since your existing controller
// layer wasn't shown to me. It assumes gin + a response.Error/response.Success
// style helper similar to the rest of the codebase (utils.ToggleStatus,
// apperror.New, etc.). Wire it up to whatever your actual controllers use
// for binding/responding, and confirm apperror.CodeBadRequest exists
// alongside CodeNotFound/CodeInternal/CodeConflict used elsewhere.

import (
	"net/http"
	"strconv"

	"mysql/request"
	"mysql/service"

	"github.com/gin-gonic/gin"
)

type ClassScheduleController struct {
	service service.ClassScheduleService
}

func NewClassScheduleController() ClassScheduleController {
	return ClassScheduleController{
		service: service.NewClassScheduleService(),
	}
}

// GET /class/:id/schedule
func (ctl *ClassScheduleController) GetByClass(c *gin.Context) {
	classID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid class id"})
		return
	}

	data, err := ctl.service.GetByClass(c.Request.Context(), classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// GET /class/:id/schedule/subjects
func (ctl *ClassScheduleController) GetAvailableSubjects(c *gin.Context) {
	classID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid class id"})
		return
	}

	data, err := ctl.service.GetAvailableSubjects(c.Request.Context(), classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// POST /class/:id/schedule
func (ctl *ClassScheduleController) Create(c *gin.Context) {
	classID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid class id"})
		return
	}

	var input request.ClassScheduleRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctl.service.Create(c.Request.Context(), classID, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

// PATCH /class-schedule/:id/toggle
func (ctl *ClassScheduleController) Toggle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := ctl.service.Toggle(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// Example route registration — place alongside your other class routes:
//
//   scheduleSvc := service.NewClassScheduleService()
//   scheduleCtl := controller.NewClassScheduleController(scheduleSvc)
//
//   classGroup.GET("/:id/schedule", scheduleCtl.GetByClass)
//   classGroup.GET("/:id/schedule/subjects", scheduleCtl.GetAvailableSubjects)
//   classGroup.POST("/:id/schedule", scheduleCtl.Create)
//   classGroup.PATCH("/schedule/:id/toggle", scheduleCtl.Toggle)
