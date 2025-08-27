package httpclient

import (
	"fmt"
	"net/http"
)

// APIError represents a structured API error response
type APIError struct {
	StatusCode int
	Message    string
	Details    string
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("API Error %d: %s - %s", e.StatusCode, e.Message, e.Details)
	}
	return fmt.Sprintf("API Error %d: %s", e.StatusCode, e.Message)
}

// IsNotFound checks if the error is a 404 Not Found error
func IsNotFound(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsServerError checks if the error is a 5xx server error
func IsServerError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode >= 500 && apiErr.StatusCode < 600
	}
	return false
}

// IsClientError checks if the error is a 4xx client error
func IsClientError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode >= 400 && apiErr.StatusCode < 500
	}
	return false
}

// NewAPIError creates a new APIError instance
func NewAPIError(statusCode int, message, details string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Details:    details,
	}
}
