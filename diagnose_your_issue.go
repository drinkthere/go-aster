package main

import (
	"context"
	"fmt"
	"reflect"

	"github.com/drinkthere/go-aster/v2"
	"github.com/drinkthere/go-aster/v2/futures"
)

// 添加这个诊断函数到你的代码中
func DiagnoseAsterClient(asterClient interface{}, cfg interface{}) {
	fmt.Println("\n=== Aster Client 诊断 ===")
	
	// 1. 检查配置
	fmt.Println("\n1. 配置检查:")
	cfgValue := reflect.ValueOf(cfg)
	if cfgValue.Kind() == reflect.Ptr {
		cfgValue = cfgValue.Elem()
	}
	
	apiKeyField := cfgValue.FieldByName("AsterAPIKey")
	secretKeyField := cfgValue.FieldByName("AsterSecretKey")
	
	if apiKeyField.IsValid() {
		apiKey := apiKeyField.String()
		fmt.Printf("   Config.AsterAPIKey: '%s' (长度: %d)\n", apiKey, len(apiKey))
		if apiKey == "" {
			fmt.Println("   ⚠️  警告: API Key 为空!")
		}
	} else {
		fmt.Println("   ❌ 错误: 找不到 AsterAPIKey 字段!")
	}
	
	if secretKeyField.IsValid() {
		secretKey := secretKeyField.String()
		fmt.Printf("   Config.AsterSecretKey: '%s' (长度: %d)\n", maskString(secretKey), len(secretKey))
		if secretKey == "" {
			fmt.Println("   ⚠️  警告: Secret Key 为空!")
		}
	} else {
		fmt.Println("   ❌ 错误: 找不到 AsterSecretKey 字段!")
	}
	
	// 2. 检查客户端
	fmt.Println("\n2. 客户端检查:")
	clientValue := reflect.ValueOf(asterClient)
	if clientValue.Kind() == reflect.Ptr {
		clientValue = clientValue.Elem()
	}
	
	futuresClientField := clientValue.FieldByName("FuturesClient")
	if !futuresClientField.IsValid() {
		fmt.Println("   ❌ 错误: 找不到 FuturesClient 字段!")
		return
	}
	
	if futuresClientField.IsNil() {
		fmt.Println("   ❌ 错误: FuturesClient 是 nil!")
		fmt.Println("   建议: 确保调用了 Init 方法")
		return
	}
	
	// 获取 BaseClient
	baseClient := futuresClientField.Interface().(*aster.BaseClient)
	fmt.Printf("   FuturesClient != nil: %v\n", baseClient != nil)
	fmt.Printf("   FuturesClient.APIKey: '%s' (长度: %d)\n", baseClient.APIKey, len(baseClient.APIKey))
	fmt.Printf("   FuturesClient.SecretKey 是否为空: %v\n", baseClient.SecretKey == "")
	fmt.Printf("   FuturesClient.BaseURL: %s\n", baseClient.BaseURL)
	
	// 3. 测试调用
	fmt.Println("\n3. 测试 API 调用:")
	service := &futures.ChangeLeverageService{C: baseClient}
	_, err := service.
		Symbol("ETHUSDT").
		Leverage(20).
		Do(context.Background())
		
	if err != nil {
		fmt.Printf("   ❌ API 调用失败: %v\n", err)
		
		// 特定错误分析
		if err.Error() == "API key is required" {
			fmt.Println("\n   错误分析: 'API key is required'")
			fmt.Println("   可能原因:")
			fmt.Println("   1. APIKey 在某处被清空了")
			fmt.Println("   2. 使用了错误的 client 实例")
			fmt.Println("   3. client 被重新初始化了")
			
			// 再次检查当前状态
			fmt.Printf("\n   当前 APIKey: '%s'\n", baseClient.APIKey)
			fmt.Printf("   APIKey == \"\": %v\n", baseClient.APIKey == "")
		}
	} else {
		fmt.Println("   ✅ API 调用成功!")
	}
}

func maskString(s string) string {
	if len(s) <= 10 {
		return "***"
	}
	return s[:5] + "..." + s[len(s)-5:]
}

// 使用示例
func main() {
	// 模拟你的结构
	type Config struct {
		AsterAPIKey    string
		AsterSecretKey string
	}
	
	type AsterClient struct {
		FuturesClient *aster.BaseClient
	}
	
	// 场景1: 正常情况
	fmt.Println("=== 场景1: 正常情况 ===")
	config1 := &Config{
		AsterAPIKey:    "1b0b0731bcf2d5d7ba00485dc6eaa114fa049f9875203eb58d4bd246f043e993",
		AsterSecretKey: "f01bac3bf14d96b1ba8c8ae6fa4fb3f731bed5e612e531db0fc44c19bede3c1c",
	}
	
	client1 := &AsterClient{}
	client1.FuturesClient = aster.NewFuturesClient(config1.AsterAPIKey, config1.AsterSecretKey)
	
	DiagnoseAsterClient(client1, config1)
	
	// 场景2: API Key 为空
	fmt.Println("\n\n=== 场景2: API Key 为空 ===")
	config2 := &Config{
		AsterAPIKey:    "",
		AsterSecretKey: "test-secret",
	}
	
	client2 := &AsterClient{}
	client2.FuturesClient = aster.NewFuturesClient(config2.AsterAPIKey, config2.AsterSecretKey)
	
	DiagnoseAsterClient(client2, config2)
	
	// 场景3: 客户端未初始化
	fmt.Println("\n\n=== 场景3: 客户端未初始化 ===")
	config3 := &Config{
		AsterAPIKey:    "test-key",
		AsterSecretKey: "test-secret",
	}
	
	client3 := &AsterClient{}
	// 忘记初始化 FuturesClient
	
	DiagnoseAsterClient(client3, config3)
}