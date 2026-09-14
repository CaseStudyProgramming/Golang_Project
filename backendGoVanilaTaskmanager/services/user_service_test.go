package services

import (
	"errors"
	"taskmanager/models"
	"taskmanager/utils"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// MockUserModel is a mock implementation of UserModelInterface for testing
type MockUserModel struct {
	CreateFunc         func(user *models.User) (*models.User, error)
	GetByIDFunc        func(id int64) (*models.User, error)
	GetByEmailFunc     func(email string) (*models.User, error)
	UpdateFunc         func(user *models.User) error
	UpdateTimezoneFunc func(userID int64, timezone string) error
	DeleteFunc         func(id int64) error
}

// Ensure MockUserModel implements UserModelInterface
var _ models.UserModelInterface = (*MockUserModel)(nil)

func (m *MockUserModel) Create(user *models.User) (*models.User, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(user)
	}
	return nil, nil
}

func (m *MockUserModel) GetByID(id int64) (*models.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *MockUserModel) GetByEmail(email string) (*models.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(email)
	}
	return nil, nil
}

func (m *MockUserModel) Update(user *models.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(user)
	}
	return nil
}

func (m *MockUserModel) UpdateTimezone(userID int64, timezone string) error {
	if m.UpdateTimezoneFunc != nil {
		return m.UpdateTimezoneFunc(userID, timezone)
	}
	return nil
}

func (m *MockUserModel) Delete(id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func TestUserService_Register(t *testing.T) {
	tests := []struct {
		name        string
		request     *RegisterRequest
		mockUser    *models.User
		mockError   error
		expectError error
	}{
		{
			name: "successful registration",
			request: &RegisterRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
				Timezone: "UTC",
			},
			mockUser: &models.User{
				ID:       1,
				Name:     "Test User",
				Email:    "test@example.com",
				Timezone: "UTC",
			},
		},
		{
			name: "user already exists",
			request: &RegisterRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockUser:    &models.User{ID: 1, Email: "test@example.com"},
			expectError: ErrUserAlreadyExists,
		},
		{
			name: "registration with default timezone",
			request: &RegisterRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockUser: &models.User{
				ID:       1,
				Name:     "Test User",
				Email:    "test@example.com",
				Timezone: "UTC",
			},
		},
		{
			name: "database error on create",
			request: &RegisterRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockError:   errors.New("database error"),
			expectError: errors.New("failed to create user: database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserModel := &MockUserModel{
				GetByEmailFunc: func(email string) (*models.User, error) {
					// For user already exists error, return the mock user
					if tt.expectError == ErrUserAlreadyExists {
						return tt.mockUser, nil
					}
					// For all other cases, return "not found" to simulate new user
					return nil, errors.New("not found")
				},
				CreateFunc: func(user *models.User) (*models.User, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockUser, nil
				},
			}

			mockJWTManager := utils.NewJWTManager("test-secret")
			service := NewUserService(mockUserModel, mockJWTManager)
			response, err := service.Register(tt.request)

			if tt.expectError != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectError)
				}
				if err != nil && !errors.Is(err, tt.expectError) && err.Error() != tt.expectError.Error() {
					t.Errorf("Expected error %v, got %v", tt.expectError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if response == nil {
					t.Errorf("Expected response, got nil")
				} else {
					if response.Token == "" {
						t.Errorf("Expected token in response")
					}
					if response.User == nil {
						t.Errorf("Expected user in response")
					}
				}
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 12)

	tests := []struct {
		name        string
		request     *LoginRequest
		mockUser    *models.User
		mockError   error
		expectError error
	}{
		{
			name: "successful login",
			request: &LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockUser: &models.User{
				ID:           1,
				Name:         "Test User",
				Email:        "test@example.com",
				PasswordHash: string(hashedPassword),
				Timezone:     "UTC",
			},
		},
		{
			name: "user not found",
			request: &LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			mockError:   errors.New("not found"),
			expectError: ErrInvalidCredentials,
		},
		{
			name: "invalid password",
			request: &LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockUser: &models.User{
				ID:           1,
				Email:        "test@example.com",
				PasswordHash: string(hashedPassword),
			},
			expectError: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserModel := &MockUserModel{
				GetByEmailFunc: func(email string) (*models.User, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockUser, nil
				},
			}

			mockJWTManager := utils.NewJWTManager("test-secret")
			service := NewUserService(mockUserModel, mockJWTManager)
			response, err := service.Login(tt.request)

			if tt.expectError != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectError)
				}
				if err != nil && !errors.Is(err, tt.expectError) && err.Error() != tt.expectError.Error() {
					t.Errorf("Expected error %v, got %v", tt.expectError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if response == nil {
					t.Errorf("Expected response, got nil")
				} else {
					if response.Token == "" {
						t.Errorf("Expected token in response")
					}
					if response.User == nil {
						t.Errorf("Expected user in response")
					}
				}
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		mockUser    *models.User
		mockError   error
		expectError error
	}{
		{
			name:   "successful get user",
			userID: 1,
			mockUser: &models.User{
				ID:    1,
				Name:  "Test User",
				Email: "test@example.com",
			},
		},
		{
			name:        "user not found",
			userID:      999,
			mockError:   errors.New("not found"),
			expectError: errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserModel := &MockUserModel{
				GetByIDFunc: func(id int64) (*models.User, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockUser, nil
				},
			}

			mockJWTManager := utils.NewJWTManager("test-secret")
			service := NewUserService(mockUserModel, mockJWTManager)
			user, err := service.GetUserByID(tt.userID)

			if tt.expectError != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectError)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if user == nil {
					t.Errorf("Expected user, got nil")
				}
			}
		})
	}
}

func TestUserService_UpdateTimezone(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		request     *UpdateTimezoneRequest
		mockUser    *models.User
		mockError   error
		expectError error
	}{
		{
			name:   "successful timezone update",
			userID: 1,
			request: &UpdateTimezoneRequest{
				Timezone: "America/New_York",
			},
			mockUser: &models.User{
				ID:       1,
				Name:     "Test User",
				Email:    "test@example.com",
				Timezone: "America/New_York",
			},
		},
		{
			name:   "database error",
			userID: 1,
			request: &UpdateTimezoneRequest{
				Timezone: "America/New_York",
			},
			mockError:   errors.New("database error"),
			expectError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserModel := &MockUserModel{
				UpdateTimezoneFunc: func(userID int64, timezone string) error {
					return tt.mockError
				},
				GetByIDFunc: func(id int64) (*models.User, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockUser, nil
				},
			}

			mockJWTManager := utils.NewJWTManager("test-secret")
			service := NewUserService(mockUserModel, mockJWTManager)
			user, err := service.UpdateTimezone(tt.userID, tt.request)

			if tt.expectError != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectError)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if user == nil {
					t.Errorf("Expected user, got nil")
				}
				if user != nil && user.Timezone != tt.request.Timezone {
					t.Errorf("Expected timezone %s, got %s", tt.request.Timezone, user.Timezone)
				}
			}
		})
	}
}
