package main

import (
	"context"
	"fmt"
	"log"

	"github.com/drinkthere/go-aster/v2"
	"github.com/drinkthere/go-aster/v2/futures"
)

func main() {
	// 你的 API 凭证
	apiKey := "1b0b0731bcf2d5d7ba00485dc6eaa114fa049f9875203eb58d4bd246f043e993"
	secretKey := "f01bac3bf14d96b1ba8c8ae6fa4fb3f731bed5e612e531db0fc44c19bede3c1c"
	
	fmt.Printf("API Key: %s...\n", apiKey[:20])
	fmt.Printf("Secret Key: %s...\n", secretKey[:20])
	
	// 测试1: 使用期货客户端
	fmt.Println("\n=== 测试期货客户端 ===")
	futuresClient := aster.NewFuturesClient(apiKey, secretKey, 
		aster.WithDebug(true),
		aster.WithBaseURL("https://fapi.asterdex.com")) // 确保使用正确的期货 URL
	
	// 确认客户端设置
	fmt.Printf("客户端 API Key: %s\n", futuresClient.APIKey)
	fmt.Printf("客户端 Base URL: %s\n", futuresClient.BaseURL)
	fmt.Printf("签名类型: %v\n", futuresClient.SignatureType)
	
	// 测试简单的 API 调用
	fmt.Println("\n测试 Ping...")
	pingService := &futures.PingService{C: futuresClient}
	if err := pingService.Do(context.Background()); err != nil {
		log.Printf("Ping 失败: %v", err)
	} else {
		fmt.Println("Ping 成功!")
	}
	
	// 测试需要认证的 API
	fmt.Println("\n测试获取账户...")
	accountService := &futures.GetAccountService{C: futuresClient}
	_, err := accountService.Do(context.Background())
	if err != nil {
		log.Printf("获取账户失败: %v", err)
	}
	
	// 测试2: 使用现货客户端对比
	fmt.Println("\n=== 测试现货客户端 ===")
	spotClient := aster.NewSpot(apiKey, secretKey,
		aster.WithDebug(true))
	
	fmt.Printf("现货客户端 API Key: %s\n", spotClient.APIKey)
	fmt.Printf("现货客户端 Base URL: %s\n", spotClient.BaseURL)
	
	// 测试现货 API
	fmt.Println("\n测试现货服务器时间...")
	serverTime, err := spotClient.NewServerTimeService().Do(context.Background())
	if err != nil {
		log.Printf("现货获取时间失败: %v", err)
	} else {
		fmt.Printf("现货服务器时间: %d\n", serverTime)
	}
	
	// 测试现货账户
	fmt.Println("\n测试现货账户...")
	spotAccount, err := spotClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		log.Printf("现货获取账户失败: %v", err)
	} else {
		fmt.Printf("现货账户更新时间: %d\n", spotAccount.UpdateTime)
	}
}