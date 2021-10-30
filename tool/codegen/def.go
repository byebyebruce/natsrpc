package codegen

import "github.com/byebyebruce/natsrpc"

type FileSpec struct {
	PackageName   string
	GoPackageName string
	ServiceList   []ServiceSpec
}

type ServiceSpec struct {
	ServiceName string
	MethodList  []ServiceMethodSpec
}

type ServiceMethodSpec struct {
	MethodName     string
	InputTypeName  string
	OutputTypeName string
	MethodType     natsrpc.MethodType
}
