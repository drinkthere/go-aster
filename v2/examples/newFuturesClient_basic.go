package main

import (
	"context"
	"fmt"
	"log"

	aster "github.com/drinkthere/go-aster/v2"
	"github.com/drinkthere/go-aster/v2/futures"
)

func main() {
	// 使用aster.NewFuturesClient创建期货客户端
	client := aster.NewFuturesClient("your-api-key", "your-secret-key")

	// 示例1: 获取交易所信息
	exchangeService := &futures.ExchangeInfoService{C: client}
	exchangeInfo, err := exchangeService.Do(context.Background())
	if err != nil {
		log.Printf("获取交易所信息失败: %v", err)
	} else {
		fmt.Printf("交易对数量: %d\n", len(exchangeInfo.Symbols))
		// 显示前3个交易对
		for i, symbol := range exchangeInfo.Symbols {
			if i >= 3 {
				break
			}
			fmt.Printf("交易对: %s, 状态: %s\n", symbol.Symbol, symbol.Status)
		}
	}

	// 示例2: 获取账户信息
	accountService := &futures.GetAccountService{C: client}
	account, err := accountService.Do(context.Background())
	if err != nil {
		log.Printf("获取账户信息失败: %v", err)
	} else {
		fmt.Printf("总余额: %s USDT\n", account.TotalWalletBalance)
		fmt.Printf("可用余额: %s USDT\n", account.AvailableBalance)
	}

	// 示例3: 设置杠杆
	leverageService := &futures.ChangeLeverageService{C: client}
	leverageResp, err := leverageService.
		Symbol("BTCUSDT").
		Leverage(10).
		Do(context.Background())
	if err != nil {
		log.Printf("设置杠杆失败: %v", err)
	} else {
		fmt.Printf("BTCUSDT杠杆设置为: %d倍\n", leverageResp.Leverage)
	}

	fmt.Println("期货客户端基础示例完成!")
}