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

const (
	StatusCodeOK                  = http.StatusOK
	StatusCodeCreated             = http.StatusCreated
	StatusCodeAccepted            = http.StatusAccepted
	StatusCodeBadRequest          = http.StatusBadRequest
	StatusCodeUnauthorized        = http.StatusUnauthorized
	StatusCodeNotFound            = http.StatusNotFound
	StatusCodeInternalServerError = http.StatusInternalServerError
)

func SuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Message: message, Data: data})
}

func ErrorResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: message, Data: nil})
}

func NotFoundResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: message, Data: nil})
}

func InternalErrorResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: message, Data: nil})
}
