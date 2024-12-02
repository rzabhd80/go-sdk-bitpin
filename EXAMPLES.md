# Bitpin SDK Examples

This document provides examples for all available methods in the Bitpin Go SDK. Each example demonstrates proper usage and error handling.

## Table of Contents

- [Client Initialization](#client-initialization)
- [Authentication](#authentication)
- [Market Information](#market-information)
- [Trading Operations](#trading-operations)
- [Wallet Operations](#wallet-operations)
- [Error Handling](#error-handling)

## Client Initialization

### Basic Client
```go
package main

import (
    "github.com/amiwrpremium/go-bitpin"
    "time"
)

func main() {
    opts := bitpin.ClientOptions{
        ApiKey:      "your-api-key",
        SecretKey:   "your-secret-key",
        Timeout:     10 * time.Second,
        AutoRefresh: true,
    }
    
    client, err := bitpin.NewClient(opts)
    if err != nil {
        panic(err)
    }
}
```

### Custom HTTP Client
```go
httpClient := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        10,
        IdleConnTimeout:     30 * time.Second,
        DisableCompression:  true,
    },
}

opts := bitpin.ClientOptions{
    HttpClient:  httpClient,
    ApiKey:      "your-api-key",
    SecretKey:   "your-secret-key",
    AutoRefresh: true,
}

client, err := bitpin.NewClient(opts)
```

## Authentication

### Manual Authentication
```go
authResponse, err := client.Authenticate("your-api-key", "your-secret-key")
if err != nil {
    panic(err)
}

fmt.Printf("Access Token: %s\n", authResponse.Access)
fmt.Printf("Refresh Token: %s\n", authResponse.Refresh)
```

### Refresh Access Token
```go
err := client.RefreshAccessToken()
if err != nil {
    panic(err)
}
```

## Market Information

### Get Currencies
```go
currencies, err := client.GetCurrencies()
if err != nil {
    panic(err)
}

for _, currency := range *currencies {
    fmt.Printf("Currency: %s\n", currency.Currency)
    fmt.Printf("Name: %s\n", currency.Name)
    fmt.Printf("Tradable: %v\n", currency.Tradable)
    fmt.Printf("Precision: %s\n\n", currency.Precision)
}
```

### Get Markets
```go
markets, err := client.GetMarkets()
if err != nil {
    panic(err)
}

for _, market := range *markets {
    fmt.Printf("Symbol: %s\n", market.Symbol)
    fmt.Printf("Base/Quote: %s/%s\n", market.Base, market.Quote)
    fmt.Printf("Price Precision: %d\n", market.PricePrecision)
    fmt.Printf("Base Amount Precision: %d\n\n", market.BaseAmountPrecision)
}
```

### Get Tickers
```go
tickers, err := client.GetTickers()
if err != nil {
    panic(err)
}

for _, ticker := range *tickers {
    fmt.Printf("Symbol: %s\n", ticker.Symbol)
    fmt.Printf("Price: %s\n", ticker.Price)
    fmt.Printf("24h Change: %f\n", ticker.DailyChangePrice)
    fmt.Printf("High: %s\n", ticker.High)
    fmt.Printf("Low: %s\n\n", ticker.Low)
}
```

### Get Order Book
```go
orderBook, err := client.GetOrderBook("BTC_USDT")
if err != nil {
    panic(err)
}

fmt.Println("Asks:")
for _, ask := range orderBook.Asks {
    fmt.Printf("Price: %s, Amount: %s\n", ask[0], ask[1])
}

fmt.Println("\nBids:")
for _, bid := range orderBook.Bids {
    fmt.Printf("Price: %s, Amount: %s\n", bid[0], bid[1])
}
```

### Get Recent Trades
```go
trades, err := client.GetRecentTrades("BTC_USDT")
if err != nil {
    panic(err)
}

for _, trade := range *trades {
    fmt.Printf("ID: %s\n", trade.Id)
    fmt.Printf("Price: %s\n", trade.Price)
    fmt.Printf("Base Amount: %s\n", trade.BaseAmount)
    fmt.Printf("Quote Amount: %s\n", trade.QuoteAmount)
    fmt.Printf("Side: %s\n\n", trade.Side)
}
```

## Trading Operations

### Create Order
```go
// Limit Order
limitOrderParams := types.CreateOrderParams{
    Symbol:     "BTC_USDT",
    Type:       "limit",
    Side:       "buy",
    Price:      "50000",
    BaseAmount: "0.001",
}

limitOrder, err := client.CreateOrder(limitOrderParams)
if err != nil {
    panic(err)
}

// Market Order
marketOrderParams := types.CreateOrderParams{
    Symbol:      "BTC_USDT",
    Type:        "market",
    Side:        "sell",
    QuoteAmount: "100", // Selling BTC worth 100 USDT
}

marketOrder, err := client.CreateOrder(marketOrderParams)
if err != nil {
    panic(err)
}
```

### Cancel Order
```go
err := client.CancelOrder(123456)
if err != nil {
    panic(err)
}
```

### Get Order History
```go
params := types.GetOrdersHistoryParams{
    Symbol: "BTC_USDT",
}

orders, err := client.GetOrdersHistory(params)
if err != nil {
    panic(err)
}

for _, order := range *orders {
    fmt.Printf("Order ID: %d\n", order.Id)
    fmt.Printf("Symbol: %s\n", order.Symbol)
    fmt.Printf("Type: %s\n", order.Type)
    fmt.Printf("Side: %s\n", order.Side)
    fmt.Printf("Price: %s\n", order.Price)
    fmt.Printf("Status: %s\n\n", order.State)
}
```

### Get Open Orders
```go
params := types.GetOrdersHistoryParams{
    Symbol: "BTC_USDT",
}

openOrders, err := client.GetOpenOrders(params)
if err != nil {
    panic(err)
}

for _, order := range *openOrders {
    fmt.Printf("Order ID: %d\n", order.Id)
    fmt.Printf("Symbol: %s\n", order.Symbol)
    fmt.Printf("Price: %s\n", order.Price)
    fmt.Printf("Amount: %s\n\n", order.BaseAmount)
}
```

### Get Order Status
```go
orderIds := []string{"123456", "789012"}
status, err := client.GetOrderStatuses(orderIds)
if err != nil {
    panic(err)
}

fmt.Printf("Order ID: %d\n", status.Id)
fmt.Printf("State: %s\n", status.State)
fmt.Printf("Filled Amount: %s\n", status.DealedBaseAmount)
```

### Get User Trades
```go
params := types.GetUserTradesParams{
    Symbol: "BTC_USDT",
    Offset: 0,
}

trades, err := client.GetUserTrades(params)
if err != nil {
    panic(err)
}

for _, trade := range *trades {
    fmt.Printf("Trade ID: %d\n", trade.Id)
    fmt.Printf("Symbol: %s\n", trade.Symbol)
    fmt.Printf("Price: %s\n", trade.Price)
    fmt.Printf("Amount: %s\n", trade.BaseAmount)
    fmt.Printf("Side: %s\n\n", trade.Side)
}
```

## Wallet Operations

### Get Wallets
```go
params := types.GetWalletParams{
    Assets:  []string{"BTC", "USDT"},
}

wallets, err := client.GetWallets(params)
if err != nil {
    panic(err)
}

for _, wallet := range *wallets {
    fmt.Printf("Asset: %s\n", wallet.Asset)
    fmt.Printf("Balance: %s\n", wallet.Balance)
    fmt.Printf("Frozen: %s\n", wallet.Frozen)
    fmt.Printf("Service: %s\n\n", wallet.Service)
}
```

## Error Handling

### API Error Handling
```go
if err != nil {
    switch e := err.(type) {
    case *bitpin.APIError:
        fmt.Printf("API Error - Status: %d, Message: %s\n", 
            e.StatusCode, e.Message)
        
        switch e.StatusCode {
        case 401:
            fmt.Println("Authentication failed. Please check your credentials.")
        case 403:
            fmt.Println("Permission denied. Please check your API key permissions.")
        case 429:
            fmt.Println("Rate limit exceeded. Please wait before trying again.")
        default:
            fmt.Printf("Unexpected API error: %v\n", e)
        }
    default:
        fmt.Printf("Non-API error occurred: %v\n", err)
    }
}
```

### Complete Error Handling Example
```go
func placeOrder(client *bitpin.Client, symbol string, side string, amount string) {
    params := types.CreateOrderParams{
        Symbol:     symbol,
        Type:       "market",
        Side:       side,
        BaseAmount: amount,
    }
    
    order, err := client.CreateOrder(params)
    if err != nil {
        if apiErr, ok := err.(*bitpin.APIError); ok {
            switch apiErr.StatusCode {
            case 400:
                fmt.Printf("Invalid order parameters: %s\n", apiErr.Message)
            case 401:
                fmt.Println("Authentication failed. Refreshing tokens...")
                if err := client.RefreshAccessToken(); err != nil {
                    fmt.Printf("Token refresh failed: %v\n", err)
                    return
                }
                // Retry the order after token refresh
                order, err = client.CreateOrder(params)
            case 429:
                fmt.Println("Rate limit hit. Implementing exponential backoff...")
                // Implement backoff logic here
            default:
                fmt.Printf("API error: %v\n", apiErr)
            }
        } else {
            fmt.Printf("Network or client error: %v\n", err)
        }
        return
    }
    
    fmt.Printf("Order placed successfully! ID: %d\n", order.Id)
}
```

This document provides examples for all available methods in the Bitpin Go SDK. Each example includes proper error handling and demonstrates the recommended usage patterns. For more detailed information about the API responses and additional parameters, please refer to the [official Bitpin API documentation](https://docs.bitpin.ir/).