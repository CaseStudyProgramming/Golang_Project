package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSuccessResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
		data       interface{}
	}{
		{
			name:       "successful response with data",
			statusCode: StatusCodeOK,
			message:    "Success",
			data:       map[string]string{"key": "value"},
		},
		{
			name:       "successful response with nil data",
			statusCode: StatusCodeOK,
			message:    "Success",
			data:       nil,
		},
		{
			name:       "created response",
			statusCode: StatusCodeCreated,
			message:    "Resource created",
			data:       map[string]int64{"id": 123},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			SuccessResponse(w, tt.statusCode, tt.message, tt.data)

			// Check status code
			if w.Code != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, w.Code)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", contentType)
			}

			// Check response body
			var response APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Status != "success" {
				t.Errorf("Expected status 'success', got '%s'", response.Status)
			}
			if response.Message != tt.message {
				t.Errorf("Expected message '%s', got '%s'", tt.message, response.Message)
			}
		})
	}
}

func TestSuccessResponseWithTimezone(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
		data       interface{}
		timezone   string
	}{
		{
			name:       "response with timezone",
			statusCode: StatusCodeOK,
			message:    "Success",
			data:       map[string]string{"key": "value"},
			timezone:   "America/New_York",
		},
		{
			name:       "response with empty timezone",
			statusCode: StatusCodeOK,
			message:    "Success",
			data:       map[string]string{"key": "value"},
			timezone:   "",
		},
		{
			name:       "response with UTC timezone",
			statusCode: StatusCodeOK,
			message:    "Success",
			data:       map[string]string{"key": "value"},
			timezone:   "UTC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			SuccessResponseWithTimezone(w, tt.statusCode, tt.message, tt.data, tt.timezone)

			// Check status code
			if w.Code != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, w.Code)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", contentType)
			}

			// Check response body
			var response TimezoneAwareResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Status != "success" {
				t.Errorf("Expected status 'success', got '%s'", response.Status)
			}
			if response.Message != tt.message {
				t.Errorf("Expected message '%s', got '%s'", tt.message, response.Message)
			}
			if response.Timezone != tt.timezone {
				t.Errorf("Expected timezone '%s', got '%s'", tt.timezone, response.Timezone)
			}
		})
	}
}

func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
	}{
		{
			name:       "bad request error",
			statusCode: StatusCodeBadRequest,
			message:    "Invalid input",
		},
		{
			name:       "not found error",
			statusCode: StatusCodeNotFound,
			message:    "Resource not found",
		},
		{
			name:       "internal server error",
			statusCode: StatusCodeInternalServerError,
			message:    "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ErrorResponse(w, tt.statusCode, tt.message)

			// Check status code
			if w.Code != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, w.Code)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", contentType)
			}

			// Check response body
			var response APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Status != "error" {
				t.Errorf("Expected status 'error', got '%s'", response.Status)
			}
			if response.Message != tt.message {
				t.Errorf("Expected message '%s', got '%s'", tt.message, response.Message)
			}
			if response.Data != nil {
				t.Error("Expected data to be nil in error response")
			}
		})
	}
}

func TestErrorWithSanitization(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectStatus int
	}{
		{
			name:         "public error",
			err:          NewPublicError("Public error message", 400),
			expectStatus: StatusCodeBadRequest,
		},
		{
			name:         "internal error",
			err:          NewInternalError("Internal error message", errors.New("internal"), 500),
			expectStatus: StatusCodeInternalServerError,
		},
		{
			name:         "not found error",
			err:          WrapInternalErrorWithStatus("Resource not found", errors.New("not found"), StatusCodeNotFound),
			expectStatus: StatusCodeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ErrorWithSanitization(w, tt.err)

			// Check status code
			if w.Code != tt.expectStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectStatus, w.Code)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", contentType)
			}

			// Check response body
			var response APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Status != "error" {
				t.Errorf("Expected status 'error', got '%s'", response.Status)
			}
		})
	}
}

func TestStatusCodeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant int
		expected int
	}{
		{"StatusCodeOK", StatusCodeOK, http.StatusOK},
		{"StatusCodeCreated", StatusCodeCreated, http.StatusCreated},
		{"StatusCodeAccepted", StatusCodeAccepted, http.StatusAccepted},
		{"StatusCodeBadRequest", StatusCodeBadRequest, http.StatusBadRequest},
		{"StatusCodeUnauthorized", StatusCodeUnauthorized, http.StatusUnauthorized},
		{"StatusCodeNotFound", StatusCodeNotFound, http.StatusNotFound},
		{"StatusCodeConflict", StatusCodeConflict, http.StatusConflict},
		{"StatusCodeInternalServerError", StatusCodeInternalServerError, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, tt.constant)
			}
		})
	}
}
