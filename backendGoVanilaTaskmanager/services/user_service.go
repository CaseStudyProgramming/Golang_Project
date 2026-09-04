package services

import (
	"errors"
	"fmt"
	"taskmanager/models"
	"taskmanager/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type UserService struct {
	UserModel  models.UserModelInterface
	JWTManager *utils.JWTManager
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Timezone string `json:"timezone"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateTimezoneRequest struct {
	Timezone string `json:"timezone"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func NewUserService(userModel models.UserModelInterface, jwtManager *utils.JWTManager) *UserService {
	return &UserService{
		UserModel:  userModel,
		JWTManager: jwtManager,
	}
}

func (s *UserService) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.UserModel.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password with explicit cost factor for security
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Timezone:     req.Timezone,
	}

	// Set default timezone if not provided
	if user.Timezone == "" {
		user.Timezone = "UTC"
	}

	createdUser, err := s.UserModel.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token with timezone
	token, err := s.JWTManager.GenerateToken(createdUser.ID, createdUser.Email, createdUser.Timezone, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return &AuthResponse{
		Token: token,
		User:  createdUser,
	}, nil
}

func (s *UserService) Login(req *LoginRequest) (*AuthResponse, error) {
	// Get user by email
	user, err := s.UserModel.GetByEmail(req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token with timezone
	token, err := s.JWTManager.GenerateToken(user.ID, user.Email, user.Timezone, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *UserService) GetUserByID(userID int64) (*models.User, error) {
	return s.UserModel.GetByID(userID)
}

func (s *UserService) UpdateTimezone(userID int64, req *UpdateTimezoneRequest) (*models.User, error) {
	err := s.UserModel.UpdateTimezone(userID, req.Timezone)
	if err != nil {
		return nil, err
	}

	return s.UserModel.GetByID(userID)
}
