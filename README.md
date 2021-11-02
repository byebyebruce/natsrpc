# NATSRPC
[![Go doc](https://img.shields.io/badge/go.dev-reference-brightgreen?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/byebyebruce/natsrpc)
[![Go Report](https://goreportcard.com/badge/github.com/byebyebruce/natsrpc)](https://goreportcard.com/report/github.com/byebyebruce/natsrpc)

```
  _   _       _______ _____   _____  _____   _____ 
 | \ | |   /\|__   __/ ____| |  __ \|  __ \ / ____|
 |  \| |  /  \  | | | (___   | |__) | |__) | |     
 | . ` | / /\ \ | |  \___ \  |  _  /|  ___/| |     
 | |\  |/ ____ \| |  ____) | | | \ \| |    | |____ 
 |_| \_/_/    \_\_| |_____/  |_|  \_\_|     \_____|
```
> NATSRPC 是一个基于nats作为消息通信，用grpc接口定义文件生成client和server代码的rpc

## Feature
* 使用简单，不需要服务发现
* 代码生成器一键生成
* 支持空间隔离
* 支持定向发送也支持负载均衡(nats的同组内随机)
* 不用手动定义subject
* 使用grpc接口定义方式

## 使用
1. 引用包 `go get github.com/byebyebruce/natsrpc`
2. 编译代码生成器 `go get github.com/byebyebruce/natsrpc/tool/cmd/protoc-gen-natsrpc`
3. 定义服务接口[示例](example/pb)

4. 生成客户端和服务端代码
```shell
natsrpc_codegen -s="greeter.go" # 客户端代码只生成同步接口
natsrpc_codegen -s="greeter.go -cm=1" # 客户端代码只生成同步接口
natsrpc_codegen -s="greeter.go -cm=2" # 客户端代码生成同步异步接口
```
5. 写服务实现[示例](example/example_greeter.go)
## 示例
* [Client](example/client/main.go)
* [Server](example/server/main.go)
> 运行示例需要部署gnatsd，如果没有可以临时启动`go run tool/cmd/simple_natsserver/main.go`

## 压测工具
1. 广播 `go run bench/pub/main.go -server=nats://127.0.0.1:4222`

2. 请求 `go run bench/req/main.go -server=nats://127.0.0.1:4222`

## TODO
[ ] 不用NatsEncoder 
[ ] service 定义文件改成gRPC标准
[ ] 考虑收发包顺序