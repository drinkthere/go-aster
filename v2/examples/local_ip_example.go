package main

import (
	"fmt"
	"log"

	aster "github.com/drinkthere/go-aster/v2"
)

func main() {
	// Example 1: Create a spot client with local IP binding
	spotClient := aster.NewSpot("your-api-key", "your-secret-key", 
		aster.WithLocalAddress("192.168.1.100"), // Use specific local IP
		aster.WithDebug(true),
	)

	// Example 2: Create a futures client with local IP binding  
	futuresClient := aster.NewFutures("your-api-key", "your-secret-key",
		aster.WithLocalAddress("192.168.1.100"), // Use specific local IP
		aster.WithDebug(true),
	)

	// Example 3: WebSocket connections will also use the specified local IP
	fmt.Println("Starting WebSocket with local IP binding...")
	
	// Spot book ticker with local IP
	_, stopC, err := spotClient.WsBookTickerServe("BTCUSDT", func(event *aster.WsBookTickerEvent) {
		fmt.Printf("Spot Book Ticker: %s - Bid: %s, Ask: %s\n", 
			event.Symbol, event.BestBidPrice, event.BestAskPrice)
	}, func(err error) {
		log.Printf("Spot WebSocket error: %v", err)
	})
	
	if err != nil {
		log.Printf("Failed to start spot WebSocket: %v", err)
		return
	}

	// Futures book ticker with local IP
	_, stopC2, err := futuresClient.WsBookTickerServe("BTCUSDT", func(event *aster.WsBookTickerEvent) {
		fmt.Printf("Futures Book Ticker: %s - Bid: %s, Ask: %s\n", 
			event.Symbol, event.BestBidPrice, event.BestAskPrice)
	}, func(err error) {
		log.Printf("Futures WebSocket error: %v", err)
	})
	
	if err != nil {
		log.Printf("Failed to start futures WebSocket: %v", err)
		return
	}

	// Example 4: You can also specify different local IP for individual connections
	_, stopC3, err := futuresClient.WsBookTickerServeWithLocalAddr("ETHUSDT", func(event *aster.WsBookTickerEvent) {
		fmt.Printf("ETH Futures Book Ticker (custom IP): %s - Bid: %s, Ask: %s\n", 
			event.Symbol, event.BestBidPrice, event.BestAskPrice)
	}, func(err error) {
		log.Printf("ETH WebSocket error: %v", err)
	}, "192.168.1.101") // Different local IP for this connection
	
	if err != nil {
		log.Printf("Failed to start ETH WebSocket: %v", err)
		return
	}

	fmt.Println("WebSocket connections started. Press Enter to stop...")
	fmt.Scanln()

	// Clean shutdown
	close(stopC)
	close(stopC2) 
	close(stopC3)
	
	fmt.Println("All connections stopped.")
}