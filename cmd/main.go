package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/smarttask/api/internal/config"
	"github.com/smarttask/api/internal/controllers"
	"github.com/smarttask/api/internal/models"
	"github.com/smarttask/api/internal/repository"
	"github.com/smarttask/api/internal/routes"
	"github.com/smarttask/api/internal/services"
	aiservice "github.com/smarttask/api/pkg/ai"
)

func main() {
	// Load config
	config.Load()

	// Init DB
	models.InitDB()

	// Gin mode
	gin.SetMode(config.App.GinMode)

	// Wire dependencies
	userRepo := repository.NewUserRepository(models.DB)
	taskRepo := repository.NewTaskRepository(models.DB)
	aiSvc := aiservice.NewAIService()

	authSvc := services.NewAuthService(userRepo)
	taskSvc := services.NewTaskService(taskRepo, aiSvc)

	authCtrl := controllers.NewAuthController(authSvc)
	taskCtrl := controllers.NewTaskController(taskSvc)

	// Create Gin engine
	r := gin.New()

	// Register routes
	routes.Setup(r, authCtrl, taskCtrl, authSvc)

	addr := fmt.Sprintf(":%s", config.App.Port)
	log.Printf("🚀 SmartTask API running on http://localhost%s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
