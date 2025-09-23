package main

import (
	"context"
	"fmt"
	"log"

	"github.com/drinkthere/go-aster/v2"
	"github.com/drinkthere/go-aster/v2/futures"
)

func main() {
	// 设置你的API凭证
	apiKey := "your-api-key"
	secretKey := "your-secret-key"

	// 创建期货客户端
	client := aster.NewFuturesClient(apiKey, secretKey)

	// 修改ETHUSDT的杠杆为20倍
	leverageService := &futures.ChangeLeverageService{C: client}
	result, err := leverageService.
		Symbol("ETHUSDT").
		Leverage(20).
		Do(context.Background())

	if err != nil {
		log.Fatalf("修改杠杆失败: %v", err)
	}

	// 打印结果
	fmt.Printf("成功修改杠杆!\n")
	fmt.Printf("交易对: %s\n", result.Symbol)
	fmt.Printf("新杠杆: %d\n", result.Leverage)
	fmt.Printf("最大名义价值: %s\n", result.MaxNotionalValue)
}