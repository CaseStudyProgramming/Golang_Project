package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestNewSwaggerController(t *testing.T) {
	controller := NewSwaggerController()
	if controller == nil {
		t.Fatal("NewSwaggerController returned nil")
	}
}

func TestSwaggerController_ServeSwaggerUI(t *testing.T) {
	// Create a temporary swagger directory with index.html
	tempDir := t.TempDir()
	swaggerDir := filepath.Join(tempDir, "swagger")
	if err := os.Mkdir(swaggerDir, 0755); err != nil {
		t.Fatalf("Failed to create swagger directory: %v", err)
	}

	indexContent := []byte("<!DOCTYPE html><html><body>Swagger UI</body></html>")
	if err := os.WriteFile(filepath.Join(swaggerDir, "index.html"), indexContent, 0644); err != nil {
		t.Fatalf("Failed to create index.html: %v", err)
	}

	// Change to temp directory
	originalDir, _ := os.Getwd()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(originalDir)

	controller := NewSwaggerController()
	req := httptest.NewRequest("GET", "/swagger", nil)
	w := httptest.NewRecorder()

	controller.ServeSwaggerUI(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "text/html" {
		t.Errorf("Expected Content-Type text/html, got %s", w.Header().Get("Content-Type"))
	}

	body := w.Body.String()
	if body != string(indexContent) {
		t.Errorf("Expected body '%s', got '%s'", string(indexContent), body)
	}
}

func TestSwaggerController_ServeSwaggerUI_NotFound(t *testing.T) {
	// Change to a directory without swagger folder
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(originalDir)

	controller := NewSwaggerController()
	req := httptest.NewRequest("GET", "/swagger", nil)
	w := httptest.NewRecorder()

	controller.ServeSwaggerUI(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestSwaggerController_ServeOpenAPISpec(t *testing.T) {
	// Create a temporary swagger directory with openapi.yaml
	tempDir := t.TempDir()
	swaggerDir := filepath.Join(tempDir, "swagger")
	if err := os.Mkdir(swaggerDir, 0755); err != nil {
		t.Fatalf("Failed to create swagger directory: %v", err)
	}

	openapiContent := []byte("openapi: 3.0.0\ninfo:\n  title: Test API")
	if err := os.WriteFile(filepath.Join(swaggerDir, "openapi.yaml"), openapiContent, 0644); err != nil {
		t.Fatalf("Failed to create openapi.yaml: %v", err)
	}

	// Change to temp directory
	originalDir, _ := os.Getwd()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(originalDir)

	controller := NewSwaggerController()
	req := httptest.NewRequest("GET", "/swagger/openapi.yaml", nil)
	w := httptest.NewRecorder()

	controller.ServeOpenAPISpec(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/x-yaml" {
		t.Errorf("Expected Content-Type application/x-yaml, got %s", w.Header().Get("Content-Type"))
	}

	body := w.Body.String()
	if body != string(openapiContent) {
		t.Errorf("Expected body '%s', got '%s'", string(openapiContent), body)
	}
}

func TestSwaggerController_ServeOpenAPISpec_NotFound(t *testing.T) {
	// Change to a directory without swagger folder
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(originalDir)

	controller := NewSwaggerController()
	req := httptest.NewRequest("GET", "/swagger/openapi.yaml", nil)
	w := httptest.NewRecorder()

	controller.ServeOpenAPISpec(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}
