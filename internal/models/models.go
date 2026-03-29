package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string `gorm:"not null" json:"name"`
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`
	Tasks        []Task `gorm:"foreignKey:UserID" json:"-"`
}

type TaskStatus string
type TaskPriority string

const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"

	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
)

type Task struct {
	gorm.Model
	UserID              uint         `gorm:"not null;index" json:"user_id"`
	Title               string       `gorm:"not null" json:"title"`
	Description         string       `json:"description"`
	Status              TaskStatus   `gorm:"default:'pending'" json:"status"`
	Priority            TaskPriority `gorm:"default:'medium'" json:"priority"`
	AISuggestedPriority TaskPriority `json:"ai_suggested_priority,omitempty"`
	AIEstimatedHours    float64      `json:"ai_estimated_hours,omitempty"`
	AIAnalyzed          bool         `gorm:"default:false" json:"ai_analyzed"`
	Deadline            *time.Time   `json:"deadline,omitempty"`
	CompletedAt         *time.Time   `json:"completed_at,omitempty"`
	Tags                string       `json:"tags,omitempty"`
}

type AIAnalysisResult struct {
	Priority      TaskPriority `json:"priority"`
	EstimatedTime float64      `json:"estimated_time_hours"`
	Reasoning     string       `json:"reasoning"`
	Confidence    float64      `json:"confidence"`
}

type SignupRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreateTaskRequest struct {
	Title       string       `json:"title" binding:"required,min=3"`
	Description string       `json:"description"`
	Priority    TaskPriority `json:"priority"`
	Deadline    *time.Time   `json:"deadline"`
	Tags        string       `json:"tags"`
}

type UpdateTaskRequest struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	Deadline    *time.Time   `json:"deadline"`
	Tags        string       `json:"tags"`
}

type AnalyzeTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type TaskFilterParams struct {
	Priority string `form:"priority"`
	Status   string `form:"status"`
	Deadline string `form:"deadline"`
	Page     int    `form:"page,default=1"`
	Limit    int    `form:"limit,default=10"`
	Sort     string `form:"sort,default=created_at"`
	Order    string `form:"order,default=desc"`
}

type Analytics struct {
	TotalTasks        int64             `json:"total_tasks"`
	CompletedTasks    int64             `json:"completed_tasks"`
	PendingTasks      int64             `json:"pending_tasks"`
	InProgressTasks   int64             `json:"in_progress_tasks"`
	ProductivityScore float64           `json:"productivity_score"`
	OverdueTasks      int64             `json:"overdue_tasks"`
	WeeklyInsights    []WeeklyInsight   `json:"weekly_insights"`
	PriorityBreakdown PriorityBreakdown `json:"priority_breakdown"`
}

type WeeklyInsight struct {
	Week      string `json:"week"`
	Completed int64  `json:"completed"`
	Created   int64  `json:"created"`
}

type PriorityBreakdown struct {
	High   int64 `json:"high"`
	Medium int64 `json:"medium"`
	Low    int64 `json:"low"`
}

type PaginatedTasks struct {
	Tasks      []Task `json:"tasks"`
	Total      int64  `json:"total"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalPages int    `json:"total_pages"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}