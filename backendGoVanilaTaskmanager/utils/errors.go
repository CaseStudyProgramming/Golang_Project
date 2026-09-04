package utils

import "errors"

// Common public errors (safe to show to users)
var (
	// Validation errors
	ErrInvalidInput      = NewPublicError("Invalid input provided", 400)
	ErrMissingRequired   = NewPublicError("Required field is missing", 400)
	ErrInvalidFormat     = NewPublicError("Invalid format", 400)
	
	// Authentication errors
	ErrUnauthorized      = NewPublicError("Unauthorized access", 401)
	ErrInvalidCredentials = NewPublicError("Invalid email or password", 401)
	ErrTokenExpired      = NewPublicError("Token has expired", 401)
	
	// Authorization errors
	ErrForbidden         = NewPublicError("You do not have permission to perform this action", 403)
	
	// Resource errors
	ErrNotFound          = NewPublicError("Resource not found", 404)
	ErrConflict          = NewPublicError("Resource already exists", 409)
	
	// Business logic errors
	ErrInvalidOperation  = NewPublicError("Invalid operation", 400)
	ErrInvalidState      = NewPublicError("Invalid state for this operation", 400)
	
	// Generic internal error (message shown to users)
	ErrInternal          = NewPublicError("An internal error occurred. Please try again later", 500)
)

// WrapInternalError wraps an internal error with context
func WrapInternalError(message string, err error) *AppError {
	return NewInternalError(message, err, 500)
}

// WrapInternalErrorWithStatus wraps an internal error with custom status code
func WrapInternalErrorWithStatus(message string, err error, statusCode int) *AppError {
	return NewInternalError(message, err, statusCode)
}

// CreatePublicError creates a public error with custom message and status
func CreatePublicError(message string, statusCode int) *AppError {
	return NewPublicError(message, statusCode)
}

// IsNotFoundError checks if error is a not found error
func IsNotFoundError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Message == ErrNotFound.Message
	}
	return errors.Is(err, ErrNotFound)
}

// IsConflictError checks if error is a conflict error
func IsConflictError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Message == ErrConflict.Message
	}
	return errors.Is(err, ErrConflict)
}

// IsUnauthorizedError checks if error is an unauthorized error
func IsUnauthorizedError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Message == ErrUnauthorized.Message || appErr.Message == ErrInvalidCredentials.Message
	}
	return errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrInvalidCredentials)
}
