package main

import (
	"context"
	"fmt"
	"log"

	aster "github.com/drinkthere/go-aster/v2"
	"github.com/drinkthere/go-aster/v2/futures"
)

func main() {
	// 1. 基本的期货客户端创建（只使用API Key和Secret）
	futuresClient := aster.NewFuturesClient("your-api-key", "your-secret-key")

	// 2. 带选项的期货客户端创建
	futuresClientWithOptions := aster.NewFuturesClient("your-api-key", "your-secret-key",
		aster.WithDebug(true),                          // 开启调试模式
		aster.WithLocalAddress("192.168.1.100"),       // 绑定本地IP
		aster.WithBaseURL("https://fapi.asterdex.com"), // 自定义API地址
	)
	_ = futuresClientWithOptions // 使用变量避免编译警告

	// 3. 使用期货客户端进行REST API调用
	fmt.Println("=== 期货REST API示例 ===")

	// 获取账户信息
	accountService := &futures.GetAccountService{C: futuresClient}
	account, err := accountService.Do(context.Background())
	if err != nil {
		log.Printf("获取账户信息失败: %v", err)
	} else {
		fmt.Printf("账户余额: %+v\n", account.Assets)
		fmt.Printf("持仓信息: %+v\n", account.Positions)
	}

	// 获取交易所信息
	exchangeService := &futures.ExchangeInfoService{C: futuresClient}
	exchangeInfo, err := exchangeService.Do(context.Background())
	if err != nil {
		log.Printf("获取交易所信息失败: %v", err)
	} else {
		fmt.Printf("交易对数量: %d\n", len(exchangeInfo.Symbols))
	}

	// 下单示例
	orderService := &futures.CreateOrderService{C: futuresClient}
	order, err := orderService.
		Symbol("BTCUSDT").
		Side("BUY").
		Type("LIMIT").
		Quantity("0.001").
		Price("30000").
		TimeInForce("GTC").
		Do(context.Background())
	if err != nil {
		log.Printf("下单失败: %v", err)
	} else {
		fmt.Printf("下单成功: OrderID=%d, ClientOrderID=%s\n", order.OrderID, order.ClientOrderID)
	}

	// 查询持仓
	positionService := &futures.GetPositionRiskService{C: futuresClient}
	positions, err := positionService.Do(context.Background())
	if err != nil {
		log.Printf("查询持仓失败: %v", err)
	} else {
		fmt.Printf("持仓数量: %d\n", len(positions))
		for _, pos := range positions {
			if pos.PositionAmt != "0" {
				fmt.Printf("持仓: %s, 数量: %s, 未实现盈亏: %s\n", 
					pos.Symbol, pos.PositionAmt, pos.UnRealizedProfit)
			}
		}
	}

	// 4. 使用新的FuturesClient包装器（推荐方式）
	fmt.Println("\n=== 使用FuturesClient包装器 ===")
	
	futuresWrapper := aster.NewFutures("your-api-key", "your-secret-key",
		aster.WithDebug(true),
		aster.WithLocalAddress("192.168.1.100"),
	)

	// 5. WebSocket示例
	fmt.Println("\n=== 期货WebSocket示例 ===")

	// 订阅单个交易对的book ticker
	_, stopC1, err := futuresWrapper.WsBookTickerServe("BTCUSDT", 
		func(event *aster.WsBookTickerEvent) {
			fmt.Printf("BTC期货 Book Ticker: 买价=%s, 卖价=%s\n", 
				event.BestBidPrice, event.BestAskPrice)
		},
		func(err error) {
			log.Printf("BTC WebSocket错误: %v", err)
		})
	if err != nil {
		log.Printf("启动BTC WebSocket失败: %v", err)
	}

	// 订阅多个交易对的book ticker（组合流）
	symbols := []string{"BTCUSDT", "ETHUSDT", "ADAUSDT"}
	_, stopC2, err := futuresWrapper.WsCombinedBookTickerServe(symbols,
		func(event *aster.WsBookTickerEvent) {
			fmt.Printf("期货组合流 %s: 买价=%s, 卖价=%s\n", 
				event.Symbol, event.BestBidPrice, event.BestAskPrice)
		},
		func(err error) {
			log.Printf("组合流WebSocket错误: %v", err)
		})
	if err != nil {
		log.Printf("启动组合流WebSocket失败: %v", err)
	}

	// 订阅标记价格流
	_, stopC3, err := futuresWrapper.WsMarkPriceServe("BTCUSDT",
		func(event *aster.WsFuturesMarkPriceEvent) {
			fmt.Printf("BTC标记价格: %s, 资金费率: %s\n", 
				event.MarkPrice, event.FundingRate)
		},
		func(err error) {
			log.Printf("标记价格WebSocket错误: %v", err)
		})
	if err != nil {
		log.Printf("启动标记价格WebSocket失败: %v", err)
	}

	// 用户数据流（需要先创建listen key）
	listenKeyService := &futures.StartUserStreamService{C: futuresClient}
	listenKey, err := listenKeyService.Do(context.Background())
	if err != nil {
		log.Printf("创建用户数据流失败: %v", err)
	} else {
		_, stopC4, err := futuresWrapper.WsUserDataServe(listenKey,
			func(event *aster.WsFuturesUserDataEvent) {
				switch event.Event {
				case "ACCOUNT_UPDATE":
					fmt.Printf("账户更新: %+v\n", event.AccountUpdate)
				case "ORDER_TRADE_UPDATE":
					fmt.Printf("订单更新: %s %s %s\n", 
						event.OrderUpdate.Symbol, 
						event.OrderUpdate.Side, 
						event.OrderUpdate.OrderStatus)
				}
			},
			func(err error) {
				log.Printf("用户数据流错误: %v", err)
			})
		if err != nil {
			log.Printf("启动用户数据流失败: %v", err)
		} else {
			defer close(stopC4)
		}
	}

	// 6. 指定不同本地IP的WebSocket连接
	_, stopC5, err := futuresWrapper.WsBookTickerServeWithLocalAddr("ETHUSDT",
		func(event *aster.WsBookTickerEvent) {
			fmt.Printf("ETH期货(自定义IP): 买价=%s, 卖价=%s\n", 
				event.BestBidPrice, event.BestAskPrice)
		},
		func(err error) {
			log.Printf("ETH自定义IP WebSocket错误: %v", err)
		},
		"192.168.1.101") // 使用不同的本地IP
	if err != nil {
		log.Printf("启动ETH自定义IP WebSocket失败: %v", err)
	}

	fmt.Println("\nWebSocket连接已启动，按回车键停止...")
	fmt.Scanln()

	// 清理资源
	if stopC1 != nil {
		close(stopC1)
	}
	if stopC2 != nil {
		close(stopC2)
	}
	if stopC3 != nil {
		close(stopC3)
	}
	if stopC5 != nil {
		close(stopC5)
	}

	fmt.Println("所有连接已停止")
}