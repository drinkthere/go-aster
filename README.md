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

- Market Data endpoints
- Account and Trading endpoints
- WebSocket support
- Automatic signature generation
- Rate limit handling

## Documentation

See [examples](./examples) directory for more usage examples.