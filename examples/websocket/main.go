package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	aster "github.com/drinkthere/go-aster/v3"
)

func main() {
	// Example 1: Market Data WebSocket (no auth required)
	marketDataExample()

	// Example 2: User Data WebSocket (requires auth)
	userDataExample()
}

func marketDataExample() {
	fmt.Println("=== Market Data WebSocket Examples ===")
	
	symbol := "BTCUSDT"
	errHandler := func(err error) {
		log.Printf("WebSocket error: %v", err)
	}

	// Subscribe to depth updates
	fmt.Println("\n--- Depth Updates ---")
	wsDepthHandler := func(event *aster.WsDepthEvent) {
		fmt.Printf("Depth - Symbol: %s, Bids: %d, Asks: %d\n", 
			event.Symbol, len(event.Bids), len(event.Asks))
		if len(event.Bids) > 0 && len(event.Asks) > 0 {
			fmt.Printf("  Best bid: %s @ %s, Best ask: %s @ %s\n", 
				event.Bids[0][1], event.Bids[0][0], event.Asks[0][1], event.Asks[0][0])
		}
	}

	doneC1, stopC1, err := aster.WsDepthServe(symbol, wsDepthHandler, errHandler)
	if err != nil {
		log.Printf("Failed to start depth websocket: %v", err)
		return
	}

	// Subscribe to book ticker
	fmt.Println("\n--- Book Ticker Updates ---")
	wsBookTickerHandler := func(event *aster.WsBookTickerEvent) {
		fmt.Printf("BookTicker - Symbol: %s, Bid: %s (%s), Ask: %s (%s)\n",
			event.Symbol, event.BestBidPrice, event.BestBidQty,
			event.BestAskPrice, event.BestAskQty)
	}

	doneC2, stopC2, err := aster.WsBookTickerServe(symbol, wsBookTickerHandler, errHandler)
	if err != nil {
		log.Printf("Failed to start book ticker websocket: %v", err)
		close(stopC1)
		<-doneC1
		return
	}

	// Subscribe to aggregate trades
	fmt.Println("\n--- Aggregate Trades ---")
	wsAggTradeHandler := func(event *aster.WsAggTradeEvent) {
		side := "SELL"
		if event.IsBuyerMaker {
			side = "BUY"
		}
		fmt.Printf("Trade - Price: %s, Qty: %s, Side: %s, Time: %d\n",
			event.Price, event.Quantity, side, event.TradeTime)
	}

	doneC3, stopC3, err := aster.WsAggTradeServe(symbol, wsAggTradeHandler, errHandler)
	if err != nil {
		log.Printf("Failed to start trade websocket: %v", err)
		close(stopC1)
		close(stopC2)
		<-doneC1
		<-doneC2
		return
	}

	// Run for 15 seconds then stop all
	time.Sleep(15 * time.Second)
	fmt.Println("\n--- Stopping all market data streams ---")
	close(stopC1)
	close(stopC2)
	close(stopC3)
	<-doneC1
	<-doneC2
	<-doneC3
	fmt.Println("All market data streams stopped")
}

func userDataExample() {
	fmt.Println("\n\n=== User Data WebSocket Example ===")

	// Get credentials from environment variables
	userAddress := os.Getenv("ASTER_USER_ADDRESS")
	signerAddress := os.Getenv("ASTER_SIGNER_ADDRESS")
	privateKey := os.Getenv("ASTER_PRIVATE_KEY")

	if userAddress == "" || signerAddress == "" || privateKey == "" {
		fmt.Println("Skipping user data example - credentials not set")
		fmt.Println("Set ASTER_USER_ADDRESS, ASTER_SIGNER_ADDRESS and ASTER_PRIVATE_KEY to test user data streams")
		return
	}

	// Initialize client
	client := aster.NewClient(userAddress, signerAddress, privateKey)
	ctx := context.Background()

	// Create listen key
	listenKey, err := client.NewStartUserStreamService().Do(ctx)
	if err != nil {
		log.Printf("Failed to create listen key: %v", err)
		return
	}
	fmt.Printf("Listen key created: %s\n", listenKey)

	// Set up user data handler
	wsUserDataHandler := func(event *aster.WsUserDataEvent) {
		switch event.Event {
		case "ACCOUNT_UPDATE":
			if event.AccountUpdate != nil {
				fmt.Printf("\n=== Account Update ===\n")
				fmt.Printf("Reason: %s\n", event.AccountUpdate.Reason)
				for _, balance := range event.AccountUpdate.Balances {
					if balance.Balance != "0" {
						fmt.Printf("  %s: Wallet=%s, Cross=%s\n", 
							balance.Asset, balance.Balance, balance.CrossWalletBalance)
					}
				}
				for _, position := range event.AccountUpdate.Positions {
					if position.PositionAmount != "0" {
						fmt.Printf("  Position %s: Amount=%s, Entry=%s, PnL=%s\n",
							position.Symbol, position.PositionAmount, 
							position.EntryPrice, position.UnrealizedPnL)
					}
				}
			}
		case "ORDER_TRADE_UPDATE":
			if event.OrderUpdate != nil {
				fmt.Printf("\n=== Order Update ===\n")
				fmt.Printf("Symbol: %s, Side: %s, Type: %s\n", 
					event.OrderUpdate.Symbol, event.OrderUpdate.Side, event.OrderUpdate.Type)
				fmt.Printf("Status: %s, Exec Type: %s\n", 
					event.OrderUpdate.Status, event.OrderUpdate.ExecutionType)
				fmt.Printf("Price: %s, Qty: %s, Filled: %s\n", 
					event.OrderUpdate.Price, event.OrderUpdate.OriginalQty, 
					event.OrderUpdate.FilledAccumulatedQty)
				if event.OrderUpdate.RealizedProfit != "" {
					fmt.Printf("Realized Profit: %s\n", event.OrderUpdate.RealizedProfit)
				}
			}
		case "ACCOUNT_CONFIG_UPDATE":
			if event.AccountConfigUpdate != nil {
				fmt.Printf("\n=== Account Config Update ===\n")
				fmt.Printf("Symbol: %s, Leverage: %d, Margin Type: %s\n",
					event.AccountConfigUpdate.Symbol, 
					event.AccountConfigUpdate.Leverage, 
					event.AccountConfigUpdate.MarginType)
			}
		default:
			fmt.Printf("\n=== Unknown Event: %s ===\n", event.Event)
		}
	}

	errHandler := func(err error) {
		log.Printf("User data WebSocket error: %v", err)
	}

	// Start user data stream
	doneC, stopC, err := aster.WsUserDataServe(listenKey, wsUserDataHandler, errHandler)
	if err != nil {
		log.Printf("Failed to start user data websocket: %v", err)
		return
	}

	fmt.Println("User data stream started - listening for account updates...")

	// Keep alive routine
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := client.NewKeepaliveUserStreamService().ListenKey(listenKey).Do(ctx)
				if err != nil {
					log.Printf("Failed to keepalive user stream: %v", err)
					close(stopC)
					return
				}
				fmt.Println("Listen key renewed")
			case <-stopC:
				return
			}
		}
	}()

	// Run for 60 seconds then stop
	time.Sleep(60 * time.Second)
	fmt.Println("\n--- Stopping user data stream ---")
	close(stopC)
	<-doneC

	// Close listen key
	err = client.NewCloseUserStreamService().ListenKey(listenKey).Do(ctx)
	if err != nil {
		log.Printf("Failed to close listen key: %v", err)
	}
	fmt.Println("User data stream stopped and listen key closed")
}