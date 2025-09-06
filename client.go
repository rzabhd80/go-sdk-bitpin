package bitpin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	t "github.com/rzabhd80/go-sdk-bitpin/types"
	u "github.com/rzabhd80/go-sdk-bitpin/utils"
)

// Constants defining the API base URL and version.
const (
	// BaseUrl is the root URL for the Bitpin Market API.
	BaseUrl = "https://api.bitpin.ir"

	// Version specifies the API version.
	Version = "v1"
)

// ClientOptions represents the configuration options for creating a new API client.
// These options allow customization of the HTTP client, authentication tokens,
// API credentials, and automatic authentication/refresh behaviors.
type ClientOptions struct {
	// HttpClient is the custom HTTP client to be used for API requests.
	// If nil, the default HTTP client is used.
	HttpClient *http.Client

	// Timeout specifies the request timeout duration for the HTTP client.
	Timeout time.Duration

	// BaseUrl is the base URL of the API. Defaults to the constant BaseUrl
	// if not provided.
	BaseUrl string

	// AccessToken is the token used for authenticated API requests.
	AccessToken string

	// RefreshToken is the token used to obtain a new AccessToken when it expires.
	RefreshToken string

	// ApiKey is the API key for authentication.
	ApiKey string

	// SecretKey is the secret key for authentication.
	SecretKey string

	// AutoAuth enables automatic authentication if no valid tokens are provided.
	AutoAuth bool

	// AutoRefresh enables automatic refreshing of the access token when it expires.
	AutoRefresh bool
}

// Client represents the API client for interacting with the Bitpin Market API.
// It manages authentication, base URL, and API requests.
type Client struct {
	// HttpClient is the HTTP client used for API requests.
	// Defaults to the Go standard library's http.DefaultClient.
	HttpClient *http.Client

	// BaseUrl is the base URL of the API used by this client.
	// Defaults to the constant BaseUrl.
	BaseUrl string

	// AccessToken is the token used for authenticated API requests.
	AccessToken string

	// RefreshToken is the token used to obtain a new AccessToken when it expires.
	RefreshToken string

	// ApiKey is the API key for authentication.
	ApiKey string

	// SecretKey is the secret key for authentication.
	SecretKey string

	// AutoRefresh enables automatic refreshing of the access token when it expires.
	AutoRefresh bool
}

// NewClient initializes a new API client with the provided options.
// It configures the client for making API requests, including setting up
// authentication and handling automatic token refresh if enabled.
//
// Parameters:
//   - opts: A ClientOptions struct containing configuration for the client.
//     It includes settings for HTTP client, authentication tokens, API
//     credentials, base URL, and timeout.
//
// Returns:
//   - A pointer to a Client struct initialized with the specified options.
//   - An error if there are issues during setup, such as authentication or
//     auto-refresh errors.
//
// Behavior:
//   - If `opts.BaseUrl` is provided, it overrides the default `BaseUrl`.
//   - If `opts.HttpClient` is not provided, a default HTTP client with the
//     specified timeout is created.
//   - AccessToken and RefreshToken are set from the options.
//   - If `AutoRefresh` is enabled, the client attempts to refresh tokens on initialization.
//   - If both `ApiKey` and `SecretKey` are provided, the client attempts to authenticate.
//
// Example:
//
//	opts := ClientOptions{
//	    ApiKey:       "your-api-key",
//	    SecretKey:    "your-secret-key",
//	    Timeout:      10 * time.Second,
//	    AutoRefresh:  true,
//	}
//	client, err := NewClient(opts)
//	if err != nil {
//	    log.Fatalf("Failed to create client: %v", err)
//	}
func NewClient(opts ClientOptions) (*Client, error) {
	client := &Client{
		AutoRefresh: opts.AutoRefresh,
		BaseUrl:     BaseUrl,
	}

	if opts.BaseUrl != "" {
		client.BaseUrl = opts.BaseUrl
	}

	if opts.HttpClient != nil {
		client.HttpClient = opts.HttpClient
	} else {
		client.HttpClient = &http.Client{
			Timeout: opts.Timeout,
		}
	}

	client.AccessToken = opts.AccessToken
	client.RefreshToken = opts.RefreshToken
	client.ApiKey = opts.ApiKey
	client.SecretKey = opts.SecretKey

	if err := client.handleAutoRefresh(); err != nil {
		return nil, err
	}

	if opts.ApiKey != "" && opts.SecretKey != "" {
		if _, err := client.Authenticate(opts.ApiKey, opts.SecretKey); err != nil {
			return nil, err
		}
	}

	return client, nil
}

// assertAuth checks the authentication state of the given client by verifying
// the presence of both the access token and the refresh token.
//
// Parameters:
//   - client: A pointer to the Client struct whose authentication state is being checked.
//
// Returns:
//   - An error if either the access token or the refresh token is missing.
//   - nil if both tokens are present, indicating that the client is authenticated.
//
// Example:
//
//	err := assertAuth(client)
//	if err != nil {
//	    log.Fatalf("Authentication error: %v", err)
//	}
//
// Behavior:
//   - If `client.AccessToken` is empty, returns an error: "access token is empty".
//   - If `client.RefreshToken` is empty, returns an error: "refresh token is empty".
//   - Otherwise, returns nil to indicate the client is authenticated.
func assertAuth(client *Client) error {
	if client.AccessToken == "" {
		return &GoBitpinError{
			Message: "access token is empty",
			Err:     nil,
		}
	}
	if client.RefreshToken == "" {
		return &GoBitpinError{
			Message: "refresh token is empty",
			Err:     nil,
		}
	}
	return nil
}

// createApiURI constructs a full API URI for a given endpoint and API version.
// It combines the base URL, API version, and endpoint into a properly formatted URI.
//
// Parameters:
//   - endpoint: A string representing the specific API endpoint, such as "/orders".
//   - version: A string representing the API version. If empty, the default version
//     defined by the `Version` constant is used.
//
// Returns:
//   - A string containing the complete API URI.
//
// Behavior:
//   - If `version` is not provided, the default `Version` constant is used.
//   - Combines the `BaseUrl`, API version, and endpoint into the format:
//     "<BaseUrl>/api/<version><endpoint>".
//
// Example:
//
//	uri := client.createApiURI("/orders", "v1")
//	// Result: "https://api.bitpin.market/api/v1/orders"
//
//	uri := client.createApiURI("/trades", "")
//	// Result: "https://api.bitpin.market/api/v1/trades" (uses default version)
func (c *Client) createApiURI(endpoint string, version string) string {
	if version == "" {
		version = Version
	}
	return fmt.Sprintf("%s/api/%s%s", c.BaseUrl, version, endpoint)
}

// handleAutoRefresh ensures the client's tokens are valid and refreshes them if necessary.
// It checks both the access token and the refresh token, and performs re-authentication
// if required.
//
// Returns:
//   - An error if there are issues decoding the tokens, refreshing the access token,
//     or re-authenticating. Returns nil if all tokens are valid or successfully refreshed.
//
// Behavior:
//   - If the access token is provided, it is decoded and checked for expiration.
//   - If expired, the `RefreshAccessToken` method is called to refresh it.
//   - If the refresh token is provided, it is decoded and checked for expiration.
//   - If expired, and API credentials (`ApiKey` and `SecretKey`) are available,
//     the client re-authenticates using `Authenticate`.
//   - Returns an error if the refresh token is expired but API credentials are missing.
//
// Example:
//
//	err := client.handleAutoRefresh()
//	if err != nil {
//	    log.Fatalf("Auto-refresh failed: %v", err)
//	}
//
// Dependencies:
//   - Requires the `DecodeJWT` method to decode JWT tokens and check their expiration.
//   - Calls `RefreshAccessToken` to refresh the access token if necessary.
//   - Calls `Authenticate` with API credentials if re-authentication is needed.
//
// Errors:
//   - "error decoding access token: %v" if the access token cannot be decoded.
//   - "error refreshing access token: %v" if the access token cannot be refreshed.
//   - "error decoding refresh token: %v" if the refresh token cannot be decoded.
//   - "API key and/or secret key are empty" if re-authentication is required but credentials are missing.
//   - "error re-authenticating: %v" if re-authentication fails.
func (c *Client) handleAutoRefresh() error {
	if c.AccessToken != "" {
		decoded, err := u.DecodeJWT(c.AccessToken)
		if err != nil {
			return err
		}
		if decoded.IsExpired() {
			err = c.RefreshAccessToken()
			if err != nil {
				return err
			}
		}
	}

	if c.RefreshToken != "" {
		decoded, err := u.DecodeJWT(c.RefreshToken)
		if err != nil {
			return err
		}

		if decoded.IsExpired() {
			if c.ApiKey == "" || c.SecretKey == "" {
				return &GoBitpinError{
					Message: "API key and/or secret key are empty",
					Err:     nil,
				}
			}

			_, err = c.Authenticate(c.ApiKey, c.SecretKey)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Request sends an HTTP request to the specified URL and handles the response.
// It supports both GET and POST methods, optional authentication, and automatic
// token refresh. The request body can be serialized from a struct, and the response
// can be unmarshaled into a given result object.
//
// Parameters:
//   - method: The HTTP method for the request, such as "GET" or "POST".
//   - url: The full URL for the API endpoint.
//   - auth: A boolean indicating whether the request requires authentication.
//     If true, the method adds an Authorization header with the access token
//     and handles automatic token refresh if enabled.
//   - body: An optional request body. For GET requests, it is converted into URL
//     parameters; for POST requests, it is marshaled to JSON.
//   - result: A pointer to a variable where the response body should be unmarshaled.
//     If nil, the response body is not unmarshaled.
//
// Returns:
//   - An error if there are issues during request creation, sending, or response
//     processing. If the response status code indicates an error (non-2xx), an
//     `APIError` is returned.
//
// Example:
//
//	var result MyResponseStruct
//	err := client.Request("POST", "https://api.example.com/resource", true, requestBody, &result)
//	if err != nil {
//	    log.Fatalf("API request failed: %v", err)
//	}
//
// Behavior:
//   - For GET requests, the body is converted into URL parameters using `StructToURLParams`.
//   - For POST requests, the body is marshaled to JSON.
//   - Adds the `Authorization` header if `auth` is true and the client has valid tokens.
//   - Refreshes tokens automatically if `AutoRefresh` is enabled and tokens are expired.
//   - Handles non-2xx HTTP responses by returning an `APIError` containing the status
//     code and error message.
//   - Unmarshals the response body into the `result` parameter if provided.
//
// Errors:
//   - "error converting struct to URL params: %v" for GET body conversion errors.
//   - "error marshaling Request body: %v" for POST body marshaling errors.
//   - "error creating Request: %v" for request creation failures.
//   - "error sending Request: %v" for HTTP client errors.
//   - "error reading response body: %v" for response body read errors.
//   - "error unmarshaling response: %v" for JSON unmarshal errors.
//
// Dependencies:
//   - `StructToURLParams` for struct-to-URL parameter conversion.
//   - `handleAutoRefresh` for automatic token refresh.
//   - `assertAuth` for ensuring authentication tokens are valid.
//   - `APIError` for structured error responses.
//
// Request sends an HTTP request to the specified URL and handles the response
func (c *Client) Request(method string, url string, auth bool, body interface{}, result interface{}) error {
	var reqBody []byte
	var err error

	if method == "GET" {
		if body != nil {
			urlParams, err := u.StructToURLParams(body)
			if err != nil {
				return &RequestError{
					GoBitpinError: GoBitpinError{
						Message: "failed to convert struct to URL params",
						Err:     err,
					},
					Operation: "preparing request parameters",
				}
			}
			url += "?" + urlParams
		}
	}

	if method == "POST" {
		if body != nil {
			reqBody, err = json.Marshal(body)
			if err != nil {
				return &RequestError{
					GoBitpinError: GoBitpinError{
						Message: "failed to marshal request body",
						Err:     err,
					},
					Operation: "preparing request body",
				}
			}
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return &RequestError{
			GoBitpinError: GoBitpinError{
				Message: "failed to create request",
				Err:     err,
			},
			Operation: "creating request",
		}
	}

	req.Header.Set("Content-Type", "application/json")

	if auth {
		if c.AutoRefresh {
			if err := c.handleAutoRefresh(); err != nil {
				return &GoBitpinError{
					Message: "failed to refresh authentication",
					Err:     err,
				}
			}
		}

		if err := assertAuth(c); err != nil {
			return &GoBitpinError{
				Message: "authentication validation failed",
				Err:     err,
			}
		}

		req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return &RequestError{
			GoBitpinError: GoBitpinError{
				Message: "failed to send request",
				Err:     err,
			},
			Operation: "sending request",
		}
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RequestError{
			GoBitpinError: GoBitpinError{
				Message: "failed to read response body",
				Err:     err,
			},
			Operation: "reading response",
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseErrorResponse(resp.StatusCode, respBody)
	}

	if result != nil {
		if err = json.Unmarshal(respBody, result); err != nil {
			return &RequestError{
				GoBitpinError: GoBitpinError{
					Message: "failed to unmarshal response",
					Err:     err,
				},
				Operation: "parsing response",
			}
		}
	}

	return nil
}

// ApiRequest is a helper method for making API requests to a specific endpoint with the
// given HTTP method, API version, authentication, and request body.
// It constructs the full API URL using the client's base URL and version,
// and then delegates the request handling to the `Request` method.
//
// Parameters:
//   - method: The HTTP method for the request, such as "GET" or "POST".
//   - endpoint: The API endpoint (path) relative to the base URL, such as "/orders".
//   - version: The API version to use. If empty, the default version is used.
//   - auth: A boolean indicating whether the request requires authentication.
//     If true, the method adds an Authorization header with the access token
//     and handles automatic token refresh if enabled.
//   - body: An optional request body. For GET requests, it is converted into URL
//     parameters; for POST requests, it is marshaled to JSON.
//   - result: A pointer to a variable where the response body should be unmarshaled.
//     If nil, the response body is not unmarshaled.
//
// Returns:
//   - An error if the request fails or if the response contains an error status code.
//
// Behavior:
//   - Constructs the full API URL by combining the base URL, API version, and endpoint.
//   - Calls the `Request` method with the constructed URL and passes along other parameters.
//   - The `Request` method handles authentication, token refresh, and error responses.
//
// Example:
//
//	var response MyResponseStruct
//	err := client.ApiRequest("POST", "/orders", "v1", true, requestBody, &response)
//	if err != nil {
//	    log.Fatalf("API request failed: %v", err)
//	}
//
// Dependencies:
//   - `createApiURI` for constructing the full API URL.
//   - `Request` for handling the HTTP request and processing the response.
func (c *Client) ApiRequest(method, endpoint string, version string, auth bool, body interface{}, result interface{}) error {
	url := c.createApiURI(endpoint, version)
	return c.Request(method, url, auth, body, result)
}

// Authenticate authenticates the client using the provided API key and secret key.
// It sends a POST request to the authentication endpoint and retrieves the access
// and refresh tokens for the client.
//
// Parameters:
//   - apiKey: The API key used for authentication.
//   - secretKey: The secret key used for authentication.
//
// Returns:
//   - A pointer to an `AuthenticationResponse` struct containing the access and
//     refresh tokens if authentication is successful.
//   - An error if authentication fails or if there are issues with the request.
//
// Behavior:
//   - Validates that both `apiKey` and `secretKey` are provided. Returns an error
//     if either is empty.
//   - Sends a POST request to the `/usr/authenticate/` endpoint with the API key
//     and secret key in the request body.
//   - If the request succeeds, updates the client's `AccessToken` and `RefreshToken`
//     with the tokens from the response.
//   - If the request fails, checks for specific API errors (e.g., 401 or 429) and
//     returns detailed error messages. For other errors, wraps and returns them.
//
// Example:
//
//	authResponse, err := client.Authenticate("your-api-key", "your-secret-key")
//	if err != nil {
//	    log.Fatalf("Authentication failed: %v", err)
//	}
//	fmt.Printf("Access Token: %s\n", authResponse.Access)
//
// Dependencies:
//   - Calls `ApiRequest` to send the authentication request.
//   - Relies on `APIError` for structured error handling.
//
// Errors:
//   - "API key and/or secret key are empty" if either `apiKey` or `secretKey` is missing.
//   - "authentication failed: invalid API key or secret key" for 401 Unauthorized responses.
//   - "authentication failed: rate limit exceeded" for 429 Too Many Requests responses.
//   - "authentication failed: %v" for other API or request errors.
func (c *Client) Authenticate(apiKey, secretKey string) (*t.AuthenticationResponse, error) {
	if apiKey == "" || secretKey == "" {
		return nil, &GoBitpinError{
			Message: "API key and/or secret key are empty",
			Err:     nil,
		}
	}

	reqBody := map[string]string{
		"api_key":    apiKey,
		"secret_key": secretKey,
	}

	var authResponse t.AuthenticationResponse
	err := c.ApiRequest("POST", "/usr/authenticate/", Version, false, reqBody, &authResponse)

	if err != nil {
		// Check for specific API errors here
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			switch apiErr.StatusCode {
			case 401:
				return nil, err
			case 429:
				return nil, err
			default:
				return nil, err
			}
		}
		return nil, err
	}

	// Update the client's tokens with the newly received ones
	c.AccessToken = authResponse.Access
	c.RefreshToken = authResponse.Refresh

	return &authResponse, nil
}

// RefreshAccessToken refreshes the client's access token using the current refresh token.
// It sends a POST request to the `/usr/refresh_token/` endpoint and retrieves a new
// access token for the client.
//
// Returns:
//   - An error if the request fails or if the refresh token is invalid or expired.
//     Returns nil if the access token is successfully refreshed.
//
// Behavior:
//   - Sends a POST request with the current refresh token in the request body.
//   - Updates the client's `AccessToken` with the new token from the response.
//
// Example:
//
//	err := client.RefreshAccessToken()
//	if err != nil {
//	    log.Fatalf("Failed to refresh access token: %v", err)
//	}
//
// Dependencies:
//   - Calls `ApiRequest` to handle the HTTP request and process the response.
//
// Errors:
//   - Returns an error if the refresh token is invalid or expired, or if there are
//     issues with the request.
//
// Example Request Body:
//
//	{
//	    "refresh": "<current-refresh-token>"
//	}
//
// Example Response Body:
//
//	{
//	    "access": "<new-access-token>"
//	}
func (c *Client) RefreshAccessToken() error {
	reqBody := map[string]string{
		"refresh": c.RefreshToken,
	}

	var refreshResponse t.RefreshTokenResponse
	err := c.ApiRequest("POST", "/usr/refresh_token/", Version, false, reqBody, &refreshResponse)
	if err != nil {
		return err
	}

	// Update the bitpin_client's access token with the newly received one
	c.AccessToken = refreshResponse.Access

	return nil
}

// GetCurrencies retrieves a list of available currencies from the API.
// It sends a GET request to the `/mkt/currencies/` endpoint and returns
// the list of currencies.
//
// Returns:
//   - A pointer to a `Currencies` struct containing the list of currencies
//     available in the market.
//   - An error if the request fails or the response cannot be processed.
//
// Behavior:
//   - Sends a GET request to the `/mkt/currencies/` endpoint.
//   - Does not require authentication (`auth` is set to false).
//   - Unmarshals the response into a `Currencies` struct.
//
// Example:
//
//	currencies, err := client.GetCurrencies()
//	if err != nil {
//	    log.Fatalf("Failed to fetch currencies: %v", err)
//	}
//	for _, currency := range *currencies {
//	    fmt.Printf("Currency: %s, Name: %s\n", currency.Currency, currency.Name)
//	}
//
// Dependencies:
//   - Relies on `ApiRequest` for HTTP request handling and response processing.
//
// Errors:
//   - "error fetching currencies: %v" if the request fails or the response
//     cannot be unmarshaled.
//
// Example Response:
//
//	[
//	    {
//	        "currency": "BTC",
//	        "name": "Bitcoin",
//	        "tradable": true,
//	        "precision": "8"
//	    },
//	    {
//	        "currency": "ETH",
//	        "name": "Ethereum",
//	        "tradable": true,
//	        "precision": "8"
//	    }
//	]
func (c *Client) GetCurrencies() (*t.Currencies, error) {
	var currencies *t.Currencies
	err := c.ApiRequest("GET", "/mkt/currencies/", Version, false, nil, &currencies)
	if err != nil {
		return nil, err
	}
	return currencies, nil
}

// GetMarkets retrieves a list of available markets from the API.
// It sends a GET request to the `/mkt/markets/` endpoint and returns
// the list of markets.
//
// Returns:
//   - A pointer to a `Markets` struct containing the list of trading markets
//     available in the platform.
//   - An error if the request fails or the response cannot be processed.
//
// Behavior:
//   - Sends a GET request to the `/mkt/markets/` endpoint.
//   - Does not require authentication (`auth` is set to false).
//   - Unmarshals the response into a `Markets` struct.
//
// Example:
//
//	markets, err := client.GetMarkets()
//	if err != nil {
//	    log.Fatalf("Failed to fetch markets: %v", err)
//	}
//	for _, market := range *markets {
//	    fmt.Printf("Market: %s, Base: %s, Quote: %s\n", market.Symbol, market.Base, market.Quote)
//	}
//
// Dependencies:
//   - Relies on `ApiRequest` for HTTP request handling and response processing.
//
// Errors:
//   - "error fetching markets: %v" if the request fails or the response
//     cannot be unmarshaled.
//
// Example Response:
//
//	[
//	    {
//	        "symbol": "BTC_USDT",
//	        "name": "Bitcoin/USDT",
//	        "base": "BTC",
//	        "quote": "USDT",
//	        "tradable": true,
//	        "price_precision": 2,
//	        "base_amount_precision": 8,
//	        "quote_amount_precision": 2
//	    },
//	    {
//	        "symbol": "ETHUSDT",
//	        "name": "Ethereum/USDT",
//	        "base": "ETH",
//	        "quote": "USDT",
//	        "tradable": true,
//	        "price_precision": 2,
//	        "base_amount_precision": 8,
//	        "quote_amount_precision": 2
//	    }
//	]
func (c *Client) GetMarkets() (*t.Markets, error) {
	var markets *t.Markets
	err := c.ApiRequest("GET", "/mkt/markets/", Version, false, nil, &markets)
	if err != nil {
		return nil, err
	}
	return markets, nil
}

// GetTickers retrieves a list of market tickers from the API.
// It sends a GET request to the `/mkt/tickers/` endpoint and returns
// real-time ticker information for available markets.
//
// Returns:
//   - A pointer to a `Tickers` struct containing a list of tickers with
//     real-time price and trading data for various markets.
//   - An error if the request fails or the response cannot be processed.
//
// Behavior:
//   - Sends a GET request to the `/mkt/tickers/` endpoint.
//   - Does not require authentication (`auth` is set to false).
//   - Unmarshals the response into a `Tickers` struct.
//
// Example:
//
//	tickers, err := client.GetTickers()
//	if err != nil {
//	    log.Fatalf("Failed to fetch tickers: %v", err)
//	}
//	for _, ticker := range *tickers {
//	    fmt.Printf("Symbol: %s, Price: %s, Daily Change: %.2f\n",
//	        ticker.Symbol, ticker.Price, ticker.DailyChangePrice)
//	}
//
// Dependencies:
//   - Relies on `ApiRequest` for HTTP request handling and response processing.
//
// Errors:
//   - "error fetching tickers: %v" if the request fails or the response
//     cannot be unmarshaled.
//
// Example Response:
//
//	[
//	    {
//	        "symbol": "BTC_USDT",
//	        "price": "40000.00",
//	        "daily_change_price": -200.00,
//	        "low": "39500.00",
//	        "high": "40500.00",
//	        "timestamp": 1625247600
//	    },
//	    {
//	        "symbol": "ETHUSDT",
//	        "price": "2500.00",
//	        "daily_change_price": 50.00,
//	        "low": "2450.00",
//	        "high": "2550.00",
//	        "timestamp": 1625247600
//	    }
//	]
func (c *Client) GetTickers() (*t.Tickers, error) {
	var tickers *t.Tickers
	err := c.ApiRequest("GET", "/mkt/tickers/", Version, false, nil, &tickers)
	if err != nil {
		return nil, err
	}
	return tickers, nil
}

// GetOrderBook retrieves the order book for a specific trading symbol from the API.
// It sends a GET request to the `/mth/orderbook/<symbol>/` endpoint and returns
// detailed order book information, including asks and bids.
//
// Parameters:
//   - symbol: A string representing the trading symbol, such as "BTC_USDT".
//
// Returns:
//   - A pointer to an `OrderBook` struct containing the asks and bids for the
//     specified market symbol.
//   - An error if the request fails or the response cannot be processed.
//
// Behavior:
//   - Sends a GET request to the `/mth/orderbook/<symbol>/` endpoint.
//   - Does not require authentication (`auth` is set to false).
//   - Unmarshals the response into an `OrderBook` struct.
//
// Example:
//
//	orderBook, err := client.GetOrderBook("BTC_USDT")
//	if err != nil {
//	    log.Fatalf("Failed to fetch order book: %v", err)
//	}
//	fmt.Printf("Asks: %v\n", orderBook.Asks)
//	fmt.Printf("Bids: %v\n", orderBook.Bids)
//
// Dependencies:
//   - Relies on `ApiRequest` for HTTP request handling and response processing.
//
// Errors:
//   - "error fetching order book: %v" if the request fails or the response
//     cannot be unmarshaled.
//
// Example Response:
//
//	{
//	    "asks": [["40000.00", "0.5"], ["40010.00", "0.2"]],
//	    "bids": [["39990.00", "0.3"], ["39980.00", "1.0"]]
//	}
func (c *Client) GetOrderBook(symbol string) (*t.OrderBook, error) {
	var orderBook *t.OrderBook
	err := c.ApiRequest("GET", fmt.Sprintf("/mth/orderbook/%s/", symbol), Version, false, nil, &orderBook)
	if err != nil {
		return nil, err
	}
	return orderBook, nil
}

// GetRecentTrades retrieves the most recent trades for a specific trading symbol from the API.
// It sends a GET request to the `/mth/matches/<symbol>/` endpoint and returns
// a list of recent trades.
//
// Parameters:
//   - symbol: A string representing the trading symbol, such as "BTC_USDT".
//
// Returns:
//   - A pointer to a slice of `Trade` structs, each representing a recent trade for the
//     specified market symbol.
//   - An error if the request fails or the response cannot be processed.
//
// Behavior:
//   - Sends a GET request to the `/mth/matches/<symbol>/` endpoint.
//   - Does not require authentication (`auth` is set to false).
//   - Unmarshals the response into a slice of `Trade` structs.
//
// Example:
//
//	trades, err := client.GetRecentTrades("BTC_USDT")
//	if err != nil {
//	    log.Fatalf("Failed to fetch recent trades: %v", err)
//	}
//	for _, trade := range *trades {
//	    fmt.Printf("Trade ID: %d, Price: %s, Amount: %s\n", trade.Id, trade.Price, trade.BaseAmount)
//	}
//
// Dependencies:
//   - Relies on `ApiRequest` for HTTP request handling and response processing.
//
// Errors:
//   - "error fetching recent trades: %v" if the request fails or the response
//     cannot be unmarshaled.
//
// Example Response:
//
//	[
//	    {
//	        "id": 12345,
//	        "symbol": "BTC_USDT",
//	        "base_amount": "0.01",
//	        "quote_amount": "400.00",
//	        "price": "40000.00",
//	        "created_at": "2023-01-01T12:00:00Z",
//	        "commission": "0.01",
//	        "side": "buy",
//	        "commission_currency": "BTC",
//	        "order_id": 54321,
//	        "identifier": "abc123"
//	    },
//	    {
//	        "id": 12346,
//	        "symbol": "BTC_USDT",
//	        "base_amount": "0.02",
//	        "quote_amount": "800.00",
//	        "price": "40000.00",
//	        "created_at": "2023-01-01T12:01:00Z",
//	        "commission": "0.02",
//	        "side": "sell",
//	        "commission_currency": "BTC",
//	        "order_id": 54322,
//	        "identifier": "xyz789"
//	    }
//	]
func (c *Client) GetRecentTrades(symbol string) (*[]*t.Trade, error) {
	var trades *[]*t.Trade
	err := c.ApiRequest("GET", fmt.Sprintf("/mth/matches/%s/", symbol), Version, false, nil, &trades)
	if err != nil {
		return nil, err
	}
	return trades, nil
}

// GetWallets retrieves a list of wallets for the authenticated user from the API.
// It sends a GET request to the `/wlt/wallets/` endpoint and returns wallet information
// based on the provided parameters.
//
// Parameters:
//   - params: A `GetWalletParams` struct containing optional filters for querying
//     specific wallets, such as by asset, service, offset, or limit.
//
// Returns:
//   - A pointer to a `Wallets` struct containing the list of wallets with details
//     like balances and associated services.
//   - An error if the request fails, the user is not authenticated, or the response
//     cannot be processed.
//
// Behavior:
//   - Sends a GET request to the `/wlt/wallets/` endpoint with optional query parameters.
//   - Requires authentication (`auth` is set to true).
//   - Unmarshals the response into a `Wallets` struct.
//
// Example:
//
//	params := t.GetWalletParams{
//	    Assets:  []string{"BTC", "USDT"},
//	    Service: "spot",
//	    Limit:   10,
//	    Offset:  0,
//	}
//	wallets, err := client.GetWallets(params)
//	if err != nil {
//	    log.Fatalf("Failed to fetch wallets: %v", err)
//	}
//	for _, wallet := range *wallets {
//	    fmt.Printf("Asset: %s, Balance: %s, Frozen: %s\n", wallet.Asset, wallet.Balance, wallet.Frozen)
//	}
//
// Dependencies:
//   - Relies on `ApiRequest` for HTTP request handling and response processing.
//
// Errors:
//   - "error fetching wallets: %v" if the request fails or the response
//     cannot be unmarshaled.
//   - Returns authentication errors if the client is not properly authenticated.
//
// Example Response:
//
//	[
//	    {
//	        "id": 1,
//	        "asset": "BTC",
//	        "balance": "0.5",
//	        "frozen": "0.1",
//	        "service": "spot"
//	    },
//	    {
//	        "id": 2,
//	        "asset": "USDT",
//	        "balance": "1000.0",
//	        "frozen": "100.0",
//	        "service": "futures"
//	    }
//	]
func (c *Client) GetWallets(params t.GetWalletParams) (*t.Wallets, error) {
	var wallets *t.Wallets
	err := c.ApiRequest("GET", "/wlt/wallets/", Version, true, params, &wallets)
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

// CreateOrder submits a new order to the API based on the provided parameters.
// It sends a POST request to the `/odr/orders/` endpoint and returns the status
// of the created order.
//
// Parameters:
//   - params: A `CreateOrderParams` struct containing details about the order,
//     such as the trading symbol, order type, side (buy/sell), price, and amounts.
//
// Returns:
//   - A pointer to an `OrderStatus` struct containing the details and status of
//     the created order.
//   - An error if the request fails, the user is not authenticated, or the response
//     cannot be processed.
//
// Behavior:
//   - Sends a POST request to the `/odr/orders/` endpoint with the order details in the body.
//   - Requires authentication (`auth` is set to true).
//   - Unmarshals the response into an `OrderStatus` struct.
//
// Example:
//
//	params := t.CreateOrderParams{
//	    Symbol:     "BTC_USDT",
//	    Type:       "limit",
//	    Side:       "buy",
//	    Price:      "40000",
//	    BaseAmount: "0.01",
//	}
//	orderStatus, err := client.CreateOrder(params)
//	if err != nil {
//	    log.Fatalf("Failed to create order: %v", err)
//	}
//	fmt.Printf("Order ID: %d, Status: %s, Price: %s\n", orderStatus.Id, orderStatus.State, orderStatus.Price)
//
// Dependencies:
//   - Relies on `ApiRequest` for HTTP request handling and response processing.
//
// Errors:
//   - "error creating order: %v" if the request fails or the response cannot be unmarshaled.
//   - Returns authentication errors if the client is not properly authenticated.
//
// Example Request Body:
//
//	{
//	    "symbol": "BTC_USDT",
//	    "type": "limit",
//	    "side": "buy",
//	    "price": "40000",
//	    "base_amount": "0.01",
//	    "quote_amount": null
//	}
//
// Example Response:
//
//	{
//	    "id": 123456,
//	    "symbol": "BTC_USDT",
//	    "type": "limit",
//	    "side": "buy",
//	    "base_amount": "0.01",
//	    "quote_amount": "400.00",
//	    "price": "40000",
//	    "stop_price": null,
//	    "oco_target_price": null,
//	    "identifier": "user123",
//	    "state": "open",
//	    "created_at": "2023-01-01T12:00:00Z",
//	    "closed_at": null,
//	    "dealed_base_amount": "0.0",
//	    "dealed_quote_amount": "0.0",
//	    "req_to_cancel": false,
//	    "commission": "0.01"
//	}
func (c *Client) CreateOrder(params t.CreateOrderParams) (*t.OrderStatus, error) {
	var orderStatus *t.OrderStatus
	err := c.ApiRequest("POST", "/odr/orders/", Version, true, params, &orderStatus)
	if err != nil {
		return nil, err
	}
	return orderStatus, nil
}

// CancelOrder cancels an active order by its order ID.
// It sends a DELETE request to the `/odr/orders/<orderId>/` endpoint and returns an error
// if the cancellation fails.
//
// Parameters:
//   - orderId: The unique identifier of the order to be canceled.
//
// Returns:
//   - An error if the request fails, the user is not authenticated, or the cancellation
//     could not be processed. Returns nil if the order is successfully canceled.
//
// Behavior:
//   - Sends a DELETE request to the `/odr/orders/<orderId>/` endpoint with the order ID.
//   - Requires authentication (`auth` is set to true).
//   - If successful, the order will be canceled in the system.
//
// Example:
//
//	err := client.CancelOrder(123456)
//	if err != nil {
//	    log.Fatalf("Failed to cancel order: %v", err)
//	}
//
// Dependencies:
//   - Relies on `ApiRequest` for HTTP request handling and response processing.
//
// Errors:
//   - "error canceling order: %v" if the request fails or the response indicates an error.
//
// Example Response (for successful cancellation):
//
//	HTTP Status 204 No Content
//
// Example Response (if the order is not found):
//
//	HTTP Status 404 Not Found
func (c *Client) CancelOrder(orderId int) error {
	err := c.ApiRequest("DELETE", fmt.Sprintf("/odr/orders/%d/", orderId), Version, true, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetOrdersHistory retrieves the order history for the authenticated user.
// It sends a GET request to the `/odr/orders/` endpoint and returns a list of orders
// based on the provided filters.
//
// Parameters:
//   - params: A `GetOrdersHistoryParams` struct containing optional filters for the
//     order history query, such as symbol, side, state, type, start and end dates,
//     and pagination parameters like offset and limit.
//
// Returns:
//   - A pointer to an `OrderStatuses` struct containing the list of orders in the
//     user's order history.
//   - An error if the request fails, the user is not authenticated, or the response
//     cannot be processed.
//
// Behavior:
//   - Sends a GET request to the `/odr/orders/` endpoint with the specified filters.
//   - Requires authentication (`auth` is set to true).
//   - Unmarshals the response into an `OrderStatuses` struct.
//
// Example:
//
//	params := t.GetOrdersHistoryParams{
//	    Symbol:    "BTC_USDT",
//	    State:     "closed",
//	    Limit:     10,
//	    Offset:    0,
//	}
//	orders, err := client.GetOrdersHistory(params)
//	if err != nil {
//	    log.Fatalf("Failed to fetch order history: %v", err)
//	}
//	for _, order := range *orders {
//	    fmt.Printf("Order ID: %d, Status: %s, Price: %s\n", order.Id, order.State, order.Price)
//	}
//
// Dependencies:
//   - Relies on `ApiRequest` for HTTP request handling and response processing.
//
// Errors:
//   - "error fetching order history: %v" if the request fails or the response
//     cannot be unmarshaled.
//
// Example Response:
//
//	[
//	    {
//	        "id": 123456,
//	        "symbol": "BTC_USDT",
//	        "base_amount": "0.01",
//	        "quote_amount": "400.00",
//	        "price": "40000.00",
//	        "side": "buy",
//	        "state": "closed",
//	        "created_at": "2023-01-01T12:00:00Z",
//	        "closed_at": "2023-01-01T12:05:00Z",
//	        "commission": "0.01",
//	        "commission_currency": "BTC",
//	        "order_id": 654321,
//	        "identifier": "user123"
//	    },
//	    {
//	        "id": 123457,
//	        "symbol": "BTC_USDT",
//	        "base_amount": "0.02",
//	        "quote_amount": "800.00",
//	        "price": "40000.00",
//	        "side": "sell",
//	        "state": "closed",
//	        "created_at": "2023-01-01T13:00:00Z",
//	        "closed_at": "2023-01-01T13:10:00Z",
//	        "commission": "0.02",
//	        "commission_currency": "BTC",
//	        "order_id": 654322,
//	        "identifier": "user456"
//	    }
//	]
func (c *Client) GetOrdersHistory(params t.GetOrdersHistoryParams) (*t.OrderStatuses, error) {
	var orders *t.OrderStatuses
	err := c.ApiRequest("GET", "/odr/orders/", Version, true, params, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetOpenOrders retrieves a list of active (open) orders for the authenticated user.
// It sends a GET request to the `/odr/orders/` endpoint with the state filter set to "active"
// and returns the list of open orders based on the provided parameters.
//
// Parameters:
//   - params: A `GetOrdersHistoryParams` struct containing optional filters for querying
//     specific open orders, such as symbol, side, type, start and end dates, and pagination
//     parameters like offset and limit. The `State` field is automatically set to "active".
//
// Returns:
//   - A pointer to an `OrderStatuses` struct containing the list of open orders for the user.
//   - An error if the request fails, the user is not authenticated, or the response
//     cannot be processed.
//
// Behavior:
//   - Sends a GET request to the `/odr/orders/` endpoint with the provided filters.
//   - The `State` parameter is automatically set to "active" to filter for open orders.
//   - Requires authentication (`auth` is set to true).
//   - Unmarshals the response into an `OrderStatuses` struct.
//
// Example:
//
//	params := t.GetOrdersHistoryParams{
//	    Symbol:    "BTC_USDT",
//	    Side:      "buy",
//	    Limit:     10,
//	    Offset:    0,
//	}
//	openOrders, err := client.GetOpenOrders(params)
//	if err != nil {
//	    log.Fatalf("Failed to fetch open orders: %v", err)
//	}
//	for _, order := range *openOrders {
//	    fmt.Printf("Order ID: %d, Status: %s, Price: %s\n", order.Id, order.State, order.Price)
//	}
//
// Dependencies:
//   - Relies on `ApiRequest` for HTTP request handling and response processing.
//
// Errors:
//   - "error fetching order history: %v" if the request fails or the response
//     cannot be unmarshaled.
//
// Example Response:
//
//	[
//	    {
//	        "id": 123456,
//	        "symbol": "BTC_USDT",
//	        "base_amount": "0.01",
//	        "quote_amount": "400.00",
//	        "price": "40000.00",
//	        "side": "buy",
//	        "state": "active",
//	        "created_at": "2023-01-01T12:00:00Z",
//	        "closed_at": null,
//	        "commission": "0.01",
//	        "commission_currency": "BTC",
//	        "order_id": 654321,
//	        "identifier": "user123"
//	    },
//	    {
//	        "id": 123457,
//	        "symbol": "BTC_USDT",
//	        "base_amount": "0.02",
//	        "quote_amount": "800.00",
//	        "price": "40000.00",
//	        "side": "sell",
//	        "state": "active",
//	        "created_at": "2023-01-01T13:00:00Z",
//	        "closed_at": null,
//	        "commission": "0.02",
//	        "commission_currency": "BTC",
//	        "order_id": 654322,
//	        "identifier": "user456"
//	    }
//	]
func (c *Client) GetOpenOrders(params t.GetOrdersHistoryParams) (*t.OrderStatuses, error) {
	var orders *t.OrderStatuses
	params.State = "active" // Automatically filter for active (open) orders
	err := c.ApiRequest("GET", "/odr/orders/", Version, true, params, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetOrderStatuses retrieves the statuses of multiple orders using their order IDs.
// It sends a GET request to the `/odr/orders/<orderIds>/` endpoint and returns the
// statuses of the specified orders.
//
// Parameters:
//   - orderIds: A slice of strings representing the unique IDs of the orders whose
//     statuses are to be fetched.
//
// Returns:
//   - A pointer to an `OrderStatus` struct containing the status and details of
//     the specified orders.
//   - An error if the request fails, the user is not authenticated, or the response
//     cannot be processed.
//
// Behavior:
//   - Sends a GET request to the `/odr/orders/<orderIds>/` endpoint, where the
//     `orderIds` are joined into a comma-separated string.
//   - Requires authentication (`auth` is set to true).
//   - Unmarshals the response into an `OrderStatus` struct.
//
// Example:
//
//	orderIds := []string{"123456", "789012"}
//	orderStatuses, err := client.GetOrderStatuses(orderIds)
//	if err != nil {
//	    log.Fatalf("Failed to fetch order statuses: %v", err)
//	}
//	fmt.Printf("Order ID: %s, Status: %s, Price: %s\n", orderStatuses.Id, orderStatuses.State, orderStatuses.Price)
//
// Dependencies:
//   - Relies on `ApiRequest` for HTTP request handling and response processing.
//
// Errors:
//   - "error fetching order statuses: %v" if the request fails or the response
//     cannot be unmarshaled.
//
// Example Response:
//
//	{
//	    "id": 123456,
//	    "symbol": "BTC_USDT",
//	    "base_amount": "0.01",
//	    "quote_amount": "400.00",
//	    "price": "40000.00",
//	    "side": "buy",
//	    "state": "closed",
//	    "created_at": "2023-01-01T12:00:00Z",
//	    "closed_at": "2023-01-01T12:05:00Z",
//	    "commission": "0.01",
//	    "commission_currency": "BTC",
//	    "order_id": 654321,
//	    "identifier": "user123"
//	}
func (c *Client) GetOrderStatuses(orderIds []string) (*t.OrderStatus, error) {
	var orders *t.OrderStatus
	err := c.ApiRequest("GET", fmt.Sprintf("/odr/orders/%v/", strings.Join(orderIds, ",")), Version, true, nil, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetUserTrades retrieves a list of trades made by the authenticated user.
// It sends a GET request to the `/odr/fills/` endpoint and returns the list of user trades
// based on the provided parameters.
//
// Parameters:
//   - params: A `GetUserTradesParams` struct containing optional filters for querying
//     specific trades, such as symbol, side (buy/sell), offset, limit, and other parameters.
//
// Returns:
//   - A pointer to a `UserTrades` struct containing the list of trades made by the user.
//   - An error if the request fails, the user is not authenticated, or the response
//     cannot be processed.
//
// Behavior:
//   - Sends a GET request to the `/odr/fills/` endpoint with the provided filters.
//   - Requires authentication (`auth` is set to true).
//   - Unmarshals the response into a `UserTrades` struct containing the user's trade history.
//
// Example:
//
//	params := t.GetUserTradesParams{
//	    Symbol: "BTC_USDT",
//	    Side:   "buy",
//	    Limit:  10,
//	    Offset: 0,
//	}
//	trades, err := client.GetUserTrades(params)
//	if err != nil {
//	    log.Fatalf("Failed to fetch user trades: %v", err)
//	}
//	for _, trade := range *trades {
//	    fmt.Printf("Trade ID: %s, Price: %s, Amount: %s\n", trade.Id, trade.Price, trade.BaseAmount)
//	}
//
// Dependencies:
//   - Relies on `ApiRequest` for HTTP request handling and response processing.
//
// Errors:
//   - "error fetching user trades: %v" if the request fails or the response
//     cannot be unmarshaled.
//
// Example Response:
//
//	[
//	    {
//	        "id": "12345",
//	        "symbol": "BTC_USDT",
//	        "base_amount": "0.01",
//	        "quote_amount": "400.00",
//	        "price": "40000.00",
//	        "created_at": "2023-01-01T12:00:00Z",
//	        "commission": "0.01",
//	        "side": "buy",
//	        "commission_currency": "BTC",
//	        "order_id": 54321,
//	        "identifier": "abc123"
//	    },
//	    {
//	        "id": "12346",
//	        "symbol": "BTC_USDT",
//	        "base_amount": "0.02",
//	        "quote_amount": "800.00",
//	        "price": "40000.00",
//	        "created_at": "2023-01-01T12:01:00Z",
//	        "commission": "0.02",
//	        "side": "sell",
//	        "commission_currency": "BTC",
//	        "order_id": 54322,
//	        "identifier": "xyz789"
//	    }
//	]
func (c *Client) GetUserTrades(params t.GetUserTradesParams) (*t.UserTrades, error) {
	var trades *t.UserTrades
	err := c.ApiRequest("GET", "/odr/fills/", Version, true, params, &trades)
	if err != nil {
		return nil, err
	}
	return trades, nil
}
