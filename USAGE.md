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