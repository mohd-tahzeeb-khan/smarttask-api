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

	_ "github.com/smarttask/api/docs"
)

// @title           SmartTask AI API
// @version         1.0
// @description     AI-Powered Task & Productivity API — production-ready, startup-grade backend.

// @contact.name   Tahzeeb Khan
// @contact.email  tahzeeb@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer eyJhbG..."

func main() {
	config.Load()
	models.InitDB()

	gin.SetMode(config.App.GinMode)

	userRepo := repository.NewUserRepository(models.DB)
	taskRepo := repository.NewTaskRepository(models.DB)
	aiSvc := aiservice.NewAIService()

	authSvc := services.NewAuthService(userRepo)
	taskSvc := services.NewTaskService(taskRepo, aiSvc)

	authCtrl := controllers.NewAuthController(authSvc)
	taskCtrl := controllers.NewTaskController(taskSvc)

	r := gin.New()

	routes.Setup(r, authCtrl, taskCtrl, authSvc)

	addr := fmt.Sprintf(":%s", config.App.Port)
	log.Printf("🚀 SmartTask API running on http://localhost%s", addr)
	log.Printf("📖 Swagger UI: http://localhost%s/swagger/index.html", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
