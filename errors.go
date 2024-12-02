package bitpin

import (
	"encoding/json"
	"fmt"
)

// GoBitpinError is the base error type for all errors in the SDK
type GoBitpinError struct {
	Message string
	Err     error
}

func (e *GoBitpinError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap implements the errors.Unwrap interface
func (e *GoBitpinError) Unwrap() error {
	return e.Err
}

// RequestError represents errors that occur during HTTP request creation or sending
type RequestError struct {
	GoBitpinError
	Operation string // e.g., "creating request", "sending request"
}

// APIError represents errors returned by the Bitpin API
type APIError struct {
	GoBitpinError
	StatusCode int
	Details    map[string][]string // Store field-specific errors
}

// parseErrorResponse attempts to parse various error response formats from the API
func parseErrorResponse(statusCode int, respBody []byte) *APIError {
	var details map[string][]string

	// Try parsing as map[string][]string first (for field-specific errors)
	if err := json.Unmarshal(respBody, &details); err == nil {
		return &APIError{
			GoBitpinError: GoBitpinError{
				Message: fmt.Sprintf("API error (status %d): %s", statusCode, formatErrorDetails(details)),
			},
			StatusCode: statusCode,
			Details:    details,
		}
	}

	// Try parsing as map[string]string (for simple key-value errors)
	var simpleDetails map[string]string
	if err := json.Unmarshal(respBody, &simpleDetails); err == nil {
		details = make(map[string][]string)
		for k, v := range simpleDetails {
			details[k] = []string{v}
		}
		return &APIError{
			GoBitpinError: GoBitpinError{
				Message: fmt.Sprintf("API error (status %d): %s", statusCode, formatErrorDetails(details)),
			},
			StatusCode: statusCode,
			Details:    details,
		}
	}

	// Try parsing as ErrorResponse struct
	var errResp struct {
		Detail   string            `json:"detail"`
		Code     string            `json:"code"`
		Messages map[string]string `json:"messages"`
	}
	if err := json.Unmarshal(respBody, &errResp); err == nil {
		details = make(map[string][]string)
		if errResp.Detail != "" {
			details["detail"] = []string{errResp.Detail}
		}
		if errResp.Code != "" {
			details["code"] = []string{errResp.Code}
		}
		for k, v := range errResp.Messages {
			details[k] = []string{v}
		}
		return &APIError{
			GoBitpinError: GoBitpinError{
				Message: fmt.Sprintf("API error (status %d): %s", statusCode, formatErrorDetails(details)),
			},
			StatusCode: statusCode,
			Details:    details,
		}
	}

	// If all parsing attempts fail, return error with raw response
	return &APIError{
		GoBitpinError: GoBitpinError{
			Message: fmt.Sprintf("API error (status %d): %s", statusCode, string(respBody)),
		},
		StatusCode: statusCode,
		Details:    map[string][]string{"raw": {string(respBody)}},
	}
}

// formatErrorDetails creates a human-readable error message from the error details
func formatErrorDetails(details map[string][]string) string {
	msg := ""
	for field, errors := range details {
		if len(errors) > 0 {
			if msg != "" {
				msg += "; "
			}
			msg += fmt.Sprintf("%s: %v", field, errors[0])
		}
	}
	return msg
}
