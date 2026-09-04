package utils

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

// ErrorType represents the category of error
type ErrorType int

const (
	// ErrorTypePublic are safe to show to users
	ErrorTypePublic ErrorType = iota
	// ErrorTypeInternal should not be shown to users (logged internally only)
	ErrorTypeInternal
)

// AppError represents an application error with type handling
type AppError struct {
	Message    string
	Err        error
	ErrorType  ErrorType
	StatusCode int
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewPublicError creates a user-facing error (safe to display)
func NewPublicError(message string, statusCode int) *AppError {
	return &AppError{
		Message:    message,
		ErrorType:  ErrorTypePublic,
		StatusCode: statusCode,
	}
}

// NewInternalError creates an internal error (logged, not shown to users)
func NewInternalError(message string, err error, statusCode int) *AppError {
	return &AppError{
		Message:    message,
		Err:        err,
		ErrorType:  ErrorTypeInternal,
		StatusCode: statusCode,
	}
}

// SanitizeError removes sensitive information from error messages
func SanitizeError(err error) string {
	if err == nil {
		return ""
	}

	errorMsg := err.Error()

	// Patterns to match sensitive information
	sensitivePatterns := []struct {
		pattern     string
		replacement string
	}{
		// Database connection strings and credentials
		{`(?i)password\s*=\s*[^\s,)]+`, "password=***"},
		{`(?i)user\s*=\s*[^\s,)]+`, "user=***"},
		{`(?i)host\s*=\s*[^\s,)]+`, "host=***"},
		{`(?i)dbname\s*=\s*[^\s,)]+`, "dbname=***"},
		{`(?i)database\s*=\s*[^\s,)]+`, "database=***"},

		// File paths (Windows and Unix)
		{`[A-Za-z]:\\Users\\[^\\]+\\`, "C:\\Users\\***\\"},
		{`/home/[^/]+/`, "/home/***/"},

		// SQL error details
		{`(?i)SQLSTATE\s+\w+`, "SQLSTATE=***"},
		{`(?i)near\s+["'][^"']+["']`, "near '***'"},
		{`(?i)syntax\s+error\s+near`, "syntax error near"},

		// IP addresses
		{`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`, "***.***.***.***"},

		// Email addresses
		{`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`, "***@***.***"},

		// API keys and tokens (common patterns)
		{`(?i)api[_-]?key\s*[:=]\s*[^\s,)]+`, "api_key=***"},
		{`(?i)token\s*[:=]\s*[^\s,)]+`, "token=***"},
		{`(?i)secret\s*[:=]\s*[^\s,)]+`, "secret=***"},
		{`[A-Za-z0-9]{32,}`, "***"},
	}

	// Apply all sanitization patterns
	for _, p := range sensitivePatterns {
		re := regexp.MustCompile(p.pattern)
		errorMsg = re.ReplaceAllString(errorMsg, p.replacement)
	}

	// Additional cleanup: remove stack traces if present
	if strings.Contains(errorMsg, "goroutine") || strings.Contains(errorMsg, "created by") {
		lines := strings.Split(errorMsg, "\n")
		var sanitizedLines []string
		for _, line := range lines {
			if !strings.Contains(line, "goroutine") &&
				!strings.Contains(line, "created by") &&
				!strings.Contains(line, "panic:") &&
				!strings.Contains(line, "runtime/") {
				sanitizedLines = append(sanitizedLines, line)
			}
		}
		errorMsg = strings.Join(sanitizedLines, "\n")
	}

	return errorMsg
}

// GetUserMessage returns a safe message for the user
func GetUserMessage(err error) string {
	if err == nil {
		return ""
	}

	// Check if it's an AppError
	if appErr, ok := err.(*AppError); ok {
		if appErr.ErrorType == ErrorTypePublic {
			return appErr.Message
		}
		// For internal errors, return a generic message
		return "An internal error occurred. Please try again later."
	}

	// For regular errors, sanitize and return if it looks safe
	sanitized := SanitizeError(err)

	// If the sanitized error still contains technical details, return generic message
	technicalKeywords := []string{"sql", "database", "connection", "timeout", "socket", "driver", "syntax"}
	for _, keyword := range technicalKeywords {
		if strings.Contains(strings.ToLower(sanitized), keyword) {
			return "An internal error occurred. Please try again later."
		}
	}

	return sanitized
}

// LogInternalError logs internal errors while keeping user messages generic
func LogInternalError(err error, context string) {
	if err == nil {
		return
	}

	// Log the full error for debugging
	log.Printf("[INTERNAL ERROR] %s: %v", context, err)
}

// HandleAPIError is a helper function to handle errors in API controllers
// It returns the appropriate status code and user-safe message
func HandleAPIError(err error) (int, string) {
	if err == nil {
		return 200, "Success"
	}

	// Check if it's an AppError
	if appErr, ok := err.(*AppError); ok {
		// Log internal errors
		if appErr.ErrorType == ErrorTypeInternal {
			LogInternalError(appErr.Err, appErr.Message)
		}
		return appErr.StatusCode, GetUserMessage(appErr)
	}

	// For regular errors, treat as internal
	LogInternalError(err, "API Error")
	return 500, "An internal error occurred. Please try again later."
}
