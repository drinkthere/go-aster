package main

import (
	"context"
	"fmt"
	"log"

	"github.com/drinkthere/go-aster/v2"
	"github.com/drinkthere/go-aster/v2/futures"
)

func main() {
	// API credentials
	apiKey := "your-api-key"
	secretKey := "your-secret-key"

	// Create futures client
	client := aster.NewFuturesClient(apiKey, secretKey)

	// Example 1: Change leverage for a specific symbol
	fmt.Println("=== Example 1: Change Leverage for ETHUSDT ===")
	leverageService := &futures.ChangeLeverageService{C: client}
	result, err := leverageService.
		Symbol("ETHUSDT").
		Leverage(10).
		Do(context.Background())
	
	if err != nil {
		log.Printf("Error changing leverage: %v", err)
		return
	}

	fmt.Printf("Successfully changed leverage for %s\n", result.Symbol)
	fmt.Printf("New leverage: %d\n", result.Leverage)
	fmt.Printf("Max notional value: %s\n", result.MaxNotionalValue)

	// Example 2: Change leverage with error handling
	fmt.Println("\n=== Example 2: Change Leverage with Error Handling ===")
	leverageService2 := &futures.ChangeLeverageService{C: client}
	result2, err := leverageService2.
		Symbol("BTCUSDT").
		Leverage(25).
		Do(context.Background())
	
	if err != nil {
		// Handle specific API errors
		if apiErr, ok := err.(*aster.APIError); ok {
			fmt.Printf("API Error Code: %d\n", apiErr.Code)
			fmt.Printf("API Error Message: %s\n", apiErr.Message)
			
			// Common error codes
			switch apiErr.Code {
			case -4028:
				fmt.Println("Leverage value is not valid")
			case -1102:
				fmt.Println("Mandatory parameter was not sent")
			case -1121:
				fmt.Println("Invalid symbol")
			default:
				fmt.Printf("Other API error: %s\n", apiErr.Message)
			}
		} else {
			log.Printf("Network or other error: %v", err)
		}
		return
	}

	fmt.Printf("Successfully changed leverage for %s to %d\n", result2.Symbol, result2.Leverage)

	// Example 3: Query current leverage settings
	fmt.Println("\n=== Example 3: Query Current Leverage Settings ===")
	// First, get account information to see position details
	accountService := &futures.GetAccountService{C: client}
	account, err := accountService.Do(context.Background())
	if err != nil {
		log.Printf("Error getting account: %v", err)
		return
	}

	// Show leverage for active positions
	fmt.Println("Current leverage settings for positions:")
	for _, pos := range account.Positions {
		if pos.PositionAmt != "0" || pos.Symbol == "ETHUSDT" || pos.Symbol == "BTCUSDT" {
			fmt.Printf("  %s: Leverage = %s, Position = %s\n", 
				pos.Symbol, pos.Leverage, pos.PositionAmt)
		}
	}

	// Example 4: Batch change leverage for multiple symbols
	fmt.Println("\n=== Example 4: Batch Change Leverage ===")
	symbols := []string{"ETHUSDT", "BTCUSDT", "BNBUSDT"}
	newLeverage := 15

	for _, symbol := range symbols {
		service := &futures.ChangeLeverageService{C: client}
		result, err := service.
			Symbol(symbol).
			Leverage(newLeverage).
			Do(context.Background())
		
		if err != nil {
			fmt.Printf("  Failed to change leverage for %s: %v\n", symbol, err)
			continue
		}
		
		fmt.Printf("  %s: Changed to %dx (max notional: %s)\n", 
			symbol, result.Leverage, result.MaxNotionalValue)
	}

	// Example 5: With debug mode to see request details
	fmt.Println("\n=== Example 5: Change Leverage with Debug Mode ===")
	debugClient := aster.NewFuturesClient(apiKey, secretKey, 
		aster.WithDebug(true))
	
	debugService := &futures.ChangeLeverageService{C: debugClient}
	_, err = debugService.
		Symbol("SOLUSDT").
		Leverage(20).
		Do(context.Background())
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Successfully changed leverage for SOLUSDT")
	}
}