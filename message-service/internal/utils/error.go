package utils

import "fmt"

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func NewAPIError(code int, message string, details string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.Details)
}
