package utils

import (
	"errors"
	"testing"
)

func TestWrapInternalError(t *testing.T) {
	originalErr := errors.New("database connection failed")
	wrappedErr := WrapInternalError("Failed to connect to database", originalErr)

	if wrappedErr == nil {
		t.Fatal("WrapInternalError returned nil")
	}

	if wrappedErr.Message != "Failed to connect to database" {
		t.Errorf("Expected message 'Failed to connect to database', got '%s'", wrappedErr.Message)
	}

	if wrappedErr.StatusCode != 500 {
		t.Errorf("Expected status code 500, got %d", wrappedErr.StatusCode)
	}

	if wrappedErr.Err == nil {
		t.Error("Expected wrapped error to contain original error")
	}

	if !errors.Is(wrappedErr.Err, originalErr) {
		t.Error("Expected wrapped error to contain original error")
	}
}

func TestWrapInternalErrorWithStatus(t *testing.T) {
	originalErr := errors.New("resource not found")
	wrappedErr := WrapInternalErrorWithStatus("User not found", originalErr, 404)

	if wrappedErr == nil {
		t.Fatal("WrapInternalErrorWithStatus returned nil")
	}

	if wrappedErr.Message != "User not found" {
		t.Errorf("Expected message 'User not found', got '%s'", wrappedErr.Message)
	}

	if wrappedErr.StatusCode != 404 {
		t.Errorf("Expected status code 404, got %d", wrappedErr.StatusCode)
	}

	if wrappedErr.Err == nil {
		t.Error("Expected wrapped error to contain original error")
	}
}

func TestCreatePublicError(t *testing.T) {
	publicErr := CreatePublicError("Custom public error", 400)

	if publicErr == nil {
		t.Fatal("CreatePublicError returned nil")
	}

	if publicErr.Message != "Custom public error" {
		t.Errorf("Expected message 'Custom public error', got '%s'", publicErr.Message)
	}

	if publicErr.StatusCode != 400 {
		t.Errorf("Expected status code 400, got %d", publicErr.StatusCode)
	}

	if publicErr.Err != nil {
		t.Error("Expected public error to not contain internal error")
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name         string
		err          *AppError
		expectedMsg  string
		expectedCode int
	}{
		{"ErrInvalidInput", ErrInvalidInput, "Invalid input provided", 400},
		{"ErrMissingRequired", ErrMissingRequired, "Required field is missing", 400},
		{"ErrInvalidFormat", ErrInvalidFormat, "Invalid format", 400},
		{"ErrUnauthorized", ErrUnauthorized, "Unauthorized access", 401},
		{"ErrInvalidCredentials", ErrInvalidCredentials, "Invalid email or password", 401},
		{"ErrTokenExpired", ErrTokenExpired, "Token has expired", 401},
		{"ErrForbidden", ErrForbidden, "You do not have permission to perform this action", 403},
		{"ErrNotFound", ErrNotFound, "Resource not found", 404},
		{"ErrConflict", ErrConflict, "Resource already exists", 409},
		{"ErrInvalidOperation", ErrInvalidOperation, "Invalid operation", 400},
		{"ErrInvalidState", ErrInvalidState, "Invalid state for this operation", 400},
		{"ErrInternal", ErrInternal, "An internal error occurred. Please try again later", 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Message != tt.expectedMsg {
				t.Errorf("Expected message '%s', got '%s'", tt.expectedMsg, tt.err.Message)
			}
			if tt.err.StatusCode != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, tt.err.StatusCode)
			}
		})
	}
}
