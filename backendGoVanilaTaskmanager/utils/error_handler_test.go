package utils

import (
	"errors"
	"fmt"
	"testing"
)

func TestSanitizeError(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected string
	}{
		{
			name:     "password in connection string",
			input:    errors.New("connection failed: password=secret123 host=localhost"),
			expected: "connection failed: password=*** host=***",
		},
		{
			name:     "file path",
			input:    errors.New("cannot open file C:\\Users\\admin\\secret.txt"),
			expected: "cannot open file C:\\Users\\***\\secret.txt",
		},
		{
			name:     "SQL error",
			input:    errors.New("SQLSTATE 23505 near 'SELECT * FROM users'"),
			expected: "SQLSTATE=*** near '***'",
		},
		{
			name:     "IP address",
			input:    errors.New("connection from 192.168.1.1 failed"),
			expected: "connection from ***.***.***.*** failed",
		},
		{
			name:     "email address",
			input:    errors.New("user admin@example.com not found"),
			expected: "user ***@***.*** not found",
		},
		{
			name:     "API key",
			input:    errors.New("invalid api_key=abc123xyz789"),
			expected: "invalid api_key=***",
		},
		{
			name:     "simple error",
			input:    errors.New("invalid input"),
			expected: "invalid input",
		},
		{
			name:     "nil error",
			input:    nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeError(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetUserMessage(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected string
	}{
		{
			name:     "public error",
			input:    NewPublicError("Invalid input", 400),
			expected: "Invalid input",
		},
		{
			name:     "internal error",
			input:    NewInternalError("Database failed", errors.New("connection error"), 500),
			expected: "An internal error occurred. Please try again later.",
		},
		{
			name:     "SQL error should be generic",
			input:    errors.New("SQL connection failed: database=prod host=localhost"),
			expected: "An internal error occurred. Please try again later.",
		},
		{
			name:     "simple wrapped error",
			input:    fmt.Errorf("failed to create task: %w", errors.New("validation error")),
			expected: "failed to create task: validation error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetUserMessage(tt.input)
			if result != tt.expected {
				t.Errorf("GetUserMessage() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHandleAPIError(t *testing.T) {
	tests := []struct {
		name         string
		input        error
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "public error",
			input:        NewPublicError("Invalid input", 400),
			expectedCode: 400,
			expectedMsg:  "Invalid input",
		},
		{
			name:         "internal error",
			input:        NewInternalError("Database failed", errors.New("connection error"), 500),
			expectedCode: 500,
			expectedMsg:  "An internal error occurred. Please try again later.",
		},
		{
			name:         "regular error",
			input:        errors.New("some error"),
			expectedCode: 500,
			expectedMsg:  "An internal error occurred. Please try again later.",
		},
		{
			name:         "nil error",
			input:        nil,
			expectedCode: 200,
			expectedMsg:  "Success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, msg := HandleAPIError(tt.input)
			if code != tt.expectedCode {
				t.Errorf("HandleAPIError() code = %v, want %v", code, tt.expectedCode)
			}
			if msg != tt.expectedMsg {
				t.Errorf("HandleAPIError() msg = %v, want %v", msg, tt.expectedMsg)
			}
		})
	}
}

func TestIsNotFoundError(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected bool
	}{
		{
			name:     "not found error",
			input:    ErrNotFound,
			expected: true,
		},
		{
			name:     "other error",
			input:    errors.New("some error"),
			expected: false,
		},
		{
			name:     "wrapped not found",
			input:    NewPublicError("Resource not found", 404),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotFoundError(tt.input)
			if result != tt.expected {
				t.Errorf("IsNotFoundError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsUnauthorizedError(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected bool
	}{
		{
			name:     "unauthorized error",
			input:    ErrUnauthorized,
			expected: true,
		},
		{
			name:     "invalid credentials",
			input:    ErrInvalidCredentials,
			expected: true,
		},
		{
			name:     "other error",
			input:    errors.New("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsUnauthorizedError(tt.input)
			if result != tt.expected {
				t.Errorf("IsUnauthorizedError() = %v, want %v", result, tt.expected)
			}
		})
	}
}
