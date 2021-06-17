# NATSRPC
NATSRPC 是一个基于nats的简单rpc

## Feature
* 使用简单，不需要服务发现
* 代码生成器生成client和server代码
* 支持空间隔离
* 支持定向发送也支持负载均衡(nats的同组内随机)
* 不用手动定义subject
* 支持单协程回调(适用于逻辑单协程模型)

## 使用
1. 引用包 `go get github.com/byebyebruce/natsrpc`
2. 编译代码生成器 go get github.com/byebyebruce/natsrpc/cmd
3. 编写service
```go
package helloworld

import (
	"context"

	"github.com/byebyebruce/natsrpc/testdata/pb"
)

// Greeter hello
type Greeter interface {
	HiAll(ctx context.Context, req *pb.HelloRequest)
	AreYouOK(ctx context.Context, req *pb.HelloRequest, repl *pb.HelloReply)
}
```
4. 生成代码
```shell
natsrpc_codegen -s="testdata/greeter.go"
```

## 示例
* [Client](example/client/main.go)
* [Server](example/server/main.go)
* [API](example/api/greeter.go)
> 运行示例需要部署gnatsd，如果没有可以临时启动`go run cmd/simple_natsserver/main.go`