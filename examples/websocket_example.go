package main

import (
	"fmt"
	"log"
	"time"

	"github.com/drinkthere/go-aster/v2"
	"github.com/drinkthere/go-aster/v2/common"
)

func main() {
	// Example 1: Spot Book Ticker WebSocket
	fmt.Println("=== Spot Book Ticker WebSocket ===")
	spotClient := aster.NewSpot("", "") // No auth needed for public streams
	
	doneC, stopC, err := spotClient.WsBookTickerServe("BTCUSDT", func(event *aster.WsBookTickerEvent) {
		fmt.Printf("Book Ticker - Symbol: %s, Best Bid: %s @ %s, Best Ask: %s @ %s\n",
			event.Symbol, event.BestBidPrice, event.BestBidQty, event.BestAskPrice, event.BestAskQty)
	}, func(err error) {
		log.Printf("Error: %v", err)
	})
	
	if err != nil {
		log.Fatal(err)
	}
	
	// Let it run for 10 seconds
	go func() {
		time.Sleep(10 * time.Second)
		close(stopC)
	}()
	
	// Wait for connection to close
	<-doneC
	fmt.Println("Book ticker stream closed")
	
	// Example 2: Futures Kline WebSocket
	fmt.Println("\n=== Futures Kline WebSocket ===")
	doneC2, stopC2, err := aster.WsFuturesKlineServe("BTCUSDT", string(common.Interval1m), 
		func(event *aster.WsFuturesKlineEvent) {
			k := event.Kline
			fmt.Printf("Kline - Symbol: %s, Time: %d, Open: %s, High: %s, Low: %s, Close: %s, Volume: %s, Final: %v\n",
				k.Symbol, k.StartTime, k.Open, k.High, k.Low, k.Close, k.Volume, k.IsFinal)
		}, 
		func(err error) {
			log.Printf("Error: %v", err)
		})
	
	if err != nil {
		log.Fatal(err)
	}
	
	// Let it run for 10 seconds
	go func() {
		time.Sleep(10 * time.Second)
		close(stopC2)
	}()
	
	// Wait for connection to close
	<-doneC2
	fmt.Println("Kline stream closed")
	
	// Example 3: User Data Stream (requires authentication)
	fmt.Println("\n=== User Data Stream Example ===")
	
	// For spot user data
	spotAuthClient := aster.NewSpot("your-api-key", "your-secret-key")
	
	// First, create a listen key
	listenKey, err := spotAuthClient.NewStartUserStreamService().Do(nil)
	if err != nil {
		fmt.Printf("Error creating listen key (expected if using demo keys): %v\n", err)
		return
	}
	
	fmt.Printf("Listen key created: %s\n", listenKey)
	
	// Connect to user data stream
	doneC3, stopC3, err := spotAuthClient.WsUserDataServe(listenKey, 
		func(event *aster.WsSpotUserDataEvent) {
			switch event.Event {
			case "outboundAccountPosition":
				fmt.Println("Account update received")
				if event.AccountUpdate != nil {
					for _, balance := range event.AccountUpdate.Balances {
						fmt.Printf("  Asset: %s, Free: %s, Locked: %s\n", 
							balance.Asset, balance.Free, balance.Locked)
					}
				}
			case "executionReport":
				fmt.Println("Order update received")
				if event.OrderUpdate != nil {
					fmt.Printf("  Symbol: %s, Side: %s, Type: %s, Status: %s, Price: %s, Qty: %s\n",
						event.OrderUpdate.Symbol, event.OrderUpdate.Side, event.OrderUpdate.OrderType,
						event.OrderUpdate.OrderStatus, event.OrderUpdate.OriginalPrice, event.OrderUpdate.OriginalQuantity)
				}
			}
		},
		func(err error) {
			log.Printf("User data error: %v", err)
		})
	
	if err != nil {
		log.Fatal(err)
	}
	
	// Keepalive for listen key
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := spotAuthClient.NewKeepaliveUserStreamService().ListenKey(listenKey).Do(nil)
				if err != nil {
					log.Printf("Error keeping listen key alive: %v", err)
				}
			case <-stopC3:
				return
			}
		}
	}()
	
	// Let it run for 30 seconds
	time.Sleep(30 * time.Second)
	close(stopC3)
	
	// Wait for connection to close
	<-doneC3
	fmt.Println("User data stream closed")
	
	// Close the listen key
	err = spotAuthClient.NewCloseUserStreamService().ListenKey(listenKey).Do(nil)
	if err != nil {
		log.Printf("Error closing listen key: %v", err)
	}
}