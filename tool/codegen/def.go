package codegen

import "github.com/byebyebruce/natsrpc"

type FileSpec struct {
	PackageName   string
	GoPackageName string
	ServiceList   []ServiceSpec
}

type ServiceSpec struct {
	ServiceName  string
	MethodList   []ServiceMethodSpec
	ServiceAsync bool // service 异步handler
	ClientAsync  bool // client 异步handler
}

type ServiceMethodSpec struct {
	MethodName     string
	InputTypeName  string
	OutputTypeName string
	MethodType     natsrpc.MethodType
}
