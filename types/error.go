package types

// ErrorResponse represents a standardized structure for API error responses.
// It includes details about the error, an error code, and additional messages
// providing context or specific field-related errors.
type ErrorResponse struct {
	// Detail provides a high-level description of the error, such as
	// "Invalid request parameters" or "Authentication failed."
	Detail string `json:"detail"`

	// Code is a unique identifier for the error type, which can be used to
	// programmatically handle specific error cases. For example, "400_BAD_REQUEST"
	// or "AUTH_ERROR".
	Code string `json:"code"`

	// Messages is a map of field-specific error messages, where the key is
	// the field name and the value is a message describing the issue.
	// For example, {"username": "This field is required", "password": "Too short"}.
	Messages map[string]string `json:"messages"`
}
