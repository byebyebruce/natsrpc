package natsrpc

import (
	"reflect"
)

type ServiceDesc struct {
	ServiceName string
	Methods     []MethodDesc
	Metadata    string
}

func (s ServiceDesc) PublishMethods() []MethodDesc {
	var ret []MethodDesc
	for _, v := range s.Methods {
		if v.IsPublish {
			ret = append(ret, v)
		}
	}
	return ret
}

type MethodDesc struct {
	MethodName  string
	Handler     Handler
	IsPublish   bool
	RequestType reflect.Type
}

func (md MethodDesc) NewRequest() any {
	return reflect.New(md.RequestType).Interface()
}
