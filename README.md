```
  _   _       _______ _____   _____  _____   _____ 
 | \ | |   /\|__   __/ ____| |  __ \|  __ \ / ____|
 |  \| |  /  \  | | | (___   | |__) | |__) | |     
 | . ` | / /\ \ | |  \___ \  |  _  /|  ___/| |     
 | |\  |/ ____ \| |  ____) | | | \ \| |    | |____ 
 |_| \_/_/    \_\_| |_____/  |_|  \_\_|     \_____|
```

## What is NATSRPC
> NATSRPC 是一个基于[NATS](https://nats.io/)作为消息通信，使用[gRPC](https://www.grpc.io/)的方式来定义接口的RPC框架

![GitHub release (with filter)](https://img.shields.io/github/v/release/byebyebruce/natsrpc)
![](https://hits.sh/github.com/byebyebruce/natsrpc/doc/hits.svg?label=visit)

## Motivation  
NATS收发消息需要手动定义subject，request，reply，handler等繁琐且易出错的代码。
gRPC需要连接到可知endpoint才能发送请求。
NATRPC的目的就是要像gRPC一样定义接口，像NATS一样不关心具体网络位置，只需要监听和发送就能完成RPC调用。

## Feature
* 使用gRPC接口定义方式，使用简单，一键生成代码
* 支持空间隔离,也可以指定id发送
* 多服务可以负载均衡(nats的同组内随机)
* 支持Header和返回Error
* 支持单协程和多协程handle
* 支持中间件
* 支持延迟回复消息
* 支持自定义编码器

## How It Works
上层通过Server、Service、Client对nats.Conn和Subscription进行封装。  
底层通过nats的request和publish来传输消息。一个Service会创建一个以service name为subject的Subscription，如果有publish方法会在创建一个用于接收publish的sub。  
Client发请求时会的subject是service 的name，并且nats msg的header传递method name。  
Service收到消息后取出method name，然后调用对应的handler，handler返回的结果会通过nats msg的reply subject返回给Client。

## Install Tools
1. protoc(v3.17.3) [Linux](https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-linux-x86_64.zip)/[MacOS](https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-osx-x86_64.zip)/[Windows](https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-win64.zip)
2. protoc-gen-go `go install github.com/golang/protobuf/protoc-gen-go@latest`
3. protoc-gen-natsrpc `go install github.com/byebyebruce/natsrpc/cmd/protoc-gen-natsrpc@v0.7.0`

## Quick Start
* [nats-server](https://github.com/nats-io/nats-server/releases)>=2.2.0
1. 引用包
   ```shell
   go get github.com/byebyebruce/natsrpc
   ```
2. 定义服务接口 example.proto
    ```
    syntax = "proto3";

    package example;
    option go_package = "github.com/byebyebruce/natsrpc/example;example";

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
   
3. 生成客户端和服务端代码
    ```shell
    protoc --proto_path=. --go_out=paths=source_relative:. --natsrpc_out=paths=source_relative:. *.proto
    ```
4. Server端实现接口并创建服务
   ```
   type HelloSvc struct {
   }

   func (s *HelloSvc) Hello(ctx context.Context, req *example.HelloRequest) (*example.HelloReply, error) {
       return &example.HelloReply{
           Message: "hello " + req.Name,
       }, nil
   }

   func main() {
       conn, err := nats.Connect(*nats_url)
       defer conn.Close()

       server, err := natsrpc.NewServer(conn)
       defer server.Close(context.Background())

       svc, err := example.RegisterGreetingNRServer(server, &HelloSvc{})
       defer svc.Close()
       
       select{
       }
   }

   ```
   
5. Client 调用 rpc
   ```
   client:=natsrpc.NewClient(conn)
   
   cli := example.NewGreeterNRClient(client)
   rsp,err:=cli.Hello(context.Background(), &example.HelloRequest{Name: "natsrpc"})
   ```
 
## Examples
[here](./example)

## Bench Tool
1. 请求 `go run ./example/tool/request_bench -url=nats://127.0.0.1:4222`
2. 广播 `go run ./example/tool/publish_bench -url=nats://127.0.0.1:4222`

## TODO
- [x] service 定义文件改成gRPC标准
- [x] 支持返回错误
- [x] 支持Header
- [x] 生成Client接口
- [x] 支持中间件
- [x] 默认多线程，同时支持单一个线程
- [ ] 支持goroutine池
- [ ] 支持字节池
