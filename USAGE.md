# 如何在项目中使用 go-aster SDK

## 方法1: 使用特定版本 (推荐)

在你的项目中运行：

```bash
go get github.com/drinkthere/go-aster@v0.0.2
```

或者在 `go.mod` 文件中添加：

```go
require github.com/drinkthere/go-aster v0.0.2
```

然后运行 `go mod download`

## 方法2: 使用最新版本

```bash
go get github.com/drinkthere/go-aster@latest
```

## 示例代码

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    aster "github.com/drinkthere/go-aster/v3"
)

func main() {
    // 创建客户端
    client := aster.NewClient(userAddress, signerAddress, privateKey)
    
    // 获取服务器时间
    serverTime, err := client.NewServerTimeService().Do(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Server time: %d\n", serverTime)
}
```

## 版本管理

- 使用 `go get github.com/drinkthere/go-aster@v0.0.2` 获取特定版本
- 使用 `go get -u github.com/drinkthere/go-aster` 更新到最新版本
- 使用 `go list -m -versions github.com/drinkthere/go-aster` 查看所有可用版本

## WebSocket 使用示例

### 订阅市场数据

```go
// 订阅深度数据
doneC, stopC, err := aster.WsDepthServe("BTCUSDT", depthHandler, errHandler)

// 订阅 Book Ticker（最优买卖价）
doneC, stopC, err := aster.WsBookTickerServe("BTCUSDT", bookTickerHandler, errHandler)

// 订阅多个交易对
symbols := []string{"BTCUSDT", "ETHUSDT"}
doneC, stopC, err := aster.WsCombinedBookTickerServe(symbols, bookTickerHandler, errHandler)
```

### 订阅用户数据

```go
// 1. 创建 listen key
listenKey, err := client.NewStartUserStreamService().Do(ctx)

// 2. 订阅用户数据流
doneC, stopC, err := aster.WsUserDataServe(listenKey, userDataHandler, errHandler)

// 3. 定期更新 listen key (每30分钟)
client.NewKeepaliveUserStreamService().ListenKey(listenKey).Do(ctx)

// 4. 关闭时删除 listen key
client.NewCloseUserStreamService().ListenKey(listenKey).Do(ctx)
```

### 处理器示例

```go
// Book Ticker 处理器
bookTickerHandler := func(event *aster.WsBookTickerEvent) {
    fmt.Printf("Symbol: %s, Bid: %s @ %s, Ask: %s @ %s\n",
        event.Symbol, event.BestBidPrice, event.BestBidQty,
        event.BestAskPrice, event.BestAskQty)
}

// 用户数据处理器
userDataHandler := func(event *aster.WsUserDataEvent) {
    switch event.Event {
    case "ORDER_TRADE_UPDATE":
        // 处理订单更新
        fmt.Printf("Order %s: %s\n", event.OrderUpdate.Status, event.OrderUpdate.Symbol)
    case "ACCOUNT_UPDATE":
        // 处理账户更新（仓位、余额变化）
        for _, pos := range event.AccountUpdate.Positions {
            fmt.Printf("Position %s: %s\n", pos.Symbol, pos.PositionAmount)
        }
    }
}
```