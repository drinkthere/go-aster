package main

import (
	"context"
	"fmt"
	"log"

	"github.com/drinkthere/go-aster/v2"
	"github.com/drinkthere/go-aster/v2/common"
	"github.com/drinkthere/go-aster/v2/futures"
)

func main() {
	// Create a new futures client with Web3 signature
	userAddress := "0xYourUserAddress"
	signerAddress := "0xYourSignerAddress"
	privateKey := "yourPrivateKeyWithout0x"
	
	client := aster.NewFuturesClient(userAddress, signerAddress, privateKey)
	
	// Example 1: Get server time
	fmt.Println("=== Server Time ===")
	pingService := &futures.PingService{C: client}
	err := pingService.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Ping successful")
	
	// Example 2: Get exchange info
	fmt.Println("\n=== Exchange Info ===")
	exchangeService := &futures.ExchangeInfoService{C: client}
	exchangeInfo, err := exchangeService.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Exchange timezone: %s\n", exchangeInfo.Timezone)
	fmt.Printf("Number of symbols: %d\n", len(exchangeInfo.Symbols))
	
	// Example 3: Get market depth
	fmt.Println("\n=== Market Depth ===")
	depthService := &futures.DepthService{C: client}
	depthService.Symbol("BTCUSDT").Limit(5)
	depth, err := depthService.Do(context.Background())
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
	
	// Example 4: Get mark price
	fmt.Println("\n=== Mark Price ===")
	markPriceService := &futures.MarkPriceService{C: client}
	markPriceService.Symbol("BTCUSDT")
	markPrices, err := markPriceService.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if len(markPrices) > 0 {
		mp := markPrices[0]
		fmt.Printf("Symbol: %s\n", mp.Symbol)
		fmt.Printf("Mark Price: %s\n", mp.MarkPrice)
		fmt.Printf("Index Price: %s\n", mp.IndexPrice)
		fmt.Printf("Funding Rate: %s\n", mp.FundingRate)
		fmt.Printf("Next Funding Time: %d\n", mp.NextFundingTime)
	}
	
	// Example 5: Get commission rate
	fmt.Println("\n=== Commission Rate ===")
	commissionService := &futures.CommissionRateService{C: client}
	commissionService.Symbol("BTCUSDT")
	commission, err := commissionService.Do(context.Background())
	if err != nil {
		fmt.Printf("Error getting commission rate (expected if using demo keys): %v\n", err)
	} else {
		fmt.Printf("Symbol: %s\n", commission.Symbol)
		fmt.Printf("Maker Commission Rate: %s\n", commission.MakerCommissionRate)
		fmt.Printf("Taker Commission Rate: %s\n", commission.TakerCommissionRate)
	}
	
	// Example 6: Get account info (requires authentication)
	fmt.Println("\n=== Account Info ===")
	accountService := &futures.GetAccountService{C: client}
	account, err := accountService.Do(context.Background())
	if err != nil {
		fmt.Printf("Error getting account (expected if using demo keys): %v\n", err)
		return
	}
	
	fmt.Printf("Total Initial Margin: %s\n", account.TotalInitialMargin)
	fmt.Printf("Total Maint Margin: %s\n", account.TotalMaintMargin)
	fmt.Printf("Total Wallet Balance: %s\n", account.TotalWalletBalance)
	
	// Example: Create an order (commented out to prevent accidental orders)
	/*
	orderService := &futures.CreateOrderService{C: client}
	order, err := orderService.
		Symbol("BTCUSDT").
		Side(common.SideTypeBuy).
		Type(common.OrderTypeLimit).
		TimeInForce(common.TimeInForceTypeGTC).
		Quantity("0.001").
		Price("40000").
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Order created: %+v\n", order)
	*/
}