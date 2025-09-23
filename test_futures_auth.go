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
	apiKey := "1b0b0731bcf2d5d7ba00485dc6eaa114fa049f9875203eb58d4bd246f043e993"
	secretKey := "f01bac3bf14d96b1ba8c8ae6fa4fb3f731bed5e612e531db0fc44c19bede3c1c"
	
	// 创建带调试模式的客户端
	client := aster.NewFuturesClient(apiKey, secretKey,
		aster.WithDebug(true), // 开启调试模式查看请求详情
		aster.WithHTTPClient(&http.Client{
			Timeout: 10 * time.Second,
		}))
	
	fmt.Println("=== 测试 Aster 期货 API ===")
	fmt.Println("使用 HMAC 签名方式")
	fmt.Printf("API Key: %s...\n", apiKey[:20])
	
	// 测试1: 先测试不需要签名的接口
	fmt.Println("\n1. 测试 Ping (不需要签名)")
	pingService := &futures.PingService{C: client}
	if err := pingService.Do(context.Background()); err != nil {
		log.Printf("Ping 失败: %v", err)
	} else {
		fmt.Println("✓ Ping 成功!")
	}
	
	// 测试2: 测试获取服务器时间
	fmt.Println("\n2. 测试获取服务器时间")
	timeService := &futures.ServerTimeService{C: client}
	serverTime, err := timeService.Do(context.Background())
	if err != nil {
		log.Printf("获取服务器时间失败: %v", err)
	} else {
		fmt.Printf("✓ 服务器时间: %d\n", serverTime)
		// 计算时间偏差
		localTime := time.Now().UnixNano() / 1e6
		timeDiff := serverTime - localTime
		fmt.Printf("  时间偏差: %d ms\n", timeDiff)
		
		// 如果时间偏差过大，可能导致签名失败
		if abs(timeDiff) > 5000 {
			fmt.Println("  警告：时间偏差过大，可能导致签名验证失败!")
		}
	}
	
	// 测试3: 测试交易规则（不需要签名）
	fmt.Println("\n3. 测试获取交易规则")
	exchangeService := &futures.ExchangeInfoService{C: client}
	exchangeInfo, err := exchangeService.Do(context.Background())
	if err != nil {
		log.Printf("获取交易规则失败: %v", err)
	} else {
		fmt.Printf("✓ 获取成功，交易对数量: %d\n", len(exchangeInfo.Symbols))
		// 查找 ETHUSDT
		for _, symbol := range exchangeInfo.Symbols {
			if symbol.Symbol == "ETHUSDT" {
				fmt.Printf("  ETHUSDT 状态: %s\n", symbol.Status)
				break
			}
		}
	}
	
	// 测试4: 测试需要签名的接口（获取账户信息）
	fmt.Println("\n4. 测试获取账户信息 (需要签名)")
	accountService := &futures.GetAccountService{C: client}
	account, err := accountService.Do(context.Background())
	if err != nil {
		log.Printf("获取账户信息失败: %v", err)
		if apiErr, ok := err.(*common.APIError); ok {
			log.Printf("API 错误码: %d, 错误信息: %s", apiErr.Code, apiErr.Message)
			
			// 分析可能的错误原因
			switch apiErr.Code {
			case -1022:
				fmt.Println("  可能原因: 签名无效，请检查 Secret Key")
			case -2015:
				fmt.Println("  可能原因: API Key 无效或无权限")
			case -1021:
				fmt.Println("  可能原因: 时间戳超出接收窗口")
			}
		}
	} else {
		fmt.Println("✓ 账户信息获取成功!")
		fmt.Printf("  总余额: %s\n", account.TotalWalletBalance)
		fmt.Printf("  可用余额: %s\n", account.AvailableBalance)
	}
	
	// 测试5: 测试取消所有订单
	fmt.Println("\n5. 测试取消 ETHUSDT 所有订单")
	cancelService := &futures.CancelAllOpenOrdersService{C: client}
	cancelService.Symbol("ETHUSDT")
	if err := cancelService.Do(context.Background()); err != nil {
		log.Printf("取消订单失败: %v", err)
		if apiErr, ok := err.(*common.APIError); ok {
			log.Printf("API 错误码: %d, 错误信息: %s", apiErr.Code, apiErr.Message)
			if apiErr.Code == -1022 {
				fmt.Println("  错误: 签名验证失败")
				fmt.Println("  请检查:")
				fmt.Println("  1. API Key 和 Secret Key 是否正确")
				fmt.Println("  2. 是否使用了正确的认证方式")
				fmt.Println("  3. 系统时间是否准确")
			}
		}
	} else {
		fmt.Println("✓ 取消订单成功!")
	}
	
	// 测试6: 获取持仓信息
	fmt.Println("\n6. 测试获取持仓信息")
	positionService := &futures.GetPositionRiskService{C: client}
	positions, err := positionService.Do(context.Background())
	if err != nil {
		log.Printf("获取持仓信息失败: %v", err)
	} else {
		fmt.Printf("✓ 获取成功，持仓数量: %d\n", len(positions))
		for _, pos := range positions {
			if pos.PositionAmt != "0" {
				fmt.Printf("  %s: %s\n", pos.Symbol, pos.PositionAmt)
			}
		}
	}
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}