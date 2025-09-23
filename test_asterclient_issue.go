package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/drinkthere/go-aster/v2"
	"github.com/drinkthere/go-aster/v2/futures"
)

// 测试场景1：配置值为空
func testEmptyConfig() {
	fmt.Println("\n=== 测试场景1: 配置值为空 ===")
	
	// 模拟空的配置
	apiKey := ""
	secretKey := ""
	
	client := aster.NewFuturesClient(apiKey, secretKey)
	
	fmt.Printf("Client APIKey: '%s'\n", client.APIKey)
	fmt.Printf("Client SecretKey: '%s'\n", client.SecretKey)
	
	service := &futures.ChangeLeverageService{C: client}
	_, err := service.
		Symbol("ETHUSDT").
		Leverage(20).
		Do(context.Background())
		
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	}
}

// 测试场景2：客户端未正确初始化
func testNilClient() {
	fmt.Println("\n=== 测试场景2: 客户端为nil ===")
	
	var client *aster.BaseClient
	
	if client == nil {
		fmt.Println("客户端是 nil，跳过测试以避免 panic")
		return
	}
	
	service := &futures.ChangeLeverageService{C: client}
	_, err := service.
		Symbol("ETHUSDT").
		Leverage(20).
		Do(context.Background())
		
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	}
}

// 测试场景3：正确的配置但可能有空格
func testConfigWithSpaces() {
	fmt.Println("\n=== 测试场景3: 配置值有空格 ===")
	
	// 模拟可能有空格的配置
	apiKey := " 1b0b0731bcf2d5d7ba00485dc6eaa114fa049f9875203eb58d4bd246f043e993 "
	secretKey := " f01bac3bf14d96b1ba8c8ae6fa4fb3f731bed5e612e531db0fc44c19bede3c1c "
	
	fmt.Printf("原始 APIKey: '%s' (长度: %d)\n", apiKey, len(apiKey))
	
	client := aster.NewFuturesClient(apiKey, secretKey)
	
	fmt.Printf("Client APIKey: '%s' (长度: %d)\n", client.APIKey, len(client.APIKey))
	
	service := &futures.ChangeLeverageService{C: client}
	_, err := service.
		Symbol("ETHUSDT").
		Leverage(20).
		Do(context.Background())
		
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Println("成功!")
	}
}

// 测试场景4：检查是否是指针问题
type AsterClientWrapper struct {
	Client *aster.BaseClient
}

func testPointerIssue() {
	fmt.Println("\n=== 测试场景4: 指针问题 ===")
	
	wrapper := &AsterClientWrapper{}
	
	// 忘记初始化 Client
	if wrapper.Client == nil {
		fmt.Println("wrapper.Client 是 nil!")
	}
	
	// 正确初始化
	wrapper.Client = aster.NewFuturesClient(
		"1b0b0731bcf2d5d7ba00485dc6eaa114fa049f9875203eb58d4bd246f043e993",
		"f01bac3bf14d96b1ba8c8ae6fa4fb3f731bed5e612e531db0fc44c19bede3c1c")
	
	service := &futures.ChangeLeverageService{C: wrapper.Client}
	_, err := service.
		Symbol("ETHUSDT").
		Leverage(20).
		Do(context.Background())
		
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Println("成功!")
	}
}

// 添加调试函数来检查你的实际场景
func debugYourScenario(apiKey, secretKey string) {
	fmt.Println("\n=== 调试你的实际场景 ===")
	
	// 打印配置值（隐藏敏感信息）
	if apiKey == "" {
		fmt.Println("警告: API Key 为空!")
	} else {
		fmt.Printf("API Key: %s... (长度: %d)\n", apiKey[:10], len(apiKey))
	}
	
	if secretKey == "" {
		fmt.Println("警告: Secret Key 为空!")
	} else {
		fmt.Printf("Secret Key: %s... (长度: %d)\n", secretKey[:10], len(secretKey))
	}
	
	// 创建客户端
	client := aster.NewFuturesClient(apiKey, secretKey)
	
	// 验证客户端状态
	fmt.Printf("\n客户端状态:\n")
	fmt.Printf("  Client != nil: %v\n", client != nil)
	if client != nil {
		fmt.Printf("  Client.APIKey: '%s'\n", client.APIKey)
		fmt.Printf("  Client.SecretKey 是否为空: %v\n", client.SecretKey == "")
		fmt.Printf("  Client.BaseURL: %s\n", client.BaseURL)
	}
	
	// 尝试调用
	fmt.Printf("\n尝试调用 ChangeLeverageService...\n")
	service := &futures.ChangeLeverageService{C: client}
	
	// 打印 service 状态
	fmt.Printf("Service.C != nil: %v\n", service.C != nil)
	
	_, err := service.
		Symbol("ETHUSDT").
		Leverage(20).
		Do(context.Background())
		
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Println("成功!")
	}
}

func main() {
	// 运行所有测试场景
	testEmptyConfig()
	testNilClient()
	testConfigWithSpaces()
	testPointerIssue()
	
	// 测试你的实际场景（使用你的实际配置值）
	// 请替换为你的实际配置获取方式
	// debugYourScenario(cfg.AsterAPIKey, cfg.AsterSecretKey)
	
	// 或者使用测试值
	fmt.Println("\n" + strings.Repeat("=", 50))
	debugYourScenario(
		"1b0b0731bcf2d5d7ba00485dc6eaa114fa049f9875203eb58d4bd246f043e993",
		"f01bac3bf14d96b1ba8c8ae6fa4fb3f731bed5e612e531db0fc44c19bede3c1c")
}