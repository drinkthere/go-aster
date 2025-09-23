package main

import (
	"context"
	"fmt"
	"log"
	"os"

	aster "github.com/drinkthere/go-aster/v3"
	"github.com/drinkthere/go-aster/v3/futures"
)

func main() {
	// Get credentials from environment variables
	userAddress := os.Getenv("ASTER_USER_ADDRESS")
	signerAddress := os.Getenv("ASTER_SIGNER_ADDRESS")
	privateKey := os.Getenv("ASTER_PRIVATE_KEY")

	if userAddress == "" || signerAddress == "" || privateKey == "" {
		log.Fatal("Please set ASTER_USER_ADDRESS, ASTER_SIGNER_ADDRESS and ASTER_PRIVATE_KEY environment variables")
	}

	// Initialize client with authentication
	client := aster.NewClient(userAddress, signerAddress, privateKey)
	
	// Enable debug mode to see requests
	client.Debug = true

	// Example 1: Get account information
	fmt.Println("=== Account Information ===")
	account, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		log.Printf("Get account error: %v", err)
	} else {
		fmt.Printf("Can trade: %v\n", account.CanTrade)
		fmt.Printf("Total wallet balance: %s\n", account.TotalWalletBalance)
		fmt.Printf("Total margin balance: %s\n", account.TotalMarginBalance)
		fmt.Printf("Available balance: %s\n", account.AvailableBalance)
	}

	// Example 2: Get balances
	fmt.Println("\n=== Balances ===")
	balances, err := client.NewGetBalanceService().Do(context.Background())
	if err != nil {
		log.Printf("Get balance error: %v", err)
	} else {
		for _, balance := range balances {
			if balance.Balance != "0" {
				fmt.Printf("%s: %s (available: %s)\n", 
					balance.Asset, balance.Balance, balance.AvailableBalance)
			}
		}
	}

	// Example 3: Get positions
	fmt.Println("\n=== Positions ===")
	positions, err := client.NewGetPositionsService().Do(context.Background())
	if err != nil {
		log.Printf("Get positions error: %v", err)
	} else {
		for _, position := range positions {
			if position.PositionAmt != "0" {
				fmt.Printf("Symbol: %s, Side: %s, Amount: %s, PnL: %s\n",
					position.Symbol, position.PositionSide, 
					position.PositionAmt, position.UnrealizedProfit)
			}
		}
	}

	// Example 4: Place a limit order
	fmt.Println("\n=== Place Order ===")
	symbol := "BTCUSDT"
	order, err := client.NewCreateOrderService().
		Symbol(symbol).
		Side(futures.SideTypeBuy).
		Type(futures.OrderTypeLimit).
		TimeInForce(futures.TimeInForceTypeGTC).
		Quantity("0.001").
		Price("50000").
		Do(context.Background())
	if err != nil {
		log.Printf("Create order error: %v", err)
	} else {
		fmt.Printf("Order placed successfully!\n")
		fmt.Printf("Order ID: %d\n", order.OrderID)
		fmt.Printf("Client order ID: %s\n", order.ClientOrderID)
		fmt.Printf("Status: %s\n", order.Status)
	}

	// Example 5: Get open orders
	fmt.Println("\n=== Open Orders ===")
	openOrders, err := client.NewListOpenOrdersService().
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		log.Printf("List open orders error: %v", err)
	} else {
		fmt.Printf("Number of open orders: %d\n", len(openOrders))
		for _, order := range openOrders {
			fmt.Printf("  Order %d: %s %s %s @ %s, Status: %s\n",
				order.OrderID, order.Symbol, order.Side, 
				order.OrigQty, order.Price, order.Status)
		}
	}

	// Example 6: Cancel an order (if you have the order ID)
	// orderID := int64(123456) // Replace with actual order ID
	// canceledOrder, err := client.NewCancelOrderService().
	//     Symbol(symbol).
	//     OrderID(orderID).
	//     Do(context.Background())
	// if err != nil {
	//     log.Printf("Cancel order error: %v", err)
	// } else {
	//     fmt.Printf("Order canceled: %d\n", canceledOrder.OrderID)
	// }
}