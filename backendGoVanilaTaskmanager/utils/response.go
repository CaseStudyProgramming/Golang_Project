package utils

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// TimezoneAwareResponse adds timezone information to the response
type TimezoneAwareResponse struct {
	APIResponse
	Timezone string `json:"timezone,omitempty"`
}

// universal status code
const (
	StatusCodeOK                  = http.StatusOK
	StatusCodeCreated             = http.StatusCreated
	StatusCodeAccepted            = http.StatusAccepted
	StatusCodeBadRequest          = http.StatusBadRequest
	StatusCodeUnauthorized        = http.StatusUnauthorized
	StatusCodeNotFound            = http.StatusNotFound
	StatusCodeConflict            = http.StatusConflict
	StatusCodeInternalServerError = http.StatusInternalServerError
)

// success response for universal status code
func SuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Message: message, Data: data})
}

// success response with timezone support
func SuccessResponseWithTimezone(w http.ResponseWriter, statusCode int, message string, data interface{}, timezone string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := TimezoneAwareResponse{
		APIResponse: APIResponse{Status: "success", Message: message, Data: data},
		Timezone:    timezone,
	}
	json.NewEncoder(w).Encode(response)
}

// error response for universal status code
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: message, Data: nil})
}

// ErrorWithSanitization handles errors with automatic sanitization and logging
func ErrorWithSanitization(w http.ResponseWriter, err error) {
	statusCode, userMessage := HandleAPIError(err)
	ErrorResponse(w, statusCode, userMessage)
}
