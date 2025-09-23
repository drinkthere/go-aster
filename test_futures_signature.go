package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/drinkthere/go-aster/v2"
	"github.com/drinkthere/go-aster/v2/common"
	"github.com/drinkthere/go-aster/v2/futures"
)

func main() {
	// 你的 API 凭证
	apiKey := "cc7ac0adbd279bf10db8206d8e0d8e71ce33aa6a301f64cafc257889c35361ac"
	secretKey := "cc7ac0adbd279bf10db8206d8e0d8e71ce33aa6a301f64cafc257889c35361ac"
	
	// 创建带调试模式的客户端
	client := aster.NewFuturesClient(apiKey, secretKey,
		aster.WithDebug(true), // 开启调试模式查看请求详情
		aster.WithHTTPClient(&http.Client{
			Timeout: 10 * time.Second,
		}))
	
	// 测试1: 先测试不需要签名的接口
	fmt.Println("=== 测试 Ping (不需要签名) ===")
	pingService := &futures.PingService{C: client}
	if err := pingService.Do(context.Background()); err != nil {
		log.Printf("Ping 失败: %v", err)
	} else {
		fmt.Println("Ping 成功!")
	}
	
	// 测试2: 测试获取服务器时间
	fmt.Println("\n=== 测试获取服务器时间 ===")
	timeService := &futures.ServerTimeService{C: client}
	serverTime, err := timeService.Do(context.Background())
	if err != nil {
		log.Printf("获取服务器时间失败: %v", err)
	} else {
		fmt.Printf("服务器时间: %d\n", serverTime)
		// 计算时间偏差
		localTime := time.Now().UnixNano() / 1e6
		timeDiff := serverTime - localTime
		fmt.Printf("时间偏差: %d ms\n", timeDiff)
	}
	
	// 测试3: 测试需要签名的接口（获取账户信息）
	fmt.Println("\n=== 测试获取账户信息 (需要签名) ===")
	accountService := &futures.GetAccountService{C: client}
	account, err := accountService.Do(context.Background())
	if err != nil {
		log.Printf("获取账户信息失败: %v", err)
		// 打印具体错误
		if apiErr, ok := err.(*common.APIError); ok {
			log.Printf("API 错误码: %d, 错误信息: %s", apiErr.Code, apiErr.Message)
		}
	} else {
		fmt.Printf("账户信息获取成功!\n")
		fmt.Printf("总余额: %s\n", account.TotalWalletBalance)
	}
	
	// 测试4: 测试取消所有订单
	fmt.Println("\n=== 测试取消所有订单 ===")
	cancelService := &futures.CancelAllOpenOrdersService{C: client}
	cancelService.Symbol("ETHUSDT")
	if err := cancelService.Do(context.Background()); err != nil {
		log.Printf("取消订单失败: %v", err)
		// 打印具体错误
		if apiErr, ok := err.(*common.APIError); ok {
			log.Printf("API 错误码: %d, 错误信息: %s", apiErr.Code, apiErr.Message)
		}
	} else {
		fmt.Println("取消订单成功!")
	}
}