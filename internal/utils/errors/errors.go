package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError represents a standardized error response
type APIError struct {
	StatusCode  int    `json:"status_code"`
	Message     string `json:"message"`
	ErrorDetail string `json:"error,omitempty"`
}

// NewAPIError creates a new APIError instance
func NewAPIError(statusCode int, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("status %d: %s", e.StatusCode, e.Message)
}

// WriteError writes an error response in JSON format
func WriteError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(NewAPIError(statusCode, message))
}

// ValidationError represents input validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a collection of validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// Add adds a new validation error
func (v *ValidationErrors) Add(field, message string) {
	v.Errors = append(v.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// HasErrors checks if there are any validation errors
func (v *ValidationErrors) HasErrors() bool {
	return len(v.Errors) > 0
}

// WriteValidationError writes validation errors in JSON format
func WriteValidationError(w http.ResponseWriter, errors *ValidationErrors) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(struct {
		*APIError
		ValidationErrors []ValidationError `json:"validation_errors"`
	}{
		APIError:         NewAPIError(http.StatusBadRequest, "Validation failed"),
		ValidationErrors: errors.Errors,
	})
}
