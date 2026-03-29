package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smarttask/api/internal/controllers"
	"github.com/smarttask/api/internal/middleware"
	"github.com/smarttask/api/internal/models"
	"github.com/smarttask/api/internal/services"
)

func Setup(
	r *gin.Engine,
	authCtrl *controllers.AuthController,
	taskCtrl *controllers.TaskController,
	authSvc *services.AuthService,
) {
	rl := middleware.NewRateLimiter(100, time.Minute)

	r.Use(middleware.Logger())
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.RateLimit(rl))
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.APIResponse{
			Success: true,
			Message: "SmartTask API is running 🚀",
			Data: gin.H{
				"version": "1.0.0",
				"status":  "healthy",
				"time":    time.Now().UTC(),
			},
		})
	})

	api := r.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/signup", authCtrl.Signup)
		auth.POST("/login", authCtrl.Login)
		auth.GET("/me", middleware.AuthMiddleware(authSvc), authCtrl.Me)
	}

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(authSvc))
	{
		tasks := protected.Group("/tasks")
		{
			tasks.POST("", taskCtrl.Create)
			tasks.GET("", taskCtrl.List)
			tasks.GET("/:id", taskCtrl.Get)
			tasks.PUT("/:id", taskCtrl.Update)
			tasks.DELETE("/:id", taskCtrl.Delete)
		}

		ai := protected.Group("/ai")
		{
			ai.POST("/analyze-task", taskCtrl.AnalyzeTask)
		}

		protected.GET("/analytics", taskCtrl.Analytics)
	}
}
