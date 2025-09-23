package main

import (
	"fmt"
	"log"
	"time"

	aster "github.com/drinkthere/go-aster/v3"
)

func main() {
	wsURL := "wss://fstream.asterdex.com"
	symbol := "btcusdt"

	// Example 1: Subscribe to depth updates
	fmt.Println("=== Subscribing to Depth Updates ===")
	wsDepthHandler := func(event *aster.WsDepthEvent) {
		fmt.Printf("Depth update - Symbol: %s, Time: %d\n", event.Symbol, event.Time)
		fmt.Printf("  Bids: %d, Asks: %d\n", len(event.Bids), len(event.Asks))
		if len(event.Bids) > 0 {
			fmt.Printf("  Best bid: %s @ %s\n", event.Bids[0][1], event.Bids[0][0])
		}
		if len(event.Asks) > 0 {
			fmt.Printf("  Best ask: %s @ %s\n", event.Asks[0][1], event.Asks[0][0])
		}
	}

	errHandler := func(err error) {
		log.Printf("WebSocket error: %v", err)
	}

	doneC, stopC, err := aster.WsDepthServe(wsURL, symbol, wsDepthHandler, errHandler)
	if err != nil {
		log.Fatal(err)
	}

	// Run for 10 seconds then stop
	go func() {
		time.Sleep(10 * time.Second)
		close(stopC)
	}()

	<-doneC
	fmt.Println("Depth subscription stopped")

	// Example 2: Subscribe to aggregate trades
	fmt.Println("\n=== Subscribing to Aggregate Trades ===")
	wsAggTradeHandler := func(event *aster.WsAggTradeEvent) {
		side := "SELL"
		if event.IsBuyerMaker {
			side = "BUY"
		}
		fmt.Printf("Trade - Symbol: %s, Price: %s, Quantity: %s, Side: %s\n",
			event.Symbol, event.Price, event.Quantity, side)
	}

	doneC2, stopC2, err := aster.WsAggTradeServe(wsURL, symbol, wsAggTradeHandler, errHandler)
	if err != nil {
		log.Fatal(err)
	}

	// Run for 10 seconds then stop
	go func() {
		time.Sleep(10 * time.Second)
		close(stopC2)
	}()

	<-doneC2
	fmt.Println("Trade subscription stopped")

	// Example 3: Subscribe to klines
	fmt.Println("\n=== Subscribing to Klines ===")
	wsKlineHandler := func(event *aster.WsKlineEvent) {
		k := event.Kline
		fmt.Printf("Kline - Symbol: %s, Interval: %s, O: %s, H: %s, L: %s, C: %s, V: %s\n",
			k.Symbol, k.Interval, k.Open, k.High, k.Low, k.Close, k.Volume)
	}

	doneC3, stopC3, err := aster.WsKlineServe(wsURL, symbol, "1m", wsKlineHandler, errHandler)
	if err != nil {
		log.Fatal(err)
	}

	// Run for 30 seconds then stop
	go func() {
		time.Sleep(30 * time.Second)
		close(stopC3)
	}()

	<-doneC3
	fmt.Println("Kline subscription stopped")
}