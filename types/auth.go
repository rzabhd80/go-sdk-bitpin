package types

// AuthenticationParams represents the parameters required for user authentication.
// This is typically used to authenticate with an API by providing credentials.
type AuthenticationParams struct {
	// ApiKey is the user's API key, used as an identifier for authentication.
	ApiKey string `json:"api_key"`

	// SecretKey is the user's secret key, used alongside the API key to authenticate
	// and authorize API requests.
	SecretKey string `json:"secret_key"`
}

// AuthenticationResponse represents the response received after a successful
// authentication request. It includes tokens used for accessing the API.
type AuthenticationResponse struct {
	// Refresh is the token used to obtain a new access token when the current
	// access token expires.
	Refresh string `json:"refresh"`

	// Access is the token used to authenticate API requests. It has a limited
	// lifespan and must be refreshed periodically.
	Access string `json:"access"`
}

// RefreshTokenParams represents the parameters required to refresh an access token.
// This struct is used to request a new access token using a valid refresh token.
type RefreshTokenParams struct {
	// Refresh is the refresh token obtained during authentication, used to
	// request a new access token.
	Refresh string `json:"refresh"`
}

// RefreshTokenResponse represents the response received after requesting a new
// access token using a refresh token.
type RefreshTokenResponse struct {
	// Access is the new access token generated after a refresh request. It is
	// used to authenticate API requests and replaces the expired token.
	Access string `json:"access"`
}
