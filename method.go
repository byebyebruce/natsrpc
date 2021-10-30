package natsrpc

import (
	"context"
	"errors"
	"reflect"
)

var (
	errorFuncType = errors.New(`method must like:
func (s *MyService)Notify(ctx context.Context,req *proto.Request)
func (s *MyService)Request(ctx context.Context,req *proto.Request)(*proto.Reply, error))`)
)

type handler func(svc interface{}, ctx context.Context, req interface{}) (interface{}, error)

// method 方法
type method struct {
	handle  handler      // handler
	name    string       // func name
	reqType reflect.Type // request type
}

// 构造一个 request
func (m *method) newRequest() interface{} {
	return reflect.New(m.reqType.Elem()).Interface()
}

// parseMethod 解析方法
func parseMethod(svc interface{}) (map[string]*method, error) {
	typ := reflect.TypeOf(svc)
	ret := make(map[string]*method)
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)

		if !isExportedOrBuiltinType(m.Type) {
			continue
		}

		if pM, err := genMethod(m); nil != err {
			return ret, err
		} else {
			ret[pM.name] = pM
		}
	}
	return ret, nil
}

// genMethod 生成方法
func genMethod(m reflect.Method) (*method, error) {
	const paraNum = 3 // ptr, ctx, req
	var (
		reqType reflect.Type
	)
	mType := m.Type

	// 检查参数
	numArgs := mType.NumIn()
	numRets := mType.NumOut()

	// 第1个参数必须是context
	if !IsContextType(mType.In(1)) {
		return nil, errorFuncType
	}

	// 第2个参数是pb类型
	reqType = mType.In(2)
	if !IsProtoPtrType(reqType) {
		return nil, errorFuncType
	}

	if numArgs != paraNum {
		return nil, errorFuncType
	}

	// 第一个返回值必须是pb类型
	// 检查参数
	switch {
	case numRets == 0: // notify
	case mType.NumOut() == 2: // request
		if numArgs > paraNum {
			return nil, errorFuncType
		}
		if !IsProtoPtrType(mType.Out(0)) {
			return nil, errorFuncType
		}
		if !IsErrorType(mType.Out(1)) {
			return nil, errorFuncType
		}
	default:
		return nil, errorFuncType
	}

	h := func(svc interface{}, ctx context.Context, req interface{}) (interface{}, error) {
		repVal := m.Func.Call([]reflect.Value{reflect.ValueOf(svc), reflect.ValueOf(ctx), reflect.ValueOf(req)})
		if len(repVal) == 0 {
			return nil, nil
		} else {
			var err error
			if errInter := repVal[1].Interface(); errInter != nil {
				err = errInter.(error)
			}
			return repVal[0].Interface(), err
		}
	}

	ret := &method{
		name:    m.Name,
		reqType: reqType,
		handle:  h,
	}
	return ret, nil
}
