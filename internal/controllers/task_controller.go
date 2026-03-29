package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smarttask/api/internal/middleware"
	"github.com/smarttask/api/internal/models"
	"github.com/smarttask/api/internal/services"
)

type TaskController struct {
	taskSvc *services.TaskService
}

func NewTaskController(taskSvc *services.TaskService) *TaskController {
	return &TaskController{taskSvc: taskSvc}
}

func (ctrl *TaskController) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	task, err := ctrl.taskSvc.Create(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Task created with AI analysis",
		Data:    task,
	})
}

func (ctrl *TaskController) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var params models.TaskFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	result, err := ctrl.taskSvc.List(userID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: result})
}

func (ctrl *TaskController) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := parseID(c, "id")
	if err != nil {
		return
	}

	task, err := ctrl.taskSvc.GetByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: task})
}

func (ctrl *TaskController) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := parseID(c, "id")
	if err != nil {
		return
	}

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	task, err := ctrl.taskSvc.Update(id, userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "Task updated", Data: task})
}

func (ctrl *TaskController) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := parseID(c, "id")
	if err != nil {
		return
	}

	if err := ctrl.taskSvc.Delete(id, userID); err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "Task deleted"})
}

func (ctrl *TaskController) AnalyzeTask(c *gin.Context) {
	var req models.AnalyzeTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	result, err := ctrl.taskSvc.AnalyzeTask(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "AI analysis complete",
		Data:    result,
	})
}

func (ctrl *TaskController) Analytics(c *gin.Context) {
	userID := middleware.GetUserID(c)
	analytics, err := ctrl.taskSvc.GetAnalytics(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: analytics})
}

func parseID(c *gin.Context, param string) (uint, error) {
	idStr := c.Param(param)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: "invalid ID"})
		return 0, err
	}
	return uint(id), nil
}
