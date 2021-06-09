package natsrpc

import (
	"context"
	"errors"
	"reflect"

	"github.com/golang/protobuf/proto"
)

var (
	errorFuncType = errors.New(`method must be a function likes:
func (s *MyService)Notify(ctx context.Context,req *proto.Request)
func (s *MyService)Request(ctx context.Context,req *proto.Request, resp *proto.Reply)`)
)

type fn func(ctx context.Context, data []byte) (interface{}, error)

// method 方法
type method struct {
	handle fn     // handler
	name   string // func name
}

// parseMethod 解析方法
func parseMethod(i interface{}) ([]*method, error) {
	var ret []*method
	typ := reflect.TypeOf(i)
	val := reflect.ValueOf(i)
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)

		if !isExportedOrBuiltinType(m.Type) {
			continue
		}

		if pM, err := genMethod(val, m); nil != err {
			return ret, err
		} else {
			ret = append(ret, pM)
		}
	}
	return ret, nil
}

// genMethod 生成方法
func genMethod(val reflect.Value, m reflect.Method) (*method, error) {
	const paraNum = 3
	var (
		ctxType  reflect.Type
		reqType  reflect.Type
		respType reflect.Type
	)
	mType := m.Type
	numArgs := mType.NumIn()

	// 检查参数
	switch numArgs {
	case paraNum: // notify
	case paraNum + 1: // request
		// 如果有第3个参数说明是请求
		respType = mType.In(3)
		if respType.Kind() != reflect.Ptr {
			return nil, errorFuncType
		}
		if _, ok := reflect.New(respType.Elem()).Interface().(proto.Message); !ok {
			return nil, errorFuncType
		}
	default:
		return nil, errorFuncType
	}

	// 第1个参数必须是context
	ctxType = mType.In(1)
	if ctxType.Kind() != reflect.Interface {
		return nil, errorFuncType
	}
	if ctxType.Name() != "Context" {
		return nil, errorFuncType
	}

	// 第2个参数是pb类型
	reqType = mType.In(2)
	if reqType.Kind() != reflect.Ptr {
		return nil, errorFuncType
	}
	if _, ok := reflect.New(reqType.Elem()).Interface().(proto.Message); !ok {
		return nil, errorFuncType
	}

	f := m.Func

	h := func(ctx context.Context, data []byte) (interface{}, error) {
		ctxVal := reflect.ValueOf(ctx)
		reqVal := reflect.New(reqType.Elem())
		reqPB := reqVal.Interface().(proto.Message)
		if err := proto.Unmarshal(data, reqPB); nil != err {
			return nil, err
		}

		if nil == respType {
			f.Call([]reflect.Value{val, ctxVal, reqVal})
			return nil, nil
		} else {
			respVal := reflect.New(respType.Elem())
			f.Call([]reflect.Value{val, ctxVal, reqVal, respVal})
			return respVal.Interface(), nil
		}
	}
	return &method{handle: h, name: typeName(reqType)}, nil
}
