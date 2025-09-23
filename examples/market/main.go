package main

import (
	"context"
	"fmt"
	"log"

	aster "github.com/your-org/go-aster/v3"
	"github.com/your-org/go-aster/v3/futures"
)

func main() {
	// Initialize client (no auth needed for market data)
	client := aster.NewClient("", "", "")

	// Example 1: Ping server
	fmt.Println("=== Ping Server ===")
	err := client.NewPingService().Do(context.Background())
	if err != nil {
		log.Printf("Ping error: %v", err)
	} else {
		fmt.Println("Ping successful")
	}

	// Example 2: Get server time
	fmt.Println("\n=== Server Time ===")
	serverTime, err := client.NewServerTimeService().Do(context.Background())
	if err != nil {
		log.Printf("Server time error: %v", err)
	} else {
		fmt.Printf("Server time: %d\n", serverTime)
	}

	// Example 3: Get exchange info
	fmt.Println("\n=== Exchange Info ===")
	exchangeInfo, err := client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		log.Printf("Exchange info error: %v", err)
	} else {
		fmt.Printf("Timezone: %s\n", exchangeInfo.Timezone)
		fmt.Printf("Number of symbols: %d\n", len(exchangeInfo.Symbols))
		if len(exchangeInfo.Symbols) > 0 {
			fmt.Printf("First symbol: %s\n", exchangeInfo.Symbols[0].Symbol)
		}
	}

	// Example 4: Get order book depth
	fmt.Println("\n=== Order Book Depth ===")
	symbol := "BTCUSDT"
	depth, err := client.NewDepthService().
		Symbol(symbol).
		Limit(10).
		Do(context.Background())
	if err != nil {
		log.Printf("Depth error: %v", err)
	} else {
		fmt.Printf("Order book for %s:\n", symbol)
		fmt.Printf("Bids (top 3):\n")
		for i := 0; i < 3 && i < len(depth.Bids); i++ {
			fmt.Printf("  Price: %s, Quantity: %s\n", depth.Bids[i].Price, depth.Bids[i].Quantity)
		}
		fmt.Printf("Asks (top 3):\n")
		for i := 0; i < 3 && i < len(depth.Asks); i++ {
			fmt.Printf("  Price: %s, Quantity: %s\n", depth.Asks[i].Price, depth.Asks[i].Quantity)
		}
	}

	// Example 5: Get recent trades
	fmt.Println("\n=== Recent Trades ===")
	trades, err := client.NewRecentTradesListService().
		Symbol(symbol).
		Limit(5).
		Do(context.Background())
	if err != nil {
		log.Printf("Trades error: %v", err)
	} else {
		fmt.Printf("Recent trades for %s:\n", symbol)
		for _, trade := range trades {
			fmt.Printf("  Price: %s, Quantity: %s, Time: %d\n", 
				trade.Price, trade.Quantity, trade.Time)
		}
	}

	// Example 6: Get klines
	fmt.Println("\n=== Klines ===")
	klines, err := client.NewKlinesService().
		Symbol(symbol).
		Interval(futures.KlineInterval1h).
		Limit(5).
		Do(context.Background())
	if err != nil {
		log.Printf("Klines error: %v", err)
	} else {
		fmt.Printf("1h klines for %s:\n", symbol)
		for _, kline := range klines {
			fmt.Printf("  Open: %s, High: %s, Low: %s, Close: %s, Volume: %s\n",
				kline.Open, kline.High, kline.Low, kline.Close, kline.Volume)
		}
	}

	// Example 7: Get ticker price
	fmt.Println("\n=== Ticker Price ===")
	tickers, err := client.NewTickerPriceService().
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		log.Printf("Ticker price error: %v", err)
	} else {
		for _, ticker := range tickers {
			fmt.Printf("%s price: %s\n", ticker.Symbol, ticker.Price)
		}
	}
}