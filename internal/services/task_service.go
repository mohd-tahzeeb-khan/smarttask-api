package services

import (
	"errors"
	"time"

	aiservice "github.com/smarttask/api/pkg/ai"

	"github.com/smarttask/api/internal/models"
	"github.com/smarttask/api/internal/repository"
	"gorm.io/gorm"
)

type TaskService struct {
	taskRepo *repository.TaskRepository
	ai       *aiservice.AIService
}

func NewTaskService(taskRepo *repository.TaskRepository, ai *aiservice.AIService) *TaskService {
	return &TaskService{taskRepo: taskRepo, ai: ai}
}

func (s *TaskService) Create(userID uint, req models.CreateTaskRequest) (*models.Task, error) {
	task := &models.Task{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Deadline:    req.Deadline,
		Tags:        req.Tags,
		Status:      models.StatusPending,
	}

	if task.Priority == "" {
		task.Priority = models.PriorityMedium
	}

	analysis, err := s.ai.AnalyzeTask(task.Title, task.Description)
	if err == nil {
		task.AISuggestedPriority = analysis.Priority
		task.AIEstimatedHours = analysis.EstimatedTime
		task.AIAnalyzed = true
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) GetByID(id, userID uint) (*models.Task, error) {
	task, err := s.taskRepo.FindByID(id, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("task not found")
	}
	return task, err
}

func (s *TaskService) Update(id, userID uint, req models.UpdateTaskRequest) (*models.Task, error) {
	task, err := s.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Priority != "" {
		task.Priority = req.Priority
	}
	if req.Status != "" {
		task.Status = req.Status
		if req.Status == models.StatusDone {
			now := time.Now()
			task.CompletedAt = &now
		}
	}
	if req.Deadline != nil {
		task.Deadline = req.Deadline
	}
	if req.Tags != "" {
		task.Tags = req.Tags
	}

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) Delete(id, userID uint) error {
	return s.taskRepo.Delete(id, userID)
}

func (s *TaskService) List(userID uint, params models.TaskFilterParams) (*models.PaginatedTasks, error) {
	return s.taskRepo.FindAll(userID, params)
}

func (s *TaskService) AnalyzeTask(req models.AnalyzeTaskRequest) (*models.AIAnalysisResult, error) {
	return s.ai.AnalyzeTask(req.Title, req.Description)
}

func (s *TaskService) GetAnalytics(userID uint) (*models.Analytics, error) {
	return s.taskRepo.GetAnalytics(userID)
}
