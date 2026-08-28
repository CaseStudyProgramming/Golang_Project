package middlewares

import (
	"encoding/json"
	"log"
	"net/http"
)

// Recovery middleware to recover from panics and prevent server crashes
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				log.Printf("PANIC recovered: %v", err)

				// Set response headers
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				// Send error response
				response := map[string]interface{}{
					"status":  "error",
					"message": "Internal server error",
					"data":    nil,
				}
				json.NewEncoder(w).Encode(response)
			}
		}()

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
