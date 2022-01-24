# NATSRPC

```
  _   _       _______ _____   _____  _____   _____ 
 | \ | |   /\|__   __/ ____| |  __ \|  __ \ / ____|
 |  \| |  /  \  | | | (___   | |__) | |__) | |     
 | . ` | / /\ \ | |  \___ \  |  _  /|  ___/| |     
 | |\  |/ ____ \| |  ____) | | | \ \| |    | |____ 
 |_| \_/_/    \_\_| |_____/  |_|  \_\_|     \_____|
```
> NATSRPC 是一个基于[NATS](https://nats.io/)作为消息通信，使用[gRPC](https://www.grpc.io/)的方式来定义接口的RPC框架
## Why NATSRPC  
NATS收发消息需要手动定义subject，request，reply，handler等繁琐且易出错的代码。gRPC需要用服务发现到endpoint才能发送请求。NATRPC的目的就是要像gRPC一样定义接口，像NATS一样不关心具体网络位置，只需要监听和发送就能完成RPC调用。
## Feature
* 使用简单，不需要服务发现
* 使用gRPC接口定义方式，接口定义清晰，学习成本低
* 代码生成器一键生成
* 支持空间隔离
* 支持定向发送也支持负载均衡(nats的同组内随机)
* 支持Header和返回Error

## 安装工具
* protoc(v3.17.3) [Linux](https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-linux-x86_64.zip)/[MacOS](https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-osx-x86_64.zip)/[Windows](https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-win64.zip)

* protoc插件
```
go install github.com/golang/protobuf/protoc-gen-go@v1.3.5
go install github.com/byebyebruce/natsrpc/tool/cmd/protoc-gen-natsrpc@latest
```

## 快速使用
* 启动nats-server(没有部署好的nats-server可以`go run example/cmd/simple_natsserver/main.go`)
1. 创建工程
`go mod init natsrpc_test`
2. 引用包 `go get github.com/byebyebruce/natsrpc`
3. 定义服务接口
```
syntax = "proto3";

package natsrpc_test;
option go_package = "github.com/byebyebruce/natsrpc/example/natsrpc_test;natsrpc_test";

message HelloRequest {
  string name = 1;
}

message HelloReply {
  string message = 1;
}

service Greeter {
  rpc Hello (HelloRequest) returns (HelloReply) {}
}
```

4. 生成客户端和服务端代码
```shell
protoc --proto_path=. --go_out=paths=source_relative:. --natsrpc_out=paths=source_relative:. *.proto
```
5. 实现接口
```
type Greeter interface {
	// Hello
	Hello(ctx context.Context, req *natsrpc_test.HelloRequest) (*natsrpc_test.HelloReply, error)
}
```
6. [main.go](example/0.main_test.go) 启动server和client
## 更多示例
1. [请求](example/1.request_test.go)
2. [广播](example/2.publish_test.go)
3. [异步请求](example/3.asyncclient_test.go)
4. [异步回复](example/4.asyncservice_test.go)
5. [请求头](example/5.header_test.go)

## 压测工具
1. 广播 `go run bench/pub/main.go -server=nats://127.0.0.1:4222`

2. 请求 `go run bench/req/main.go -server=nats://127.0.0.1:4222`

## TODO
- [x] service 定义文件改成gRPC标准
- [x] 支持返回错误
- [x] 支持Header
- [x] 生成Client接口
- [ ] 支持goroutine池
- [ ] 支持中间件
