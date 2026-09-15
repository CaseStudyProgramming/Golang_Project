package utils

import (
	"testing"
	"time"
)

func TestNewJWTManager(t *testing.T) {
	secretKey := "test-secret-key"
	manager := NewJWTManager(secretKey)

	if manager == nil {
		t.Fatal("NewJWTManager returned nil")
	}

	if manager.secretKey != secretKey {
		t.Error("NewJWTManager did not set secretKey correctly")
	}
}

func TestJWTManager_GenerateToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key")

	tests := []struct {
		name        string
		userID      int64
		email       string
		timezone    string
		expiration  time.Duration
		expectError bool
	}{
		{
			name:        "successful token generation",
			userID:      1,
			email:       "test@example.com",
			timezone:    "UTC",
			expiration:  time.Hour,
			expectError: false,
		},
		{
			name:        "token generation without timezone",
			userID:      1,
			email:       "test@example.com",
			timezone:    "",
			expiration:  time.Hour,
			expectError: false,
		},
		{
			name:        "token with different expiration",
			userID:      1,
			email:       "test@example.com",
			timezone:    "America/New_York",
			expiration:  24 * time.Hour,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := manager.GenerateToken(tt.userID, tt.email, tt.timezone, tt.expiration)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if token == "" {
					t.Error("Expected token but got empty string")
				}
			}
		})
	}
}

func TestJWTManager_ValidateToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key")

	// Generate a valid token for testing
	validToken, _ := manager.GenerateToken(1, "test@example.com", "UTC", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		expectError bool
		errorType   error
	}{
		{
			name:        "valid token",
			tokenString: validToken,
			expectError: false,
		},
		{
			name:        "invalid token format",
			tokenString: "invalid.token.format",
			expectError: true,
		},
		{
			name:        "empty token",
			tokenString: "",
			expectError: true,
		},
		{
			name:        "malformed token",
			tokenString: "not.a.valid.jwt",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := manager.ValidateToken(tt.tokenString)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if tt.errorType != nil && err != tt.errorType {
					// Check if it's the expected error type or a JWT validation error
					if err != ErrInvalidToken && err != ErrExpiredToken {
						t.Errorf("Expected error type %v, got %v", tt.errorType, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if claims == nil {
					t.Error("Expected claims but got nil")
				}
				if claims != nil {
					if claims.UserID != 1 {
						t.Errorf("Expected UserID 1, got %d", claims.UserID)
					}
					if claims.Email != "test@example.com" {
						t.Errorf("Expected email test@example.com, got %s", claims.Email)
					}
					if claims.Timezone != "UTC" {
						t.Errorf("Expected timezone UTC, got %s", claims.Timezone)
					}
				}
			}
		})
	}
}

func TestJWTManager_ValidateToken_WrongSecret(t *testing.T) {
	manager1 := NewJWTManager("secret1")
	manager2 := NewJWTManager("secret2")

	// Generate token with first manager
	token, _ := manager1.GenerateToken(1, "test@example.com", "UTC", time.Hour)

	// Try to validate with second manager (different secret)
	_, err := manager2.ValidateToken(token)
	if err == nil {
		t.Error("Expected error when validating token with different secret")
	}
}

func TestJWTManager_ValidateToken_ExpiredToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key")

	// Generate an expired token
	expiredToken, _ := manager.GenerateToken(1, "test@example.com", "UTC", -time.Hour)

	_, err := manager.ValidateToken(expiredToken)
	if err == nil {
		t.Error("Expected error for expired token")
	}
}

func TestJWTManager_GenerateAndValidate(t *testing.T) {
	manager := NewJWTManager("test-secret-key")

	userID := int64(123)
	email := "user@example.com"
	timezone := "America/Los_Angeles"
	expiration := 2 * time.Hour

	// Generate token
	token, err := manager.GenerateToken(userID, email, timezone, expiration)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate token
	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Verify claims
	if claims.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("Expected email %s, got %s", email, claims.Email)
	}
	if claims.Timezone != timezone {
		t.Errorf("Expected timezone %s, got %s", timezone, claims.Timezone)
	}
}
