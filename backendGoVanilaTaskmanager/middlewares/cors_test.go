package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS(t *testing.T) {
	tests := []struct {
		name           string
		allowedOrigins []string
		origin         string
		method         string
		expectedStatus int
		expectedHeader string
	}{
		{
			name:           "Allowed origin gets CORS headers",
			allowedOrigins: []string{"http://localhost:5173", "http://localhost:3000"},
			origin:         "http://localhost:5173",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectedHeader: "http://localhost:5173",
		},
		{
			name:           "Preflight request for allowed origin",
			allowedOrigins: []string{"http://localhost:5173", "http://localhost:3000"},
			origin:         "http://localhost:3000",
			method:         "OPTIONS",
			expectedStatus: http.StatusNoContent,
			expectedHeader: "http://localhost:3000",
		},
		{
			name:           "Disallowed origin gets no CORS headers",
			allowedOrigins: []string{"http://localhost:5173"},
			origin:         "http://evil.com",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectedHeader: "",
		},
		{
			name:           "Preflight request for disallowed origin gets forbidden",
			allowedOrigins: []string{"http://localhost:5173"},
			origin:         "http://evil.com",
			method:         "OPTIONS",
			expectedStatus: http.StatusForbidden,
			expectedHeader: "",
		},
		{
			name:           "Wildcard allows all origins",
			allowedOrigins: []string{"*"},
			origin:         "http://any-origin.com",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectedHeader: "http://any-origin.com", // Returns specific origin even with wildcard for security
		},
		{
			name:           "No origin header uses first allowed origin",
			allowedOrigins: []string{"http://localhost:5173"},
			origin:         "",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectedHeader: "http://localhost:5173", // Should set first allowed origin as fallback for same-origin requests
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a simple handler that returns OK
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Apply CORS middleware
			handler := CORS(tt.allowedOrigins)(nextHandler)

			// Create test request
			req := httptest.NewRequest(tt.method, "http://example.com/test", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			// Create test recorder
			rr := httptest.NewRecorder()

			// Execute request
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Check CORS header
			header := rr.Header().Get("Access-Control-Allow-Origin")
			if header != tt.expectedHeader {
				t.Errorf("Expected Access-Control-Allow-Origin header %q, got %q", tt.expectedHeader, header)
			}
		})
	}
}

func TestCORSSecurityHeaders(t *testing.T) {
	allowedOrigins := []string{"http://localhost:5173"}
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := CORS(allowedOrigins)(nextHandler)

	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check that security headers are set
	if rr.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("Expected Access-Control-Allow-Methods header to be set")
	}
	if rr.Header().Get("Access-Control-Allow-Headers") == "" {
		t.Error("Expected Access-Control-Allow-Headers header to be set")
	}
	if rr.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Error("Expected Access-Control-Allow-Credentials to be true")
	}
}
