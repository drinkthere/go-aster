# go-aster

A Go SDK for Aster Finance Futures API v3.

## Installation

```bash
go get -u github.com/drinkthere/go-aster
```

## Quick Start

```go
import (
    "github.com/drinkthere/go-aster/v3"
)

// Create a new client
client := aster.NewClient(userAddress, signerAddress, privateKey)

// Get server time
serverTime, err := client.NewServerTimeService().Do(context.Background())
```

## Features

- Market Data endpoints (depth, trades, klines, ticker)
- Account and Trading endpoints (orders, positions, balances)
- WebSocket support
  - Market data streams (depth, trades, klines, book ticker)
  - User data streams (order updates, position updates, balance changes)
  - Combined streams for multiple symbols
- Automatic signature generation with Web3
- Rate limit handling
- Reconnection support for WebSocket

## Documentation

See [examples](./examples) directory for more usage examples.