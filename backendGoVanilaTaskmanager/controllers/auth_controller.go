package controllers

import (
	"encoding/json"
	"net/http"
	"taskmanager/models"
	"taskmanager/services"
	"taskmanager/utils"
)

type AuthController struct {
	userService *services.UserService
}

type UserServiceInterface interface {
	Register(req *services.RegisterRequest) (*services.AuthResponse, error)
	Login(req *services.LoginRequest) (*services.AuthResponse, error)
	GetUserByID(userID int64) (*models.User, error)
}

func NewAuthController(userService *services.UserService) *AuthController {
	return &AuthController{userService: userService}
}

// POST /auth/register
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req services.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Name == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Name is required")
		return
	}
	if req.Email == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Email is required")
		return
	}
	if req.Password == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Password is required")
		return
	}
	if len(req.Password) < 6 {
		utils.ErrorResponse(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	response, err := c.userService.Register(&req)
	if err != nil {
		if err == services.ErrUserAlreadyExists {
			utils.ErrorResponse(w, http.StatusConflict, err.Error())
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "User registered successfully", response)
}

// POST /auth/login
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req services.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Email == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Email is required")
		return
	}
	if req.Password == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Password is required")
		return
	}

	response, err := c.userService.Login(&req)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			utils.ErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Login successful", response)
}

// GET /auth/me
func (c *AuthController) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	user, err := c.userService.GetUserByID(userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "User retrieved successfully", user)
}

// POST /auth/logout
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	// Since we're using stateless JWT tokens, logout is handled client-side
	// by removing the token. This endpoint is mainly for consistency.
	utils.SuccessResponse(w, http.StatusOK, "Logout successful", nil)
}
