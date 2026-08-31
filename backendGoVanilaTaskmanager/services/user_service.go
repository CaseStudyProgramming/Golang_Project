package services

import (
	"errors"
	"taskmanager/models"
	"taskmanager/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists = errors.New("user with this email already exists")
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
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	createdUser, err := s.UserModel.Create(user)
	if err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := s.JWTManager.GenerateToken(createdUser.ID, createdUser.Email, 24*time.Hour)
	if err != nil {
		return nil, err
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

	// Generate JWT token
	token, err := s.JWTManager.GenerateToken(user.ID, user.Email, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *UserService) GetUserByID(userID int64) (*models.User, error) {
	return s.UserModel.GetByID(userID)
}
