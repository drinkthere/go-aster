package main

import (
	"context"
	"fmt"
	"time"

	"github.com/drinkthere/go-aster/v2"
	"github.com/drinkthere/go-aster/v2/futures"
	"github.com/juju/ratelimit"
)

// 模拟你的代码结构
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
	cli.FuturesClient = aster.NewFuturesClient(cfg.AsterAPIKey, cfg.AsterSecretKey)
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
	// 测试场景1：正常使用（你的方式）
	fmt.Println("=== 测试场景1: 你的代码结构 ===")
	
	globalConfig := Config{
		AsterAPIKey:    "1b0b0731bcf2d5d7ba00485dc6eaa114fa049f9875203eb58d4bd246f043e993",
		AsterSecretKey: "f01bac3bf14d96b1ba8c8ae6fa4fb3f731bed5e612e531db0fc44c19bede3c1c",
		APILimit10s:    100,
		APILimit60s:    1000,
		LimitProcess:   10,
	}
	
	// 打印配置
	fmt.Printf("Config APIKey: %s...\n", globalConfig.AsterAPIKey[:20])
	fmt.Printf("Config SecretKey: %s...\n", globalConfig.AsterSecretKey[:20])
	
	// 你的方式 - 创建 asterClient
	asterClient := AsterClient{}  // 注意：不是指针
	asterClient.Init(&globalConfig)
	
	// 检查客户端状态
	fmt.Printf("\nasterClient.FuturesClient != nil: %v\n", asterClient.FuturesClient != nil)
	if asterClient.FuturesClient != nil {
		fmt.Printf("asterClient.FuturesClient.APIKey: %s...\n", asterClient.FuturesClient.APIKey[:20])
	}
	
	// 使用服务
	instID := "ETHUSDT"
	instCfg := struct{ Leverage int }{Leverage: 20}
	
	service := &futures.ChangeLeverageService{C: asterClient.FuturesClient}
	_, err := service.
		Symbol(instID).
		Leverage(instCfg.Leverage).
		Do(context.Background())
		
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Println("成功!")
	}
	
	// 测试场景2：使用指针
	fmt.Println("\n=== 测试场景2: 使用指针 ===")
	
	asterClientPtr := &AsterClient{}  // 使用指针
	asterClientPtr.Init(&globalConfig)
	
	service2 := &futures.ChangeLeverageService{C: asterClientPtr.FuturesClient}
	_, err = service2.
		Symbol(instID).
		Leverage(instCfg.Leverage).
		Do(context.Background())
		
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Println("成功!")
	}
	
	// 测试场景3：可能的问题 - 配置传递
	fmt.Println("\n=== 测试场景3: 检查配置传递问题 ===")
	
	// 模拟可能的问题：配置对象的值在传递过程中丢失
	testConfig := &Config{}  // 空配置
	
	fmt.Printf("空配置 APIKey: '%s'\n", testConfig.AsterAPIKey)
	fmt.Printf("空配置 SecretKey: '%s'\n", testConfig.AsterSecretKey)
	
	asterClient3 := &AsterClient{}
	asterClient3.Init(testConfig)
	
	service3 := &futures.ChangeLeverageService{C: asterClient3.FuturesClient}
	_, err = service3.
		Symbol("ETHUSDT").
		Leverage(20).
		Do(context.Background())
		
	if err != nil {
		fmt.Printf("预期的错误: %v\n", err)
	}
	
	// 测试场景4：检查全局变量问题
	fmt.Println("\n=== 测试场景4: 全局变量问题 ===")
	
	// 模拟全局变量
	var globalAsterClient AsterClient
	
	// 忘记初始化
	if globalAsterClient.FuturesClient == nil {
		fmt.Println("全局客户端未初始化!")
	}
	
	// 正确初始化
	globalAsterClient.Init(&globalConfig)
	
	service4 := &futures.ChangeLeverageService{C: globalAsterClient.FuturesClient}
	_, err = service4.
		Symbol("ETHUSDT").
		Leverage(20).
		Do(context.Background())
		
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Println("成功!")
	}
}