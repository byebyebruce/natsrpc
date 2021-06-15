# NATSRPC
NATSRPC 是一个基于nats的简单rpc

## FEATURE
* 使用非常简单，不需要服务发现
* 使用protobuf定义消息类型
* 支持空间隔离
* 支持定向发送也支持负载均衡(nats内置的同组内随机)
* 不用手动定义sub(使用的req的类型名，同时问题是同service只能一个方法用这个类型。后面考虑用代码生产)
* 支持单协程回调(适用于逻辑单协程模型)

## 使用
`go get github.com/byebyebruce/natsrpc`

## 示例
* [Client](example/client/main.go)
* [Server](example/server/main.go)
* [Service](example/example_service.go)
> 运行示例需要部署gnatsd，如果没有可以临时启动`go run example/nats_server/main/main.go`