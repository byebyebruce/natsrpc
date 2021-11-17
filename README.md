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
> NATSRPC 是一个基于[NATS](https://nats.io/)作为消息通信，使用gRPC接口文件生成Client和Server代码的rpc框架

## Feature
* 使用简单，不需要服务发现
* 使用gRPC接口定义方式，接口定义清晰，学习成本低
* 代码生成器一键生成
* 支持空间隔离
* 支持定向发送也支持负载均衡(nats的同组内随机)
* 不用手动定义subject

## 安装工具
* protoc v3.17.3
    * [Linux](https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-linux-x86_64.zip)
    * [MacOS](https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-osx-x86_64.zip)
    * [Windows](https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-win64.zip)

* 插件
```
go get -u github.com/golang/protobuf/protoc-gen-go@v1.3.5
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@v1.16.0
go get -u github.com/byebyebruce/natsrpc/tool/cmd/protoc-gen-natsrpc
```

## 快速使用
* 需要先启动nats-server
1. 创建工程
`go mod init natsrpc_test`
2. 引用包 `go get github.com/byebyebruce/natsrpc`
3. 定义服务接口
```
syntax = "proto3";

package natsrpc_test;
option go_package = "github.com/byebyebruce/example/natsrpc_test;natsrpc_test";

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
[更多示例](example)

## 压测工具
1. 广播 `go run bench/pub/main.go -server=nats://127.0.0.1:4222`

2. 请求 `go run bench/req/main.go -server=nats://127.0.0.1:4222`

## TODO
- [x] service 定义文件改成gRPC标准
- [x] 支持支持返回错误
- [ ] 支持goroutine池
- [ ] 收发包顺序