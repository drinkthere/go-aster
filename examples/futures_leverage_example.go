package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/drinkthere/go-aster/v2"
	"github.com/drinkthere/go-aster/v2/common"
	"github.com/drinkthere/go-aster/v2/futures"
)

func main() {
	// 创建期货客户端（需要Web3签名认证）
	userAddress := "0xYourUserAddress"     // 替换为你的用户地址
	signerAddress := "0xYourSignerAddress" // 替换为你的签名地址
	privateKey := "yourPrivateKeyWithout0x" // 替换为你的私钥（不要0x前缀）
	
	// 创建客户端
	client := aster.NewFuturesClient(userAddress, signerAddress, privateKey)
	
	// 也可以使用测试网
	// client := aster.NewFuturesClient(userAddress, signerAddress, privateKey, 
	//     aster.WithBaseURL("https://testnet.asterdex.com"))
	
	// Example 1: 查询当前持仓信息（查看当前杠杆）
	fmt.Println("=== 查询当前持仓信息 ===")
	positionService := &futures.GetPositionRiskService{C: client}
	positionService.Symbol("BTCUSDT") // 可以指定交易对，不指定则返回所有
	
	positions, err := positionService.Do(context.Background())
	if err != nil {
		log.Printf("获取持仓信息失败: %v", err)
	} else {
		for _, pos := range positions {
			fmt.Printf("交易对: %s, 持仓方向: %s, 数量: %s, 杠杆: %d\n", 
				pos.Symbol, pos.PositionSide, pos.PositionAmt, pos.Leverage)
		}
	}
	
	// Example 2: 修改杠杆倍数
	fmt.Println("\n=== 修改杠杆倍数 ===")
	symbol := "BTCUSDT"
	newLeverage := 10 // 新的杠杆倍数
	
	leverageService := &futures.ChangeLeverageService{C: client}
	leverageResult, err := leverageService.
		Symbol(symbol).
		Leverage(newLeverage).
		Do(context.Background())
	
	if err != nil {
		log.Printf("修改杠杆失败: %v", err)
		// 如果是API错误，可以获取具体错误信息
		if apiErr, ok := err.(*common.APIError); ok {
			log.Printf("API错误码: %d, 错误信息: %s", apiErr.Code, apiErr.Message)
		}
		return
	}
	
	fmt.Printf("杠杆修改成功!\n")
	fmt.Printf("交易对: %s\n", leverageResult.Symbol)
	fmt.Printf("新杠杆倍数: %d\n", leverageResult.Leverage)
	fmt.Printf("最大名义价值: %s\n", leverageResult.MaxNotionalValue)
	
	// Example 3: 批量修改多个交易对的杠杆
	fmt.Println("\n=== 批量修改多个交易对杠杆 ===")
	symbols := []string{"BTCUSDT", "ETHUSDT", "BNBUSDT"}
	targetLeverage := 5
	
	for _, sym := range symbols {
		fmt.Printf("\n修改 %s 的杠杆为 %dx...\n", sym, targetLeverage)
		
		service := &futures.ChangeLeverageService{C: client}
		result, err := service.
			Symbol(sym).
			Leverage(targetLeverage).
			Do(context.Background())
		
		if err != nil {
			log.Printf("修改 %s 杠杆失败: %v", sym, err)
			continue
		}
		
		fmt.Printf("成功! 新杠杆: %d, 最大名义价值: %s\n", 
			result.Leverage, result.MaxNotionalValue)
	}
	
	// Example 4: 根据账户余额计算合适的杠杆
	fmt.Println("\n=== 根据账户余额计算合适的杠杆 ===")
	
	// 先获取账户信息
	accountService := &futures.GetAccountService{C: client}
	account, err := accountService.Do(context.Background())
	if err != nil {
		log.Printf("获取账户信息失败: %v", err)
		return
	}
	
	fmt.Printf("账户总余额: %s\n", account.TotalWalletBalance)
	fmt.Printf("可用余额: %s\n", account.AvailableBalance)
	
	// 根据风险管理原则，建议单个仓位不超过总余额的某个百分比
	// 这里只是示例，实际使用时需要根据自己的风险管理策略调整
	maxRiskPercentage := 0.1 // 10% 风险
	
	// Example 5: 修改杠杆前的检查
	fmt.Println("\n=== 修改杠杆前的安全检查 ===")
	
	// 检查是否有未成交订单
	openOrdersService := &futures.ListOpenOrdersService{C: client}
	openOrdersService.Symbol("BTCUSDT")
	openOrders, err := openOrdersService.Do(context.Background())
	
	if err == nil && len(openOrders) > 0 {
		fmt.Printf("警告: %s 有 %d 个未成交订单，建议先取消订单再修改杠杆\n", 
			"BTCUSDT", len(openOrders))
		
		// 可以选择取消所有订单
		// cancelService := &futures.CancelAllOpenOrdersService{C: client}
		// cancelService.Symbol("BTCUSDT")
		// err = cancelService.Do(context.Background())
	}
	
	// 检查是否有持仓
	posService := &futures.GetPositionRiskService{C: client}
	posService.Symbol("BTCUSDT")
	positions2, err := posService.Do(context.Background())
	
	if err == nil {
		for _, pos := range positions2 {
			posAmt, _ := strconv.ParseFloat(pos.PositionAmt, 64)
			if posAmt != 0 {
				fmt.Printf("警告: %s 有持仓，数量: %s，修改杠杆可能影响保证金要求\n",
					pos.Symbol, pos.PositionAmt)
			}
		}
	}
	
	// Example 6: 错误处理最佳实践
	fmt.Println("\n=== 错误处理示例 ===")
	
	// 尝试设置一个可能无效的杠杆值
	invalidService := &futures.ChangeLeverageService{C: client}
	_, err = invalidService.
		Symbol("BTCUSDT").
		Leverage(200). // 可能超过最大允许杠杆
		Do(context.Background())
	
	if err != nil {
		fmt.Println("预期的错误发生:")
		// 处理不同类型的错误
		switch e := err.(type) {
		case *common.APIError:
			fmt.Printf("API错误 - 代码: %d, 消息: %s\n", e.Code, e.Message)
			// 常见错误码：
			// -4028: 杠杆值无效
			// -2015: 无效的API-key, IP, 或权限
			// -1021: 时间戳超出接收窗口
		default:
			fmt.Printf("其他错误: %v\n", err)
		}
	}
	
	fmt.Println("\n=== 注意事项 ===")
	fmt.Println("1. 修改杠杆可能影响现有持仓的保证金要求")
	fmt.Println("2. 建议在没有未成交订单时修改杠杆")
	fmt.Println("3. 不同交易对可能有不同的最大杠杆限制")
	fmt.Println("4. 杠杆越高，风险越大，请谨慎操作")
}