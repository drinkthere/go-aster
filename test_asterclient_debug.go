package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/drinkthere/go-aster/v2"
	"github.com/drinkthere/go-aster/v2/common"
	"github.com/drinkthere/go-aster/v2/futures"
	"github.com/juju/ratelimit"
)

// 模拟你的配置结构
type Config struct {
	AsterAPIKey    string
	AsterSecretKey string
	APILimit10s    int64
	APILimit60s    int64
	LimitProcess   int
}

// 模拟你的 AsterClient
type AsterClient struct {
	FuturesClient *aster.BaseClient
	bucket10s     *ratelimit.Bucket
	bucket60s     *ratelimit.Bucket
	limitProcess  int
}

func (cli *AsterClient) Init(cfg *Config) bool {
	// 调试信息
	fmt.Printf("初始化 AsterClient:\n")
	fmt.Printf("  API Key: %s...\n", cfg.AsterAPIKey[:20])
	fmt.Printf("  Secret Key: %s...\n", cfg.AsterSecretKey[:20])
	
	cli.FuturesClient = aster.NewFuturesClient(cfg.AsterAPIKey, cfg.AsterSecretKey)
	
	// 验证客户端是否正确初始化
	if cli.FuturesClient == nil {
		fmt.Println("错误: FuturesClient 为 nil")
		return false
	}
	
	fmt.Printf("  FuturesClient.APIKey: %s...\n", cli.FuturesClient.APIKey[:20])
	fmt.Printf("  FuturesClient.BaseURL: %s\n", cli.FuturesClient.BaseURL)
	
	cli.bucket10s = ratelimit.NewBucketWithQuantum(10*time.Second, cfg.APILimit10s, cfg.APILimit10s)
	cli.bucket60s = ratelimit.NewBucketWithQuantum(60*time.Second, cfg.APILimit60s, cfg.APILimit60s)
	cli.limitProcess = cfg.LimitProcess
	return true
}

func (cli *AsterClient) CheckLimit(n int64) bool {
	if cli.bucket10s.TakeAvailable(n) < n {
		fmt.Println("[AsterClient] reach to 10s limit")
		return false
	}

	if cli.bucket60s.TakeAvailable(n) < n {
		fmt.Println("[AsterClient] reach to 60s limit")
		return false
	}

	return true
}

func main() {
	// 模拟配置
	globalConfig := &Config{
		AsterAPIKey:    "1b0b0731bcf2d5d7ba00485dc6eaa114fa049f9875203eb58d4bd246f043e993",
		AsterSecretKey: "f01bac3bf14d96b1ba8c8ae6fa4fb3f731bed5e612e531db0fc44c19bede3c1c",
		APILimit10s:    100,
		APILimit60s:    1000,
		LimitProcess:   10,
	}

	// 创建和初始化客户端
	asterClient := &AsterClient{}
	if !asterClient.Init(globalConfig) {
		log.Fatal("Failed to initialize AsterClient")
	}

	// 测试1：直接使用 FuturesClient
	fmt.Println("\n=== 测试1: 直接使用 FuturesClient ===")
	if asterClient.FuturesClient == nil {
		fmt.Println("错误: FuturesClient 为 nil")
	} else {
		fmt.Printf("FuturesClient APIKey: %s...\n", asterClient.FuturesClient.APIKey[:20])
	}

	// 测试2：使用 ChangeLeverageService
	fmt.Println("\n=== 测试2: 使用 ChangeLeverageService ===")
	instID := "ETHUSDT"
	leverage := 20

	// 检查限流
	if !asterClient.CheckLimit(1) {
		fmt.Println("Rate limit exceeded")
		return
	}

	// 创建服务
	service := &futures.ChangeLeverageService{C: asterClient.FuturesClient}
	
	// 打印调试信息
	fmt.Printf("Service.C: %v\n", service.C)
	if service.C != nil {
		fmt.Printf("Service.C.APIKey: %s...\n", service.C.APIKey[:20])
	}

	// 执行请求
	result, err := service.
		Symbol(instID).
		Leverage(leverage).
		Do(context.Background())

	if err != nil {
		fmt.Printf("修改杠杆失败: %v\n", err)
		
		// 详细错误信息
		if apiErr, ok := err.(*common.APIError); ok {
			fmt.Printf("API错误码: %d\n", apiErr.Code)
			fmt.Printf("API错误信息: %s\n", apiErr.Message)
		}
	} else {
		fmt.Printf("成功修改杠杆!\n")
		fmt.Printf("交易对: %s\n", result.Symbol)
		fmt.Printf("新杠杆: %d\n", result.Leverage)
	}

	// 测试3：检查是否是配置问题
	fmt.Println("\n=== 测试3: 检查配置 ===")
	if globalConfig.AsterAPIKey == "" {
		fmt.Println("警告: AsterAPIKey 为空!")
	}
	if globalConfig.AsterSecretKey == "" {
		fmt.Println("警告: AsterSecretKey 为空!")
	}

	// 测试4：创建一个新的客户端测试
	fmt.Println("\n=== 测试4: 新客户端测试 ===")
	testClient := aster.NewFuturesClient(
		globalConfig.AsterAPIKey, 
		globalConfig.AsterSecretKey,
		aster.WithDebug(true))
	
	testService := &futures.ChangeLeverageService{C: testClient}
	_, err = testService.Symbol("BTCUSDT").Leverage(10).Do(context.Background())
	if err != nil {
		fmt.Printf("测试失败: %v\n", err)
	} else {
		fmt.Println("测试成功!")
	}
}