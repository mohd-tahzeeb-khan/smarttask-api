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

// Signup godoc
// @Summary      Register a new user
// @Description  Create a new account with name, email and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.SignupRequest  true  "Signup payload"
// @Success      201   {object}  models.APIResponse
// @Failure      400   {object}  models.APIResponse
// @Failure      409   {object}  models.APIResponse
// @Router       /auth/signup [post]
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

// Login godoc
// @Summary      Login user
// @Description  Authenticate with email and password, returns JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.LoginRequest  true  "Login payload"
// @Success      200   {object}  models.APIResponse
// @Failure      400   {object}  models.APIResponse
// @Failure      401   {object}  models.APIResponse
// @Router       /auth/login [post]
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

// Me godoc
// @Summary      Get current user
// @Description  Returns the authenticated user's profile
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  models.APIResponse
// @Failure      401  {object}  models.APIResponse
// @Router       /auth/me [get]
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
