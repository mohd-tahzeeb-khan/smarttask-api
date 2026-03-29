package repository

import (
	"math"
	"time"

	"github.com/smarttask/api/internal/models"
	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepository) FindByID(id, userID uint) (*models.Task, error) {
	var task models.Task
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) Update(task *models.Task) error {
	return r.db.Save(task).Error
}

func (r *TaskRepository) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Task{}).Error
}

func (r *TaskRepository) FindAll(userID uint, params models.TaskFilterParams) (*models.PaginatedTasks, error) {
	query := r.db.Model(&models.Task{}).Where("user_id = ?", userID)

	if params.Priority != "" {
		query = query.Where("priority = ?", params.Priority)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	now := time.Now()
	switch params.Deadline {
	case "overdue":
		query = query.Where("deadline < ? AND status != ?", now, models.StatusDone)
	case "today":
		endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
		query = query.Where("deadline BETWEEN ? AND ?", now, endOfDay)
	case "week":
		endOfWeek := now.AddDate(0, 0, 7)
		query = query.Where("deadline BETWEEN ? AND ?", now, endOfWeek)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	allowedSorts := map[string]bool{"created_at": true, "deadline": true, "priority": true, "status": true}
	sortCol := "created_at"
	if allowedSorts[params.Sort] {
		sortCol = params.Sort
	}
	order := "desc"
	if params.Order == "asc" {
		order = "asc"
	}
	query = query.Order(sortCol + " " + order)

	if params.Limit <= 0 || params.Limit > 100 {
		params.Limit = 10
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	offset := (params.Page - 1) * params.Limit

	var tasks []models.Task
	if err := query.Offset(offset).Limit(params.Limit).Find(&tasks).Error; err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	return &models.PaginatedTasks{
		Tasks:      tasks,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

func (r *TaskRepository) GetAnalytics(userID uint) (*models.Analytics, error) {
	var analytics models.Analytics
	now := time.Now()

	r.db.Model(&models.Task{}).Where("user_id = ?", userID).Count(&analytics.TotalTasks)
	r.db.Model(&models.Task{}).Where("user_id = ? AND status = ?", userID, models.StatusDone).Count(&analytics.CompletedTasks)
	r.db.Model(&models.Task{}).Where("user_id = ? AND status = ?", userID, models.StatusPending).Count(&analytics.PendingTasks)
	r.db.Model(&models.Task{}).Where("user_id = ? AND status = ?", userID, models.StatusInProgress).Count(&analytics.InProgressTasks)
	r.db.Model(&models.Task{}).Where("user_id = ? AND deadline < ? AND status != ?", userID, now, models.StatusDone).Count(&analytics.OverdueTasks)

	r.db.Model(&models.Task{}).Where("user_id = ? AND priority = ?", userID, models.PriorityHigh).Count(&analytics.PriorityBreakdown.High)
	r.db.Model(&models.Task{}).Where("user_id = ? AND priority = ?", userID, models.PriorityMedium).Count(&analytics.PriorityBreakdown.Medium)
	r.db.Model(&models.Task{}).Where("user_id = ? AND priority = ?", userID, models.PriorityLow).Count(&analytics.PriorityBreakdown.Low)

	if analytics.TotalTasks > 0 {
		base := float64(analytics.CompletedTasks) / float64(analytics.TotalTasks) * 100
		penalty := float64(analytics.OverdueTasks) * 5
		analytics.ProductivityScore = math.Max(0, math.Min(100, base-penalty))
	}

	for i := 3; i >= 0; i-- {
		weekStart := now.AddDate(0, 0, -7*(i+1))
		weekEnd := now.AddDate(0, 0, -7*i)

		var completed, created int64
		r.db.Model(&models.Task{}).
			Where("user_id = ? AND status = ? AND completed_at BETWEEN ? AND ?", userID, models.StatusDone, weekStart, weekEnd).
			Count(&completed)
		r.db.Model(&models.Task{}).
			Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, weekStart, weekEnd).
			Count(&created)

		analytics.WeeklyInsights = append(analytics.WeeklyInsights, models.WeeklyInsight{
			Week:      weekStart.Format("Jan 02"),
			Completed: completed,
			Created:   created,
		})
	}

	return &analytics, nil
}
