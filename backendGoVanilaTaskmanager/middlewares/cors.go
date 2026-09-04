package middlewares

import (
	"net/http"
	"os"
)

// CORS middleware to handle Cross-Origin Resource Sharing
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			env := os.Getenv("APP_ENV")
			if env == "" {
				env = "development"
			}

			// Check if the origin is in the allowed list
			allowed := false
			for _, allowedOrigin := range allowedOrigins {
				// In production, reject wildcard "*" origins for security
				if env == "production" && allowedOrigin == "*" {
					continue
				}
				if allowedOrigin == "*" || allowedOrigin == origin {
					allowed = true
					break
				}
			}

			// If no origin header is present, allow by default (e.g., same-origin requests)
			if origin == "" && len(allowedOrigins) > 0 {
				allowed = true
			}

			if allowed {
				// Set CORS headers with specific origin
				if origin != "" {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				} else {
					// Fallback to first allowed origin if no Origin header
					if len(allowedOrigins) > 0 {
						w.Header().Set("Access-Control-Allow-Origin", allowedOrigins[0])
					} else {
						// Only use wildcard in development, not in production
						if env == "production" {
							w.Header().Set("Access-Control-Allow-Origin", "")
						} else {
							w.Header().Set("Access-Control-Allow-Origin", "*")
						}
					}
				}
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
				w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				if allowed {
					w.WriteHeader(http.StatusNoContent)
				} else {
					w.WriteHeader(http.StatusForbidden)
				}
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
