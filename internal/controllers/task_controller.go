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

// Create godoc
// @Summary      Create a new task
// @Description  Creates a task and auto-analyzes it with AI for priority and time estimate
// @Tags         Tasks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      models.CreateTaskRequest  true  "Task payload"
// @Success      201   {object}  models.APIResponse
// @Failure      400   {object}  models.APIResponse
// @Failure      401   {object}  models.APIResponse
// @Router       /tasks [post]
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

// List godoc
// @Summary      List all tasks
// @Description  Returns paginated tasks with optional filters
// @Tags         Tasks
// @Security     BearerAuth
// @Produce      json
// @Param        priority  query     string  false  "Filter by priority (low|medium|high)"
// @Param        status    query     string  false  "Filter by status (pending|in_progress|done)"
// @Param        deadline  query     string  false  "Filter by deadline (overdue|today|week)"
// @Param        page      query     int     false  "Page number (default 1)"
// @Param        limit     query     int     false  "Items per page (default 10)"
// @Param        sort      query     string  false  "Sort field (created_at|deadline|priority|status)"
// @Param        order     query     string  false  "Sort order (asc|desc)"
// @Success      200       {object}  models.APIResponse
// @Failure      401       {object}  models.APIResponse
// @Router       /tasks [get]
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

// Get godoc
// @Summary      Get a task by ID
// @Description  Returns a single task belonging to the authenticated user
// @Tags         Tasks
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  models.APIResponse
// @Failure      401  {object}  models.APIResponse
// @Failure      404  {object}  models.APIResponse
// @Router       /tasks/{id} [get]
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

// Update godoc
// @Summary      Update a task
// @Description  Update task fields including status, priority, title, deadline
// @Tags         Tasks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      int                       true  "Task ID"
// @Param        body  body      models.UpdateTaskRequest  true  "Update payload"
// @Success      200   {object}  models.APIResponse
// @Failure      400   {object}  models.APIResponse
// @Failure      401   {object}  models.APIResponse
// @Failure      404   {object}  models.APIResponse
// @Router       /tasks/{id} [put]
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

// Delete godoc
// @Summary      Delete a task
// @Description  Permanently deletes a task by ID
// @Tags         Tasks
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  models.APIResponse
// @Failure      401  {object}  models.APIResponse
// @Failure      404  {object}  models.APIResponse
// @Router       /tasks/{id} [delete]
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

// AnalyzeTask godoc
// @Summary      AI analyze a task
// @Description  Uses AI to suggest priority and estimate completion time for a task
// @Tags         AI
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      models.AnalyzeTaskRequest  true  "Task to analyze"
// @Success      200   {object}  models.APIResponse
// @Failure      400   {object}  models.APIResponse
// @Failure      401   {object}  models.APIResponse
// @Router       /ai/analyze-task [post]
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

// Analytics godoc
// @Summary      Get productivity analytics
// @Description  Returns full productivity dashboard — scores, weekly insights, priority breakdown
// @Tags         Analytics
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  models.APIResponse
// @Failure      401  {object}  models.APIResponse
// @Router       /analytics [get]
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
