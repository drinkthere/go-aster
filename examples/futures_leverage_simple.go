package main

import (
	"context"
	"fmt"
	"log"

	"github.com/drinkthere/go-aster/v2"
	"github.com/drinkthere/go-aster/v2/common"
	"github.com/drinkthere/go-aster/v2/futures"
)

func main() {
	// 创建期货客户端（使用 API Key + Secret Key）
	apiKey := "your-api-key"
	secretKey := "your-secret-key"
	
	client := aster.NewFuturesClient(apiKey, secretKey)
	
	// 如果需要使用测试网
	// client := aster.NewFuturesClient(apiKey, secretKey,
	//     aster.WithBaseURL("https://testnet.asterdex.com"))
	
	// 修改单个交易对的杠杆
	fmt.Println("=== 修改 BTCUSDT 杠杆 ===")
	if err := changeLeverage(client, "BTCUSDT", 20); err != nil {
		log.Printf("修改失败: %v", err)
		return
	}
	
	// 批量修改多个交易对
	fmt.Println("\n=== 批量修改杠杆 ===")
	symbols := map[string]int{
		"BTCUSDT": 10,
		"ETHUSDT": 15,
		"BNBUSDT": 8,
	}
	
	for symbol, leverage := range symbols {
		if err := changeLeverage(client, symbol, leverage); err != nil {
			log.Printf("%s 修改失败: %v", symbol, err)
		}
	}
}

// changeLeverage 修改指定交易对的杠杆倍数
func changeLeverage(client *aster.BaseClient, symbol string, leverage int) error {
	service := &futures.ChangeLeverageService{C: client}
	
	result, err := service.
		Symbol(symbol).
		Leverage(leverage).
		Do(context.Background())
	
	if err != nil {
		// 处理API错误
		if apiErr, ok := err.(*common.APIError); ok {
			return fmt.Errorf("API错误 [%d]: %s", apiErr.Code, apiErr.Message)
		}
		return err
	}
	
	fmt.Printf("✓ %s 杠杆修改成功: %dx (最大名义价值: %s)\n", 
		result.Symbol, result.Leverage, result.MaxNotionalValue)
	
	return nil
}

// 常见错误码说明：
// -4028: 杠杆值无效（超出允许范围）
// -2015: 无效的API-key, IP, 或权限
// -1021: 时间戳超出接收窗口
// -4000: 无效的参数
// -5021: 该交易对不支持修改杠杆