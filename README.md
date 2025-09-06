# Go Bitpin

[![Go Reference](https://pkg.go.dev/badge/github.com/amiwrpremium/go-bitpin.svg)](https://pkg.go.dev/github.com/amiwrpremium/go-bitpin)
[![Go Report Card](https://goreportcard.com/badge/github.com/amiwrpremium/go-bitpin)](https://goreportcard.com/report/github.com/amiwrpremium/go-bitpin)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/amiwrpremium/go-bitpin)](https://golang.org/dl/)

A comprehensive Go SDK for interacting with the [Bitpin](https://bitpin.ir/) cryptocurrency exchange API. This SDK provides a simple and intuitive way to integrate Bitpin's trading functionality into your Go applications.

## Disclaimer

This is an unofficial SDK and is not affiliated with, endorsed by, or connected to Bitpin. Use at your own risk. The author(s) of this SDK are not responsible for any losses incurred through its use.

## Features

- Complete implementation of Bitpin's REST API
- Comprehensive market data access
- Order management (create, cancel, query)
- Wallet operations
- Authentication and token management
- Rate limiting and error handling
- Type-safe request/response structs

## Installation

```bash
go get [github.com/rzabhd80/go-sdk-bitpin](https://github.com/rzabhd80/go-sdk-bitpin/)
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/rzabhd80/go-sdk-bitpin"
)

func main() {
    // Initialize client with options
    opts := bitpin.ClientOptions{
        ApiKey:      "your-api-key",
        SecretKey:   "your-secret-key",
        AutoRefresh: true,
    }
    
    client, err := bitpin.NewClient(opts)
    if err != nil {
        panic(err)
    }

    // Get market information
    markets, err := client.GetMarkets()
    if err != nil {
        panic(err)
    }

    // Print available markets
    for _, market := range *markets {
        fmt.Printf("Market: %s, Base: %s, Quote: %s\n", 
            market.Symbol, market.Base, market.Quote)
    }
}
```

## Documentation

For detailed documentation and examples, please visit the [Go Package Documentation](https://pkg.go.dev/github.com/amiwrpremium/go-bitpin).

For information about Bitpin's API, visit their [official API documentation](https://docs.bitpin.ir/).

### Examples

#### Authentication
```go
client, err := bitpin.NewClient(bitpin.ClientOptions{
    ApiKey:    "your-api-key",
    SecretKey: "your-secret-key",
})
```

#### Get Wallet Information
```go
params := types.GetWalletParams{
    Assets: []string{"BTC", "USDT"},
    Limit:  10,
}

wallets, err := client.GetWallets(params)
if err != nil {
    panic(err)
}
```

#### Place an Order
```go
orderParams := types.CreateOrderParams{
    Symbol:     "BTC_USDT",
    Type:       "limit",
    Side:       "buy",
    Price:      "50000",
    BaseAmount: "0.001",
}

order, err := client.CreateOrder(orderParams)
if err != nil {
    panic(err)
}
```

## Error Handling

The SDK uses custom error types to provide detailed information about API errors:

```go
if err != nil {
    if apiErr, ok := err.(*bitpin.APIError); ok {
        fmt.Printf("API Error: Status %d, Message: %s\n", 
            apiErr.StatusCode, apiErr.Message)
    }
}
```

## Configuration Options

The SDK can be configured with various options:

```go
opts := bitpin.ClientOptions{
    HttpClient:   &http.Client{}, // Custom HTTP client
    Timeout:      10 * time.Second,
    BaseUrl:      "https://api.bitpin.market", // Custom base URL
    AutoAuth:     true,
    AutoAuth:     true,
    AutoRefresh:  true,
}
```

## More Examples

Comprehensive examples for all methods can be found in the [EXAMPLES.md](EXAMPLES.md) file.


## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have questions, please [open an issue](https://github.com/amiwrpremium/go-bitpin/issues) on GitHub.
