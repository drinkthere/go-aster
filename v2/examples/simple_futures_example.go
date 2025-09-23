package main

import (
	"context"
	"fmt"
	"log"

	aster "github.com/drinkthere/go-aster/v2"
	"github.com/drinkthere/go-aster/v2/futures"
)

func main() {
	// 创建期货客户端
	client := aster.NewFuturesClient("your-api-key", "your-secret-key",
		aster.WithDebug(true), // 开启调试模式查看请求详情
	)

	// 示例1: 获取账户信息
	fmt.Println("=== 获取账户信息 ===")
	accountService := &futures.GetAccountService{C: client}
	account, err := accountService.Do(context.Background())
	if err != nil {
		log.Printf("获取账户信息失败: %v", err)
	} else {
		fmt.Printf("总余额: %s USDT\n", account.TotalWalletBalance)
		fmt.Printf("可用余额: %s USDT\n", account.AvailableBalance)
		
		// 显示非零持仓
		for _, pos := range account.Positions {
			if pos.PositionAmt != "0" {
				fmt.Printf("持仓: %s, 数量: %s, 方向: %s, 未实现盈亏: %s\n",
					pos.Symbol, pos.PositionAmt, pos.PositionSide, pos.UnRealizedProfit)
			}
		}
	}

	// 示例2: 修改杠杆倍数
	fmt.Println("\n=== 修改杠杆倍数 ===")
	leverageService := &futures.ChangeLeverageService{C: client}
	leverageResp, err := leverageService.
		Symbol("BTCUSDT").
		Leverage(20).
		Do(context.Background())
	if err != nil {
		log.Printf("修改杠杆失败: %v", err)
	} else {
		fmt.Printf("BTCUSDT杠杆已设置为: %d倍\n", leverageResp.Leverage)
	}

	// 示例3: 下限价单
	fmt.Println("\n=== 下限价单 ===")
	orderService := &futures.CreateOrderService{C: client}
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
		fmt.Printf("下单成功! OrderID: %d, ClientOrderID: %s\n", 
			order.OrderID, order.ClientOrderID)
	}

	// 示例4: 查询开放订单
	fmt.Println("\n=== 查询开放订单 ===")
	openOrdersService := &futures.ListOpenOrdersService{C: client}
	openOrders, err := openOrdersService.Do(context.Background())
	if err != nil {
		log.Printf("查询开放订单失败: %v", err)
	} else {
		fmt.Printf("当前开放订单数量: %d\n", len(openOrders))
		for _, order := range openOrders {
			fmt.Printf("订单: %s %s %s, 价格: %s, 数量: %s, 状态: %s\n",
				order.Symbol, order.Side, order.Type, 
				order.Price, order.OrigQty, order.Status)
		}
	}

	// 示例5: 取消所有订单
	fmt.Println("\n=== 取消所有订单 ===")
	cancelAllService := &futures.CancelAllOpenOrdersService{C: client}
	err = cancelAllService.Symbol("BTCUSDT").Do(context.Background())
	if err != nil {
		log.Printf("取消订单失败: %v", err)
	} else {
		fmt.Printf("成功取消BTCUSDT的所有订单\n")
	}

	fmt.Println("\n期货客户端示例完成!")
}