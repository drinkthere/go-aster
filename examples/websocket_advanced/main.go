package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	aster "github.com/drinkthere/go-aster/v3"
	"github.com/drinkthere/go-aster/v3/futures"
)

// PriceTracker tracks best bid/ask prices from book ticker
type PriceTracker struct {
	mu     sync.RWMutex
	prices map[string]BookPrice
}

type BookPrice struct {
	BidPrice string
	BidQty   string
	AskPrice string
	AskQty   string
	UpdateTime int64
}

func NewPriceTracker() *PriceTracker {
	return &PriceTracker{
		prices: make(map[string]BookPrice),
	}
}

func (pt *PriceTracker) UpdatePrice(symbol string, bp BookPrice) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	pt.prices[symbol] = bp
}

func (pt *PriceTracker) GetPrice(symbol string) (BookPrice, bool) {
	pt.mu.RLock()
	defer pt.mu.RUnlock()
	price, ok := pt.prices[symbol]
	return price, ok
}

func main() {
	// Initialize components
	priceTracker := NewPriceTracker()
	ctx := context.Background()
	
	// Track multiple symbols
	symbols := []string{"BTCUSDT", "ETHUSDT", "BNBUSDT"}
	
	// Start book ticker monitoring
	fmt.Println("=== Starting Book Ticker Monitoring ===")
	bookTickerStopChannels := startBookTickerMonitoring(symbols, priceTracker)
	
	// Start user data monitoring if credentials are available
	userDataStopC := startUserDataMonitoring(ctx, priceTracker)
	
	// Monitor prices and display updates
	monitorPrices(priceTracker, symbols)
	
	// Cleanup
	fmt.Println("\n=== Shutting down ===")
	for _, stopC := range bookTickerStopChannels {
		close(stopC)
	}
	if userDataStopC != nil {
		close(userDataStopC)
	}
	
	time.Sleep(2 * time.Second)
	fmt.Println("All streams closed")
}

func startBookTickerMonitoring(symbols []string, priceTracker *PriceTracker) []chan struct{} {
	var stopChannels []chan struct{}
	
	errHandler := func(err error) {
		log.Printf("Book ticker error: %v", err)
	}
	
	// Subscribe to book ticker for each symbol
	for _, symbol := range symbols {
		wsBookTickerHandler := func(event *aster.WsBookTickerEvent) {
			priceTracker.UpdatePrice(event.Symbol, BookPrice{
				BidPrice: event.BestBidPrice,
				BidQty:   event.BestBidQty,
				AskPrice: event.BestAskPrice,
				AskQty:   event.BestAskQty,
				UpdateTime: event.EventTime,
			})
		}
		
		_, stopC, err := aster.WsBookTickerServe(symbol, wsBookTickerHandler, errHandler)
		if err != nil {
			log.Printf("Failed to start book ticker for %s: %v", symbol, err)
			continue
		}
		
		stopChannels = append(stopChannels, stopC)
		fmt.Printf("Started book ticker stream for %s\n", symbol)
	}
	
	return stopChannels
}

func startUserDataMonitoring(ctx context.Context, priceTracker *PriceTracker) chan struct{} {
	// Get credentials
	userAddress := os.Getenv("ASTER_USER_ADDRESS")
	signerAddress := os.Getenv("ASTER_SIGNER_ADDRESS")
	privateKey := os.Getenv("ASTER_PRIVATE_KEY")
	
	if userAddress == "" || signerAddress == "" || privateKey == "" {
		fmt.Println("User data monitoring skipped - credentials not set")
		return nil
	}
	
	fmt.Println("\n=== Starting User Data Monitoring ===")
	
	// Initialize client
	client := aster.NewClient(userAddress, signerAddress, privateKey)
	
	// Create listen key
	listenKey, err := client.NewStartUserStreamService().Do(ctx)
	if err != nil {
		log.Printf("Failed to create listen key: %v", err)
		return nil
	}
	fmt.Printf("Listen key created: %s\n", listenKey)
	
	// User data handler
	wsUserDataHandler := func(event *aster.WsUserDataEvent) {
		switch event.Event {
		case "ACCOUNT_UPDATE":
			handleAccountUpdate(event.AccountUpdate)
		case "ORDER_TRADE_UPDATE":
			handleOrderUpdate(event.OrderUpdate, priceTracker)
		case "ACCOUNT_CONFIG_UPDATE":
			handleConfigUpdate(event.AccountConfigUpdate)
		}
	}
	
	errHandler := func(err error) {
		log.Printf("User data error: %v", err)
	}
	
	// Start user data stream
	_, stopC, err := aster.WsUserDataServe(listenKey, wsUserDataHandler, errHandler)
	if err != nil {
		log.Printf("Failed to start user data stream: %v", err)
		return nil
	}
	
	// Keep alive routine
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := client.NewKeepaliveUserStreamService().ListenKey(listenKey).Do(ctx)
				if err != nil {
					log.Printf("Failed to keepalive: %v", err)
					return
				}
			case <-stopC:
				client.NewCloseUserStreamService().ListenKey(listenKey).Do(ctx)
				return
			}
		}
	}()
	
	return stopC
}

func handleAccountUpdate(update *aster.WsAccountUpdate) {
	if update == nil {
		return
	}
	
	fmt.Printf("\n[ACCOUNT] Reason: %s\n", update.Reason)
	for _, balance := range update.Balances {
		if balance.Balance != "0" {
			fmt.Printf("  Balance %s: %s\n", balance.Asset, balance.Balance)
		}
	}
	for _, position := range update.Positions {
		if position.PositionAmount != "0" {
			fmt.Printf("  Position %s: Amount=%s, PnL=%s\n",
				position.Symbol, position.PositionAmount, position.UnrealizedPnL)
		}
	}
}

func handleOrderUpdate(update *aster.WsOrderUpdate, priceTracker *PriceTracker) {
	if update == nil {
		return
	}
	
	fmt.Printf("\n[ORDER] %s %s %s: Status=%s, Exec=%s\n",
		update.Symbol, update.Side, update.Type,
		update.Status, update.ExecutionType)
	
	// Show current market price for reference
	if price, ok := priceTracker.GetPrice(update.Symbol); ok {
		fmt.Printf("  Market: Bid=%s, Ask=%s\n", price.BidPrice, price.AskPrice)
	}
	
	fmt.Printf("  Order: Price=%s, Qty=%s, Filled=%s\n",
		update.Price, update.OriginalQty, update.FilledAccumulatedQty)
	
	if update.RealizedProfit != "" && update.RealizedProfit != "0" {
		fmt.Printf("  Realized Profit: %s\n", update.RealizedProfit)
	}
}

func handleConfigUpdate(update *aster.WsAccountConfigUpdate) {
	if update == nil {
		return
	}
	
	fmt.Printf("\n[CONFIG] %s: Leverage=%d, Margin=%s\n",
		update.Symbol, update.Leverage, update.MarginType)
}

func monitorPrices(priceTracker *PriceTracker, symbols []string) {
	fmt.Println("\n=== Price Monitoring Started ===")
	fmt.Println("Press Ctrl+C to stop\n")
	
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for i := 0; i < 12; i++ { // Run for 60 seconds
		select {
		case <-ticker.C:
			fmt.Printf("\n--- Price Update %s ---\n", time.Now().Format("15:04:05"))
			for _, symbol := range symbols {
				if price, ok := priceTracker.GetPrice(symbol); ok {
					spread := calculateSpread(price.BidPrice, price.AskPrice)
					fmt.Printf("%s: Bid=%s (%s), Ask=%s (%s), Spread=%.4f%%\n",
						symbol, price.BidPrice, price.BidQty,
						price.AskPrice, price.AskQty, spread)
				}
			}
		}
	}
}

func calculateSpread(bidStr, askStr string) float64 {
	var bid, ask float64
	fmt.Sscanf(bidStr, "%f", &bid)
	fmt.Sscanf(askStr, "%f", &ask)
	if bid > 0 {
		return ((ask - bid) / bid) * 100
	}
	return 0
}

// Example trading strategy using real-time data
func exampleTradingLogic(client *aster.Client, symbol string, price BookPrice) {
	// This is just an example - DO NOT use in production without proper risk management
	
	// Example: Place a limit buy order slightly below best bid
	var bidPrice float64
	fmt.Sscanf(price.BidPrice, "%f", &bidPrice)
	
	myPrice := fmt.Sprintf("%.2f", bidPrice * 0.999) // 0.1% below best bid
	
	order, err := client.NewCreateOrderService().
		Symbol(symbol).
		Side(futures.SideTypeBuy).
		Type(futures.OrderTypeLimit).
		TimeInForce(futures.TimeInForceTypeGTC).
		Quantity("0.001").
		Price(myPrice).
		Do(context.Background())
		
	if err != nil {
		log.Printf("Order failed: %v", err)
		return
	}
	
	fmt.Printf("Order placed: ID=%d, Price=%s\n", order.OrderID, myPrice)
}