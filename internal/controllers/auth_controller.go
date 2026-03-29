package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smarttask/api/internal/models"
	"github.com/smarttask/api/internal/services"
)

type AuthController struct {
	authSvc *services.AuthService
}

func NewAuthController(authSvc *services.AuthService) *AuthController {
	return &AuthController{authSvc: authSvc}
}

func (ctrl *AuthController) Signup(c *gin.Context) {
	var req models.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	resp, err := ctrl.authSvc.Signup(req)
	if err != nil {
		c.JSON(http.StatusConflict, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Account created successfully",
		Data:    resp,
	})
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	resp, err := ctrl.authSvc.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Login successful",
		Data:    resp,
	})
}

func (ctrl *AuthController) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	email, _ := c.Get("user_email")
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: gin.H{
			"user_id": userID,
			"email":   email,
		},
	})
}
