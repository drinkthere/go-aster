package main

import (
	"context"
	"fmt"
	"log"

	"github.com/drinkthere/go-aster/v2"
)

func main() {
	// Create a new spot client
	apiKey := "your-api-key"
	secretKey := "your-secret-key"
	
	client := aster.NewSpot(apiKey, secretKey)
	
	// Example 1: Get server time
	fmt.Println("=== Server Time ===")
	serverTime, err := client.NewServerTimeService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Server time: %d\n", serverTime)
	
	// Example 2: Get exchange info
	fmt.Println("\n=== Exchange Info ===")
	exchangeInfo, err := client.NewExchangeInfoService().Symbol("BTCUSDT").Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Exchange timezone: %s\n", exchangeInfo.Timezone)
	if len(exchangeInfo.Symbols) > 0 {
		symbol := exchangeInfo.Symbols[0]
		fmt.Printf("Symbol: %s, Status: %s\n", symbol.Symbol, symbol.Status)
	}
	
	// Example 3: Get depth
	fmt.Println("\n=== Market Depth ===")
	depth, err := client.NewDepthService().Symbol("BTCUSDT").Limit(5).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Top 5 Bids:")
	for i, bid := range depth.Bids {
		if i >= 5 {
			break
		}
		fmt.Printf("  Price: %s, Quantity: %s\n", bid[0], bid[1])
	}
	fmt.Println("Top 5 Asks:")
	for i, ask := range depth.Asks {
		if i >= 5 {
			break
		}
		fmt.Printf("  Price: %s, Quantity: %s\n", ask[0], ask[1])
	}
	
	// Example 4: Get klines
	fmt.Println("\n=== Klines ===")
	klines, err := client.NewKlinesService().
		Symbol("BTCUSDT").
		Interval("1h").
		Limit(3).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for _, kline := range klines {
		fmt.Printf("Time: %d, Open: %s, High: %s, Low: %s, Close: %s, Volume: %s\n",
			kline.OpenTime, kline.Open, kline.High, kline.Low, kline.Close, kline.Volume)
	}
	
	// Example 5: Get 24hr ticker
	fmt.Println("\n=== 24hr Ticker ===")
	tickers, err := client.NewListPriceChangeStatsService().Symbol("BTCUSDT").Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if len(tickers) > 0 {
		ticker := tickers[0]
		fmt.Printf("Symbol: %s\n", ticker.Symbol)
		fmt.Printf("Price Change: %s (%s%%)\n", ticker.PriceChange, ticker.PriceChangePercent)
		fmt.Printf("Last Price: %s\n", ticker.LastPrice)
		fmt.Printf("Volume: %s\n", ticker.Volume)
	}
}