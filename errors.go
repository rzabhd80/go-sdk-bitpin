package bitpin

import "fmt"

// APIError represents an error returned by the Bitpin API.
//
// This struct contains the HTTP status code and the error message provided
// by the API. It implements the `error` interface, allowing it to be used
// seamlessly as an error in Go's standard library functions.
//
// Fields:
//   - StatusCode: The HTTP status code returned by the API.
//   - Message: The error message returned by the API.
//
// Example:
//
//	err := &APIError{
//	    StatusCode: 400,
//	    Message:    "Invalid request",
//	}
//	fmt.Println(err) // Output: API error: 400 Invalid request
type APIError struct {
	StatusCode int
	Message    string
}

// Error implements the error interface for APIError.
//
// It returns a formatted string containing the HTTP status code and the error
// message, making it suitable for logging or debugging.
//
// Returns:
//   - A string representation of the API error in the format:
//     "API error: <StatusCode> <Message>"
func (e *APIError) Error() string {
	return fmt.Sprintf("API error: %d %s", e.StatusCode, e.Message)
}
