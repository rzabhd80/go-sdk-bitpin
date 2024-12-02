package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// JWT represents the structure of a JSON Web Token (JWT) used for authentication
// and authorization. It includes details about the token type, expiration,
// associated user, and other metadata.
type JWT struct {
	// TokenType specifies the type of the token, such as "Bearer".
	TokenType string `json:"token_type"`

	// Exp is the expiration time of the token, represented as a Unix timestamp.
	// This indicates when the token will no longer be valid.
	Exp int `json:"exp"`

	// Jti is the unique identifier for the token, often used for tracking or
	// invalidating specific tokens.
	Jti string `json:"jti"`

	// UserId is the unique identifier of the user associated with this token.
	// It links the token to a specific user in the system.
	UserId int `json:"user_id"`

	// Ip is a list of IP addresses associated with the token, which may be used
	// for security purposes, such as validating token usage from specific IPs.
	Ip []string `json:"ip"`

	// ApiCredentialId is the unique identifier for the API credential associated
	// with this token. It links the token to an API key or credential.
	ApiCredentialId int `json:"api_credential_id"`
}

// HumanReadable returns a string representation of the JWT in a human-readable format.
// It includes details such as the token type, expiration time, JTI, user ID, IP addresses,
// and API credential ID.
func (j JWT) HumanReadable() string {
	return fmt.Sprintf(
		"Token Type: %s\nExp: %d\nJti: %s\nUserId: %d\nIp: %v\nApi Credential Id: %d",
		j.TokenType, j.Exp, j.Jti, j.UserId, j.Ip, j.ApiCredentialId,
	)
}

// IsExpired checks whether the JWT has expired based on the current Unix timestamp.
// Returns true if the token's expiration time is earlier than the current time.
func (j JWT) IsExpired() bool {
	return j.Exp < int(time.Now().Unix())
}

// IsExpiredIn checks whether the JWT will expire within the specified duration from now.
// Takes a time.Duration as input and returns true if the token will expire in the given timeframe.
func (j JWT) IsExpiredIn(t time.Duration) bool {
	return j.Exp < int(time.Now().Add(t).Unix())
}

// DecodeJWT decodes a JWT string into a JWT struct.
// It parses the JWT token, extracts the claims, and maps them to the JWT struct.
//
// Parameters:
// - tokenString: The JWT string to decode.
//
// Returns:
// - A pointer to the JWT struct if decoding is successful.
// - An error if the token cannot be parsed or the claims cannot be mapped.
//
// Example:
//
//	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
//	jwt, err := DecodeJWT(tokenString)
//	if err != nil {
//		log.Fatalf("Error decoding JWT: %v", err)
//	}
func DecodeJWT(tokenString string) (*JWT, error) {
	// Parse the JWT token.
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract the claims.
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to parse claims as MapClaims")
	}

	// Marshal claims back to JSON.
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal claims: %w", err)
	}

	// Unmarshal JSON into the JWT struct.
	var payload JWT
	if err := json.Unmarshal(claimsJSON, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims into struct: %w", err)
	}

	return &payload, nil
}
