package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"taskmanager/models"
	"taskmanager/services"
	"testing"
)

// MockUserService is a mock implementation of UserServiceInterface for testing
type MockUserService struct {
	RegisterFunc       func(req *services.RegisterRequest) (*services.AuthResponse, error)
	LoginFunc          func(req *services.LoginRequest) (*services.AuthResponse, error)
	GetUserByIDFunc    func(userID int64) (*models.User, error)
	UpdateTimezoneFunc func(userID int64, req *services.UpdateTimezoneRequest) (*models.User, error)
}

// Ensure MockUserService implements UserServiceInterface
var _ UserServiceInterface = (*MockUserService)(nil)

func (m *MockUserService) Register(req *services.RegisterRequest) (*services.AuthResponse, error) {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(req)
	}
	return &services.AuthResponse{}, nil
}

func (m *MockUserService) Login(req *services.LoginRequest) (*services.AuthResponse, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(req)
	}
	return &services.AuthResponse{}, nil
}

func (m *MockUserService) GetUserByID(userID int64) (*models.User, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(userID)
	}
	return &models.User{}, nil
}

func (m *MockUserService) UpdateTimezone(userID int64, req *services.UpdateTimezoneRequest) (*models.User, error) {
	if m.UpdateTimezoneFunc != nil {
		return m.UpdateTimezoneFunc(userID, req)
	}
	return &models.User{}, nil
}

// Helper function to create request with context
func createAuthRequest(method, path string, body *bytes.Buffer) *http.Request {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body.Bytes())
	}
	req := httptest.NewRequest(method, path, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	// Add user_id to context for authentication
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	return req
}

func TestAuthController_Register(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockResponse   *services.AuthResponse
		mockError      error
		expectedStatus int
	}{
		{
			name: "successful registration",
			requestBody: services.RegisterRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
				Timezone: "UTC",
			},
			mockResponse: &services.AuthResponse{
				Token: "test-token",
				User: &models.User{
					ID:       1,
					Name:     "Test User",
					Email:    "test@example.com",
					Timezone: "UTC",
				},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "missing name",
			requestBody: services.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing email",
			requestBody: services.RegisterRequest{
				Name:     "Test User",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing password",
			requestBody: services.RegisterRequest{
				Name:  "Test User",
				Email: "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "password too short",
			requestBody: services.RegisterRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "12345",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid timezone",
			requestBody: services.RegisterRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
				Timezone: "Invalid/Timezone",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "user already exists",
			requestBody: services.RegisterRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockError:      services.ErrUserAlreadyExists,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockUserService{
				RegisterFunc: func(req *services.RegisterRequest) (*services.AuthResponse, error) {
					return tt.mockResponse, tt.mockError
				},
			}

			controller := NewAuthController(mockService)

			var body *bytes.Buffer
			if str, ok := tt.requestBody.(string); ok {
				body = bytes.NewBufferString(str)
			} else {
				jsonData, _ := json.Marshal(tt.requestBody)
				body = bytes.NewBuffer(jsonData)
			}

			req := createAuthRequest("POST", "/auth/register", body)
			w := httptest.NewRecorder()

			controller.Register(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestAuthController_Login(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockResponse   *services.AuthResponse
		mockError      error
		expectedStatus int
	}{
		{
			name: "successful login",
			requestBody: services.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockResponse: &services.AuthResponse{
				Token: "test-token",
				User: &models.User{
					ID:       1,
					Name:     "Test User",
					Email:    "test@example.com",
					Timezone: "UTC",
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "missing email",
			requestBody: services.LoginRequest{
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing password",
			requestBody: services.LoginRequest{
				Email: "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid credentials",
			requestBody: services.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockError:      services.ErrInvalidCredentials,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockUserService{
				LoginFunc: func(req *services.LoginRequest) (*services.AuthResponse, error) {
					return tt.mockResponse, tt.mockError
				},
			}

			controller := NewAuthController(mockService)

			var body *bytes.Buffer
			if str, ok := tt.requestBody.(string); ok {
				body = bytes.NewBufferString(str)
			} else {
				jsonData, _ := json.Marshal(tt.requestBody)
				body = bytes.NewBuffer(jsonData)
			}

			req := createAuthRequest("POST", "/auth/login", body)
			w := httptest.NewRecorder()

			controller.Login(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestAuthController_GetCurrentUser(t *testing.T) {
	tests := []struct {
		name           string
		mockUser       *models.User
		mockError      error
		expectedStatus int
	}{
		{
			name: "successful get current user",
			mockUser: &models.User{
				ID:       1,
				Name:     "Test User",
				Email:    "test@example.com",
				Timezone: "UTC",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user not found",
			mockError:      errors.New("user not found"),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockUserService{
				GetUserByIDFunc: func(userID int64) (*models.User, error) {
					return tt.mockUser, tt.mockError
				},
			}

			controller := NewAuthController(mockService)

			req := createAuthRequest("GET", "/auth/me", nil)
			w := httptest.NewRecorder()

			controller.GetCurrentUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestAuthController_Logout(t *testing.T) {
	controller := NewAuthController(nil)

	req := createAuthRequest("POST", "/auth/logout", nil)
	w := httptest.NewRecorder()

	controller.Logout(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthController_UpdateTimezone(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockUser       *models.User
		mockError      error
		expectedStatus int
	}{
		{
			name: "successful timezone update",
			requestBody: services.UpdateTimezoneRequest{
				Timezone: "America/New_York",
			},
			mockUser: &models.User{
				ID:       1,
				Name:     "Test User",
				Email:    "test@example.com",
				Timezone: "America/New_York",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "missing timezone",
			requestBody: services.UpdateTimezoneRequest{
				Timezone: "",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid timezone",
			requestBody: services.UpdateTimezoneRequest{
				Timezone: "Invalid/Timezone",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockUserService{
				UpdateTimezoneFunc: func(userID int64, req *services.UpdateTimezoneRequest) (*models.User, error) {
					return tt.mockUser, tt.mockError
				},
			}

			controller := NewAuthController(mockService)

			var body *bytes.Buffer
			if str, ok := tt.requestBody.(string); ok {
				body = bytes.NewBufferString(str)
			} else {
				jsonData, _ := json.Marshal(tt.requestBody)
				body = bytes.NewBuffer(jsonData)
			}

			req := createAuthRequest("PUT", "/auth/timezone", body)
			w := httptest.NewRecorder()

			controller.UpdateTimezone(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
