# Go-Aster SDK

A Go SDK for interacting with the Aster exchange API, supporting both Spot and Futures trading.

## Features

- Complete REST API coverage for Spot and Futures trading
- WebSocket support for real-time data streaming
- HMAC authentication for Spot trading
- Web3/Ethereum signature authentication for Futures trading
- Type-safe request/response structures
- Comprehensive error handling

## Installation

```bash
go get github.com/drinkthere/go-aster/v2
```

## Quick Start

### Spot Trading

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/drinkthere/go-aster/v2"
)

func main() {
    // Create a spot client
    client := aster.NewSpot("your-api-key", "your-secret-key")
    
    // Get server time
    serverTime, err := client.NewServerTimeService().Do(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Server time: %d\n", serverTime)
    
    // Get market depth
    depth, err := client.NewDepthService().Symbol("BTCUSDT").Limit(5).Do(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Bids: %v\n", depth.Bids)
    fmt.Printf("Asks: %v\n", depth.Asks)
}
```

### Futures Trading

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/drinkthere/go-aster/v2"
    "github.com/drinkthere/go-aster/v2/futures"
)

func main() {
    // Create a futures client with Web3 authentication
    client := aster.NewFuturesClient("userAddress", "signerAddress", "privateKey")
    
    // Get exchange info
    exchangeService := &futures.ExchangeInfoService{C: client}
    info, err := exchangeService.Do(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Number of symbols: %d\n", len(info.Symbols))
}
```

### WebSocket Streaming

```go
// Subscribe to book ticker updates
doneC, stopC, err := aster.WsSpotBookTickerServe("BTCUSDT", 
    func(event *aster.WsBookTickerEvent) {
        fmt.Printf("Best bid: %s @ %s, Best ask: %s @ %s\n",
            event.BestBidPrice, event.BestBidQty,
            event.BestAskPrice, event.BestAskQty)
    },
    func(err error) {
        log.Printf("Error: %v", err)
    })
```

## API Coverage

### Spot Trading
- Market Data: Depth, Trades, Klines, Tickers
- Trading: Create/Cancel/Query Orders
- Account: Balances, Trade History
- User Streams: Real-time account and order updates

### Futures Trading
- Market Data: Depth, Mark Price, Funding Rate, Klines
- Trading: Create/Cancel/Query Orders
- Account: Positions, Leverage, Margin Type
- User Streams: Real-time position and order updates

### WebSocket Streams
- Book Ticker
- Trade Streams
- Kline/Candlestick Streams
- Mark Price Streams (Futures)
- User Data Streams

## Authentication

### Spot (HMAC)
Spot trading uses HMAC-SHA256 signature:
```go
client := aster.NewSpot("api-key", "secret-key")
```

### Futures (Web3)
Futures trading uses Ethereum-style signatures:
```go
client := aster.NewFuturesClient("userAddress", "signerAddress", "privateKey")
```

## Examples

See the `examples/` directory for more comprehensive examples:
- `spot_example.go` - Spot trading examples
- `futures_example.go` - Futures trading examples
- `websocket_example.go` - WebSocket streaming examples

## Error Handling

The SDK provides typed errors for API responses:

```go
account, err := client.NewGetAccountService().Do(ctx)
if err != nil {
    if apiErr, ok := err.(*common.APIError); ok {
        fmt.Printf("API Error: Code=%d, Message=%s\n", apiErr.Code, apiErr.Message)
    } else {
        fmt.Printf("Network Error: %v\n", err)
    }
}
```

## Configuration

### Custom HTTP Client
```go
httpClient := &http.Client{
    Timeout: 30 * time.Second,
}
client := aster.NewSpot("key", "secret", aster.WithHTTPClient(httpClient))
```

### Custom Base URL
```go
client := aster.NewSpot("key", "secret", 
    aster.WithBaseURL("https://testnet-sapi.asterdex.com"))
```

### Debug Mode
```go
client := aster.NewSpot("key", "secret", aster.WithDebug(true))
```

## License

This project is licensed under the MIT License.