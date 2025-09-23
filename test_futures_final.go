package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
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
	
	fmt.Println("=== Aster 期货 API 测试 ===")
	
	// 1. 手动测试 API 调用
	fmt.Println("\n1. 手动测试 API 调用")
	testManualRequest(apiKey, secretKey)
	
	// 2. 使用 SDK 测试
	fmt.Println("\n2. 使用 SDK 测试")
	testWithSDK(apiKey, secretKey)
}

// 手动构造请求测试
func testManualRequest(apiKey, secretKey string) {
	// 测试获取账户信息
	baseURL := "https://fapi.asterdex.com"
	endpoint := "/fapi/v2/account"
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	
	// 构造参数
	params := fmt.Sprintf("timestamp=%s", timestamp)
	
	// 计算签名
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(params))
	signature := hex.EncodeToString(mac.Sum(nil))
	
	// 构造完整 URL
	fullURL := fmt.Sprintf("%s%s?%s&signature=%s", baseURL, endpoint, params, signature)
	fmt.Printf("请求 URL: %s\n", fullURL)
	
	// 创建请求
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Printf("创建请求失败: %v", err)
		return
	}
	
	// 设置头部
	req.Header.Set("X-MBX-APIKEY", apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	fmt.Printf("请求头 X-MBX-APIKEY: %s...\n", apiKey[:20])
	
	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("发送请求失败: %v", err)
		return
	}
	defer resp.Body.Close()
	
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应失败: %v", err)
		return
	}
	
	fmt.Printf("响应状态: %d\n", resp.StatusCode)
	fmt.Printf("响应内容: %s\n", string(body))
}

// 使用 SDK 测试
func testWithSDK(apiKey, secretKey string) {
	// 创建客户端
	client := aster.NewFuturesClient(apiKey, secretKey,
		aster.WithDebug(true))
	
	// 测试不需要签名的接口
	fmt.Println("\n测试 Ping...")
	pingService := &futures.PingService{C: client}
	if err := pingService.Do(context.Background()); err != nil {
		log.Printf("Ping 失败: %v", err)
	} else {
		fmt.Println("✓ Ping 成功")
	}
	
	// 测试需要签名的接口
	fmt.Println("\n测试获取账户信息...")
	accountService := &futures.GetAccountService{C: client}
	account, err := accountService.Do(context.Background())
	if err != nil {
		log.Printf("获取账户失败: %v", err)
		if apiErr, ok := err.(*common.APIError); ok {
			fmt.Printf("错误码: %d, 错误信息: %s\n", apiErr.Code, apiErr.Message)
		}
	} else {
		fmt.Printf("✓ 获取成功，总余额: %s\n", account.TotalWalletBalance)
	}
	
	// 测试取消订单
	fmt.Println("\n测试取消 ETHUSDT 订单...")
	cancelService := &futures.CancelAllOpenOrdersService{C: client}
	cancelService.Symbol("ETHUSDT")
	if err := cancelService.Do(context.Background()); err != nil {
		log.Printf("取消订单失败: %v", err)
		if apiErr, ok := err.(*common.APIError); ok {
			fmt.Printf("错误码: %d, 错误信息: %s\n", apiErr.Code, apiErr.Message)
		}
	} else {
		fmt.Println("✓ 取消订单成功")
	}
	
	// 测试修改杠杆
	fmt.Println("\n测试修改 ETHUSDT 杠杆...")
	leverageService := &futures.ChangeLeverageService{C: client}
	result, err := leverageService.Symbol("ETHUSDT").Leverage(20).Do(context.Background())
	if err != nil {
		log.Printf("修改杠杆失败: %v", err)
		if apiErr, ok := err.(*common.APIError); ok {
			fmt.Printf("错误码: %d, 错误信息: %s\n", apiErr.Code, apiErr.Message)
		}
	} else {
		fmt.Printf("✓ 修改杠杆成功，新杠杆: %d\n", result.Leverage)
	}
	
	// 额外测试：直接构造一个简单的请求
	fmt.Println("\n直接测试 BaseClient...")
	r := aster.NewRequest(http.MethodGet, "/fapi/v1/time", aster.SecTypeNone)
	data, err := client.CallAPI(context.Background(), r)
	if err != nil {
		log.Printf("直接调用失败: %v", err)
	} else {
		fmt.Printf("✓ 直接调用成功: %s\n", string(data))
	}
}