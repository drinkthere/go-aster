package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/drinkthere/go-aster/v2"
	"github.com/drinkthere/go-aster/v2/common"
	"github.com/drinkthere/go-aster/v2/futures"
)

func main() {
	// 根据你提供的信息，这可能是：
	// - 用户地址或签名地址（需要加 0x 前缀）
	// - 私钥（不需要 0x 前缀）
	
	// 尝试方案1：假设第一个是地址，第二个是私钥
	userOrSignerAddress := "0xcc7ac0adbd279bf10db8206d8e0d8e71ce33aa6a301f64cafc257889c35361ac"
	privateKey := "cc7ac0adbd279bf10db8206d8e0d8e71ce33aa6a301f64cafc257889c35361ac"
	
	// 如果上面的地址太长（正常的以太坊地址是 40 个字符），可能需要截取
	// 标准以太坊地址格式：0x + 40个十六进制字符
	if len(userOrSignerAddress) > 42 {
		fmt.Printf("警告：地址长度异常，可能需要调整\n")
	}
	
	// 创建 Web3 认证的期货客户端
	fmt.Println("=== 使用 Web3 认证 ===")
	
	// 这里假设用户地址和签名地址相同
	client := aster.NewFuturesClientWithWeb3(
		userOrSignerAddress,  // user address
		userOrSignerAddress,  // signer address
		privateKey,          // private key (without 0x)
		aster.WithDebug(true), // 开启调试
	)
	
	// 测试基本连接
	fmt.Println("\n=== 测试 Ping ===")
	pingService := &futures.PingService{C: client}
	if err := pingService.Do(context.Background()); err != nil {
		log.Printf("Ping 失败: %v", err)
	} else {
		fmt.Println("Ping 成功!")
	}
	
	// 测试需要签名的接口
	fmt.Println("\n=== 测试获取账户信息 ===")
	accountService := &futures.GetAccountService{C: client}
	account, err := accountService.Do(context.Background())
	if err != nil {
		log.Printf("获取账户信息失败: %v", err)
		if apiErr, ok := err.(*common.APIError); ok {
			log.Printf("API 错误码: %d, 错误信息: %s", apiErr.Code, apiErr.Message)
		}
	} else {
		fmt.Printf("账户信息获取成功!\n")
		fmt.Printf("总余额: %s\n", account.TotalWalletBalance)
	}
	
	// 如果上面的方案不行，尝试其他组合
	fmt.Println("\n=== 尝试其他可能的组合 ===")
	
	// 方案2：可能你的凭证实际上是两个地址
	// 或者第一个是 API key 形式的地址表示
	possibleAddresses := []string{
		"cc7ac0adbd279bf10db8206d8e0d8e71ce33aa6a301f64cafc257889c35361ac",
		"0xcc7ac0adbd279bf10db8206d8e0d8e71ce33aa6a301f64cafc257889c35361ac",
	}
	
	// 尝试提取可能的以太坊地址（前20字节）
	for i, addr := range possibleAddresses {
		addr = strings.TrimPrefix(addr, "0x")
		if len(addr) >= 40 {
			ethAddress := "0x" + addr[:40]
			fmt.Printf("尝试地址 %d: %s\n", i+1, ethAddress)
		}
	}
	
	fmt.Println("\n说明：")
	fmt.Println("1. Aster 期货 API 使用 Web3 签名认证")
	fmt.Println("2. 需要提供：用户地址、签名地址、私钥")
	fmt.Println("3. 请确认你的凭证类型和格式")
}