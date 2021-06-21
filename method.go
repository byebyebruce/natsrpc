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
func (s *MyService)Request(ctx context.Context,req *proto.Request)(*proto.Reply, error)`)
)

type fn func(ctx context.Context, data []byte) (interface{}, error)

// method 方法
type method struct {
	handle fn     // handler
	name   string // func name
}

// parseMethod 解析方法
func parseMethod(i interface{}) (map[string]*method, error) {
	ret := make(map[string]*method)
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
			ret[pM.name] = pM
		}
	}
	return ret, nil
}

// genMethod 生成方法
func genMethod(val reflect.Value, m reflect.Method) (*method, error) {
	const (
		paraNum = 3
		retNum  = 2
	)
	var (
		ctxType  reflect.Type
		reqType  reflect.Type
		respType reflect.Type
	)
	mType := m.Type

	// 检查参数
	numArgs := mType.NumIn()
	if numArgs != paraNum {
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

	// 檢查返回值
	if mType.NumOut() > 0 {
		if mType.NumOut() != retNum {
			return nil, errorFuncType
		}

		// 第一个返回值必须是pb类型
		respType = mType.Out(0)
		if respType.Kind() != reflect.Ptr {
			return nil, errorFuncType
		}
		if _, ok := reflect.New(respType.Elem()).Interface().(proto.Message); !ok {
			return nil, errorFuncType
		}

		if mType.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
			return nil, errorFuncType
		}
	}

	f := m.Func

	h := func(ctx context.Context, data []byte) (interface{}, error) {
		ctxVal := reflect.ValueOf(ctx)
		reqVal := reflect.New(reqType.Elem())
		reqPB := reqVal.Interface().(proto.Message)
		if err := proto.Unmarshal(data, reqPB); nil != err {
			return nil, err
		}

		repVal := f.Call([]reflect.Value{val, ctxVal, reqVal})
		if nil == respType {
			return nil, nil
		} else {
			var err error
			if errInter := repVal[1].Interface(); errInter != nil {
				err = errInter.(error)
			}
			return repVal[0].Interface(), err
		}
	}
	return &method{handle: h, name: m.Name}, nil
}
