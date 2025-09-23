package main

import (
	"context"
	"fmt"
	"time"

	"github.com/drinkthere/go-aster/v2"
	"github.com/drinkthere/go-aster/v2/futures"
	"github.com/juju/ratelimit"
)

type Config struct {
	AsterAPIKey    string
	AsterSecretKey string
	APILimit10s    int64
	APILimit60s    int64
	LimitProcess   int
}

type AsterClient struct {
	FuturesClient *aster.BaseClient
	bucket10s     *ratelimit.Bucket
	bucket60s     *ratelimit.Bucket
	limitProcess  int
}

func (cli *AsterClient) Init(cfg *Config) bool {
	// 调试输出
	fmt.Printf("[Init] 接收到的配置:\n")
	fmt.Printf("  APIKey: '%s'\n", cfg.AsterAPIKey)
	fmt.Printf("  SecretKey: '%s'\n", cfg.AsterSecretKey)
	
	// 创建客户端
	cli.FuturesClient = aster.NewFuturesClient(cfg.AsterAPIKey, cfg.AsterSecretKey)
	
	// 验证客户端创建后的状态
	fmt.Printf("[Init] 客户端创建后:\n")
	fmt.Printf("  FuturesClient != nil: %v\n", cli.FuturesClient != nil)
	if cli.FuturesClient != nil {
		fmt.Printf("  FuturesClient.APIKey: '%s'\n", cli.FuturesClient.APIKey)
		fmt.Printf("  FuturesClient.SecretKey是否为空: %v\n", cli.FuturesClient.SecretKey == "")
		fmt.Printf("  FuturesClient.BaseURL: %s\n", cli.FuturesClient.BaseURL)
	}
	
	// 初始化限流器（使用默认值避免panic）
	limit10s := cfg.APILimit10s
	if limit10s <= 0 {
		limit10s = 100
	}
	limit60s := cfg.APILimit60s
	if limit60s <= 0 {
		limit60s = 1000
	}
	
	cli.bucket10s = ratelimit.NewBucketWithQuantum(10*time.Second, limit10s, limit10s)
	cli.bucket60s = ratelimit.NewBucketWithQuantum(60*time.Second, limit60s, limit60s)
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

// 模拟你的使用场景
func simulateYourScenario() {
	fmt.Println("=== 模拟你的场景 ===\n")
	
	// 1. 创建配置
	globalConfig := Config{
		AsterAPIKey:    "1b0b0731bcf2d5d7ba00485dc6eaa114fa049f9875203eb58d4bd246f043e993",
		AsterSecretKey: "f01bac3bf14d96b1ba8c8ae6fa4fb3f731bed5e612e531db0fc44c19bede3c1c",
		APILimit10s:    100,
		APILimit60s:    1000,
		LimitProcess:   10,
	}
	
	fmt.Printf("globalConfig.AsterAPIKey: %s...\n\n", globalConfig.AsterAPIKey[:20])
	
	// 2. 创建 AsterClient（你的方式）
	asterClient := AsterClient{}
	asterClient.Init(&globalConfig)
	
	// 3. 检查 Init 之后的状态
	fmt.Printf("\n[Main] Init后的状态检查:\n")
	fmt.Printf("  asterClient.FuturesClient != nil: %v\n", asterClient.FuturesClient != nil)
	if asterClient.FuturesClient != nil {
		fmt.Printf("  asterClient.FuturesClient.APIKey: '%s'\n", asterClient.FuturesClient.APIKey)
		if asterClient.FuturesClient.APIKey == "" {
			fmt.Println("  警告: APIKey 为空!")
		}
	}
	
	// 4. 使用服务
	fmt.Printf("\n[Main] 调用 ChangeLeverageService...\n")
	instID := "ETHUSDT"
	instCfg := struct{ Leverage int }{Leverage: 20}
	
	service := &futures.ChangeLeverageService{C: asterClient.FuturesClient}
	
	// 再次检查
	fmt.Printf("  service.C != nil: %v\n", service.C != nil)
	if service.C != nil {
		fmt.Printf("  service.C.APIKey: '%s'\n", service.C.APIKey)
	}
	
	_, err := service.
		Symbol(instID).
		Leverage(instCfg.Leverage).
		Do(context.Background())
		
	if err != nil {
		fmt.Printf("\n[错误] %v\n", err)
		
		// 如果是 API key required 错误，进行详细诊断
		if err.Error() == "API key is required" {
			fmt.Println("\n诊断信息:")
			fmt.Printf("1. asterClient.FuturesClient: %v\n", asterClient.FuturesClient)
			if asterClient.FuturesClient != nil {
				fmt.Printf("2. APIKey 长度: %d\n", len(asterClient.FuturesClient.APIKey))
				fmt.Printf("3. APIKey 内容: '%s'\n", asterClient.FuturesClient.APIKey)
			}
		}
	} else {
		fmt.Println("\n[成功] 修改杠杆成功!")
	}
}

// 可能的问题场景
func testPossibleIssues() {
	fmt.Println("\n\n=== 测试可能的问题 ===\n")
	
	// 场景1：配置对象被修改
	fmt.Println("场景1: 配置对象被意外修改")
	config1 := &Config{
		AsterAPIKey:    "test-key",
		AsterSecretKey: "test-secret",
	}
	
	client1 := &AsterClient{}
	
	// 模拟配置被清空
	config1.AsterAPIKey = ""
	config1.AsterSecretKey = ""
	
	client1.Init(config1)
	
	// 场景2：使用全局变量但忘记初始化
	fmt.Println("\n场景2: 全局变量问题")
	var globalClient AsterClient
	// 忘记调用 Init
	
	if globalClient.FuturesClient == nil {
		fmt.Println("全局客户端未初始化!")
	}
	
	// 场景3：并发问题
	fmt.Println("\n场景3: 检查是否有并发问题")
	// 如果你在多个 goroutine 中使用同一个 client，可能会有问题
}

func main() {
	// 运行主要测试
	simulateYourScenario()
	
	// 测试可能的问题
	testPossibleIssues()
}