package middlewares

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestRecovery(t *testing.T) {
	// Capture log output
	var logBuffer strings.Builder
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)

	// Create a handler that panics
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	// Create request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Apply recovery middleware
	handler := Recovery(panicHandler)
	handler.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Check content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Check response body
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "error" {
		t.Errorf("Expected status 'error', got '%v'", response["status"])
	}

	if response["message"] != "Internal server error" {
		t.Errorf("Expected message 'Internal server error', got '%v'", response["message"])
	}

	// Check that panic was logged
	logOutput := logBuffer.String()
	if !strings.Contains(logOutput, "PANIC recovered") {
		t.Error("Expected panic to be logged")
	}
	if !strings.Contains(logOutput, "test panic") {
		t.Error("Expected panic message to be logged")
	}
}

func TestRecovery_NoPanic(t *testing.T) {
	// Capture log output
	var logBuffer strings.Builder
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)

	// Create a normal handler
	normalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Create request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Apply recovery middleware
	handler := Recovery(normalHandler)
	handler.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	if w.Body.String() != "success" {
		t.Errorf("Expected response 'success', got '%s'", w.Body.String())
	}

	// Check that no panic was logged
	logOutput := logBuffer.String()
	if strings.Contains(logOutput, "PANIC recovered") {
		t.Error("Expected no panic to be logged for normal handler")
	}
}

func TestRecovery_DifferentPanicTypes(t *testing.T) {
	// Capture log output
	var logBuffer strings.Builder
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)

	panicTypes := []interface{}{
		"string panic",
		123,
		errors.New("error panic"),
		struct{ name string }{"struct panic"},
	}

	for _, panicValue := range panicTypes {
		logBuffer.Reset()

		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(panicValue)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		handler := Recovery(panicHandler)
		handler.ServeHTTP(w, req)

		// Check that recovery middleware handled the panic
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d for panic type %T, got %d", http.StatusInternalServerError, panicValue, w.Code)
		}

		// Check that panic was logged
		logOutput := logBuffer.String()
		if !strings.Contains(logOutput, "PANIC recovered") {
			t.Error("Expected panic to be logged")
		}
	}
}

func TestRecovery_ChainedHandlers(t *testing.T) {
	// Capture log output
	var logBuffer strings.Builder
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)

	// Create handlers - first one panics, second one normal
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("chain panic")
	})

	normalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Apply recovery middleware to both
	recoveredPanicHandler := Recovery(panicHandler)
	recoveredNormalHandler := Recovery(normalHandler)

	// Test panic handler
	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	recoveredPanicHandler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d for panic handler, got %d", http.StatusInternalServerError, w.Code)
	}

	// Test normal handler
	logBuffer.Reset()
	req = httptest.NewRequest("GET", "/normal", nil)
	w = httptest.NewRecorder()
	recoveredNormalHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for normal handler, got %d", http.StatusOK, w.Code)
	}

	// Check that no panic was logged for normal handler
	logOutput := logBuffer.String()
	if strings.Contains(logOutput, "PANIC recovered") {
		t.Error("Expected no panic to be logged for normal handler in chain")
	}
}

func TestRecovery_ResponseFormat(t *testing.T) {
	// Capture log output
	var logBuffer strings.Builder
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)

	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler := Recovery(panicHandler)
	handler.ServeHTTP(w, req)

	// Check response format
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify all expected fields are present
	expectedFields := []string{"status", "message", "data"}
	for _, field := range expectedFields {
		if _, exists := response[field]; !exists {
			t.Errorf("Expected response to contain field '%s'", field)
		}
	}

	// Verify field values
	if response["status"] != "error" {
		t.Errorf("Expected status 'error', got '%v'", response["status"])
	}

	if response["data"] != nil {
		t.Error("Expected data to be nil in error response")
	}
}

func TestRecovery_MultipleRequests(t *testing.T) {
	// Capture log output
	var logBuffer strings.Builder
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)

	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("repeated panic")
	})

	handler := Recovery(panicHandler)

	// Make multiple requests that all panic
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Request %d: Expected status %d, got %d", i+1, http.StatusInternalServerError, w.Code)
		}
	}

	// Check that all panics were logged
	logOutput := logBuffer.String()
	panicCount := strings.Count(logOutput, "PANIC recovered")
	if panicCount != 3 {
		t.Errorf("Expected 3 panic logs, got %d", panicCount)
	}
}
