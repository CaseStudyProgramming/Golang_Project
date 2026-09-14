package middlewares

import (
	"net/http"
	"net/http/httptest"
	"taskmanager/utils"
	"testing"
	"time"
)

func TestNewAuthMiddleware(t *testing.T) {
	jwtManager := utils.NewJWTManager("test-secret")
	middleware := NewAuthMiddleware(jwtManager)

	if middleware == nil {
		t.Fatal("NewAuthMiddleware returned nil")
	}

	if middleware.jwtManager != jwtManager {
		t.Error("NewAuthMiddleware did not set jwtManager correctly")
	}
}

func TestAuthMiddleware_Authenticate(t *testing.T) {
	jwtManager := utils.NewJWTManager("test-secret")
	middleware := NewAuthMiddleware(jwtManager)

	// Generate a valid token for testing
	validToken, _ := jwtManager.GenerateToken(1, "test@example.com", "America/New_York", time.Hour)

	tests := []struct {
		name           string
		authHeader     string
		expectStatus   int
		expectUserID   int64
		expectEmail    string
		expectTimezone string
	}{
		{
			name:           "valid bearer token",
			authHeader:     "Bearer " + validToken,
			expectStatus:   http.StatusOK,
			expectUserID:   1,
			expectEmail:    "test@example.com",
			expectTimezone: "America/New_York",
		},
		{
			name:         "missing authorization header",
			authHeader:   "",
			expectStatus: http.StatusUnauthorized,
		},
		{
			name:         "invalid authorization format",
			authHeader:   "InvalidFormat " + validToken,
			expectStatus: http.StatusUnauthorized,
		},
		{
			name:         "invalid token",
			authHeader:   "Bearer invalid.token.here",
			expectStatus: http.StatusUnauthorized,
		},
		{
			name:         "bearer without token",
			authHeader:   "Bearer",
			expectStatus: http.StatusUnauthorized,
		},
		{
			name:         "malformed bearer token",
			authHeader:   "Bearer token part2 part3",
			expectStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test handler that checks context values
			nextHandler := func(w http.ResponseWriter, r *http.Request) {
				if tt.expectStatus == http.StatusOK {
					userID := r.Context().Value("user_id")
					email := r.Context().Value("user_email")
					timezone := r.Context().Value("timezone")

					if userID == nil {
						t.Error("Expected user_id in context")
					} else if userID != tt.expectUserID {
						t.Errorf("Expected user_id %d, got %v", tt.expectUserID, userID)
					}

					if email == nil {
						t.Error("Expected user_email in context")
					} else if email != tt.expectEmail {
						t.Errorf("Expected email %s, got %v", tt.expectEmail, email)
					}

					if timezone == nil {
						t.Error("Expected timezone in context")
					} else if timezone != tt.expectTimezone {
						t.Errorf("Expected timezone %s, got %v", tt.expectTimezone, timezone)
					}

					w.WriteHeader(http.StatusOK)
				}
			}

			// Create request
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Call middleware
			handler := middleware.Authenticate(nextHandler)
			handler(w, req)

			// Check status code
			if w.Code != tt.expectStatus {
				t.Errorf("Expected status %d, got %d", tt.expectStatus, w.Code)
			}
		})
	}
}

func TestAuthMiddleware_Authenticate_DefaultTimezone(t *testing.T) {
	jwtManager := utils.NewJWTManager("test-secret")
	middleware := NewAuthMiddleware(jwtManager)

	// Generate token without timezone
	tokenWithoutTimezone, _ := jwtManager.GenerateToken(1, "test@example.com", "", time.Hour)

	nextHandler := func(w http.ResponseWriter, r *http.Request) {
		timezone := r.Context().Value("timezone")
		if timezone == nil {
			t.Error("Expected timezone in context")
		} else if timezone != "UTC" {
			t.Errorf("Expected default timezone UTC, got %v", timezone)
		}
		w.WriteHeader(http.StatusOK)
	}

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenWithoutTimezone)

	w := httptest.NewRecorder()
	handler := middleware.Authenticate(nextHandler)
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthMiddleware_Authenticate_ExpiredToken(t *testing.T) {
	jwtManager := utils.NewJWTManager("test-secret")
	middleware := NewAuthMiddleware(jwtManager)

	// Generate expired token
	expiredToken, _ := jwtManager.GenerateToken(1, "test@example.com", "UTC", -time.Hour)

	nextHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)

	w := httptest.NewRecorder()
	handler := middleware.Authenticate(nextHandler)
	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d for expired token, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_Authenticate_ContextValues(t *testing.T) {
	jwtManager := utils.NewJWTManager("test-secret")
	middleware := NewAuthMiddleware(jwtManager)

	validToken, _ := jwtManager.GenerateToken(123, "user@test.com", "Asia/Tokyo", time.Hour)

	nextHandler := func(w http.ResponseWriter, r *http.Request) {
		// Check that all expected context values are set
		userID := r.Context().Value("user_id")
		email := r.Context().Value("user_email")
		timezone := r.Context().Value("timezone")

		if userID == nil {
			t.Error("Expected user_id in context")
		} else if userID != int64(123) {
			t.Errorf("Expected user_id 123, got %v", userID)
		}

		if email == nil {
			t.Error("Expected user_email in context")
		} else if email != "user@test.com" {
			t.Errorf("Expected email user@test.com, got %v", email)
		}

		if timezone == nil {
			t.Error("Expected timezone in context")
		} else if timezone != "Asia/Tokyo" {
			t.Errorf("Expected timezone Asia/Tokyo, got %v", timezone)
		}

		w.WriteHeader(http.StatusOK)
	}

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)

	w := httptest.NewRecorder()
	handler := middleware.Authenticate(nextHandler)
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}
