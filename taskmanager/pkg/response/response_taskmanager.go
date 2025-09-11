package response_test

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// universal status code
const (
	StatusCodeOK                  = http.StatusOK
	StatusCodeCreated             = http.StatusCreated
	StatusCodeAccepted            = http.StatusAccepted
	StatusCodeBadRequest          = http.StatusBadRequest
	StatusCodeUnauthorized        = http.StatusUnauthorized
	StatusCodeNotFound            = http.StatusNotFound
	StatusCodeInternalServerError = http.StatusInternalServerError
)

// success response for universal status code
func SuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Message: message, Data: data})
}

// error response for universal status code
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: message, Data: nil})
}
