# Aster SDK Examples

This directory contains examples of using the go-aster SDK.

## Prerequisites

Before running these examples, make sure you have:
1. Go installed (version 1.18 or higher)
2. Valid Aster API credentials (API Key and Secret Key)

## Examples

### Basic Usage
- `basic_usage.go` - Simple example showing how to get account balance

### Spot Trading
- `spot_example.go` - Comprehensive spot trading examples including:
  - Market data queries
  - Account information
  - Order placement and management
  - Trade history

### Futures Trading
- `futures_example.go` - Complete futures trading examples including:
  - Account information
  - Market data
  - Order management
  - Position management
  
- `futures_change_leverage.go` - Detailed examples of changing leverage:
  - Basic leverage change
  - Error handling
  - Batch leverage changes
  - Debug mode

- `futures_leverage_simple.go` - Simple leverage change example (Chinese)

### WebSocket
- `websocket_example.go` - Real-time data streaming examples:
  - Spot market data streams
  - Futures market data streams
  - Order book updates
  - Trade updates
  - Position updates

## Running the Examples

1. Clone the repository:
```bash
git clone https://github.com/drinkthere/go-aster.git
cd go-aster/examples
```

2. Update the API credentials in the example files:
```go
apiKey := "your-api-key"
secretKey := "your-secret-key"
```

3. Run an example:
```bash
go run futures_change_leverage.go
```

## Important Notes

1. **API Credentials**: Never commit your actual API credentials to version control. Consider using environment variables:
```go
apiKey := os.Getenv("ASTER_API_KEY")
secretKey := os.Getenv("ASTER_SECRET_KEY")
```

2. **Rate Limits**: Be aware of Aster's rate limits when running these examples, especially those that make multiple API calls.

3. **Test Environment**: Consider using Aster's testnet first to familiarize yourself with the API before using real funds.

4. **Error Handling**: Always implement proper error handling in production code. The examples show basic error handling patterns.

## Common Issues

1. **"API key is required" error**: Make sure you've set your API key correctly
2. **"Signature for this request is not valid" error**: Check that your secret key is correct
3. **Connection errors**: Ensure you have internet connectivity and Aster's API is accessible