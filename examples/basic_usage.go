package main

import (
	"context"
	"fmt"
	"log"

	"github.com/drinkthere/go-aster/v2"
	"github.com/drinkthere/go-aster/v2/futures"
)

func main() {
	// API 凭证
	apiKey := "your-api-key"
	secretKey := "your-secret-key"

	// 创建现货客户端
	spotClient := aster.NewSpot(apiKey, secretKey)
	
	// 创建期货客户端（同样使用 API Key + Secret Key）
	futuresClient := aster.NewFuturesClient(apiKey, secretKey)

	// ========== 现货示例 ==========
	fmt.Println("=== 现货交易示例 ===")
	
	// 1. 获取服务器时间
	serverTime, err := spotClient.NewServerTimeService().Do(context.Background())
	if err != nil {
		log.Printf("获取服务器时间失败: %v", err)
	} else {
		fmt.Printf("服务器时间: %d\n", serverTime)
	}

	// 2. 获取账户信息
	account, err := spotClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		log.Printf("获取现货账户失败: %v", err)
	} else {
		fmt.Printf("账户更新时间: %d\n", account.UpdateTime)
		// 显示余额不为零的资产
		for _, balance := range account.Balances {
			if balance.Free != "0" || balance.Locked != "0" {
				fmt.Printf("  %s - 可用: %s, 冻结: %s\n", balance.Asset, balance.Free, balance.Locked)
			}
		}
	}

	// 3. 获取交易对价格
	prices, err := spotClient.NewListPricesService().Symbol("BTCUSDT").Do(context.Background())
	if err != nil {
		log.Printf("获取价格失败: %v", err)
	} else if len(prices) > 0 {
		fmt.Printf("BTCUSDT 价格: %s\n", prices[0].Price)
	}

	// ========== 期货示例 ==========
	fmt.Println("\n=== 期货交易示例 ===")
	
	// 使用 BaseClient 调用期货 API
	// 注意：期货的服务需要通过 futures 包来创建
	
	// 1. Ping 测试
	fmt.Println("测试期货连接...")
	pingService := &futures.PingService{C: futuresClient}
	if err := pingService.Do(context.Background()); err != nil {
		log.Printf("期货连接失败: %v", err)
	} else {
		fmt.Println("期货连接成功!")
	}
	
	// 2. 获取交易规则信息
	infoService := &futures.ExchangeInfoService{C: futuresClient}
	info, err := infoService.Do(context.Background())
	if err != nil {
		log.Printf("获取交易规则失败: %v", err)
	} else {
		fmt.Printf("期货交易对数量: %d\n", len(info.Symbols))
	}
	
	// 3. 获取标记价格
	markService := &futures.MarkPriceService{C: futuresClient}
	markService.Symbol("BTCUSDT")
	marks, err := markService.Do(context.Background())
	if err != nil {
		log.Printf("获取标记价格失败: %v", err)
	} else if len(marks) > 0 {
		fmt.Printf("BTCUSDT 标记价格: %s\n", marks[0].MarkPrice)
		fmt.Printf("BTCUSDT 指数价格: %s\n", marks[0].IndexPrice)
	}
}