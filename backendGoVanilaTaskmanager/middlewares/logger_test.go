package middlewares

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	// Capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Create request
	req := httptest.NewRequest("GET", "/test/path", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	// Create response recorder
	w := httptest.NewRecorder()

	// Apply logger middleware
	handler := Logger(testHandler)
	handler.ServeHTTP(w, req)

	// Check that handler was called
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check that log was written
	logOutput := logBuffer.String()
	if logOutput == "" {
		t.Error("Expected log output, got empty string")
	}

	// Check log format contains expected parts
	expectedParts := []string{"GET", "/test/path", "200", "127.0.0.1:12345"}
	for _, part := range expectedParts {
		if !strings.Contains(logOutput, part) {
			t.Errorf("Expected log output to contain '%s', got: %s", part, logOutput)
		}
	}
}

func TestLogger_DifferentMethods(t *testing.T) {
	// Capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		logBuffer.Reset()

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(method, "/test", nil)
		w := httptest.NewRecorder()

		handler := Logger(testHandler)
		handler.ServeHTTP(w, req)

		logOutput := logBuffer.String()
		if !strings.Contains(logOutput, method) {
			t.Errorf("Expected log output to contain '%s', got: %s", method, logOutput)
		}
	}
}

func TestLogger_DifferentStatusCodes(t *testing.T) {
	// Capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)

	statusCodes := []int{http.StatusOK, http.StatusCreated, http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError}

	for _, statusCode := range statusCodes {
		logBuffer.Reset()

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(statusCode)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		handler := Logger(testHandler)
		handler.ServeHTTP(w, req)

		logOutput := logBuffer.String()
		// Check that the status code appears in the log output
		statusStr := fmt.Sprintf("%d", statusCode)
		if !strings.Contains(logOutput, statusStr) {
			t.Errorf("Expected log output to contain status code %d, got: %s", statusCode, logOutput)
		}
	}
}

func TestLogger_ResponseWriter(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"default status", http.StatusOK},
		{"custom status", http.StatusCreated},
		{"error status", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &responseWriter{
				ResponseWriter: httptest.NewRecorder(),
				statusCode:     http.StatusOK,
			}

			rw.WriteHeader(tt.statusCode)

			if rw.statusCode != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, rw.statusCode)
			}
		})
	}
}

func TestLogger_WriteHeader_Default(t *testing.T) {
	w := httptest.NewRecorder()
	rw := &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	// If WriteHeader is not called, status should remain default
	if rw.statusCode != http.StatusOK {
		t.Errorf("Expected default status code %d, got %d", http.StatusOK, rw.statusCode)
	}
}

func TestLogger_ChainedHandlers(t *testing.T) {
	// Capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)

	// Create multiple handlers
	handler1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("handler1"))
	})

	handler2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("handler2"))
	})

	// Chain handlers with logger
	loggedHandler1 := Logger(handler1)
	loggedHandler2 := Logger(handler2)

	// Create request
	req := httptest.NewRequest("GET", "/chained", nil)
	w := httptest.NewRecorder()

	// Call first handler
	loggedHandler1.ServeHTTP(w, req)

	// Check that log was written
	logOutput := logBuffer.String()
	if logOutput == "" {
		t.Error("Expected log output from chained handler")
	}

	// Reset for second handler
	logBuffer.Reset()
	w = httptest.NewRecorder()

	// Call second handler
	loggedHandler2.ServeHTTP(w, req)

	// Check that log was written
	logOutput = logBuffer.String()
	if logOutput == "" {
		t.Error("Expected log output from second chained handler")
	}
}

func TestLogger_WithQueryParams(t *testing.T) {
	// Capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test?param1=value1&param2=value2", nil)
	w := httptest.NewRecorder()

	handler := Logger(testHandler)
	handler.ServeHTTP(w, req)

	logOutput := logBuffer.String()
	// Check that the base path is logged (query params might or might not be included)
	if !strings.Contains(logOutput, "/test") {
		t.Errorf("Expected log output to contain base path, got: %s", logOutput)
	}
}
