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
func (s *MyService)Request(ctx context.Context,req *proto.Request, resp *proto.Reply, done func())`)
)

type request struct {
	reqVal reflect.Value
	over   chan struct{}
	reply  interface{}
	err    error
}

func (s *request) done() {
	close(s.over)
}

type fn func(context.Context, reflect.Value, *request)

type methodType int

const (
	methodType_None         methodType = iota // none
	methodType_Publish                        // publish
	methodType_Request                        // request
	methodType_AsyncRequest                   // async request
)

// method 方法
type method struct {
	mt      methodType   // 方法类型
	handle  fn           // handler
	name    string       // func name
	reqType reflect.Type // request type
}

// 构造一个 request
func (m *method) newRequest(b []byte) (*request, error) {
	reqVal := reflect.New(m.reqType.Elem())
	if len(b) > 0 {
		pb := reqVal.Interface().(proto.Message)
		if err := proto.Unmarshal(b, pb); nil != err {
			return nil, err
		}
	}
	req := &request{
		over:   make(chan struct{}),
		reqVal: reqVal,
	}
	return req, nil
}

// parseMethod 解析方法
func parseMethod(typ reflect.Type) ([]*method, error) {
	var ret []*method

	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)

		if !isExportedOrBuiltinType(m.Type) {
			continue
		}

		if pM, err := genMethod(m); nil != err {
			return ret, err
		} else {
			ret = append(ret, pM)
		}
	}
	return ret, nil
}

// genMethod 生成方法
func genMethod(m reflect.Method) (*method, error) {
	const paraNum = 3 // ptr, ctx, req
	var (
		ctxType  reflect.Type
		reqType  reflect.Type
		respType reflect.Type
	)
	mType := m.Type
	numArgs := mType.NumIn()

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

	mt := methodType_None
	// 检查参数
	switch {
	case numArgs == paraNum: // notify
		mt = methodType_None
	case mType.NumOut() == 2: // request
		mt = methodType_Request
	case numArgs == paraNum+1: // async reply
		// 如果有第3个参数说明是请求
		respType = mType.In(3)
		if respType.Kind() != reflect.Ptr {
			return nil, errorFuncType
		}
		if _, ok := reflect.New(respType.Elem()).Interface().(proto.Message); !ok {
			return nil, errorFuncType
		}
		mt = methodType_AsyncRequest
	default:
		return nil, errorFuncType
	}

	f := m.Func

	h := func(ctx context.Context, val reflect.Value, req *request) {
		ctxVal := reflect.ValueOf(ctx)

		switch mt {
		case methodType_Publish:
			f.Call([]reflect.Value{val, ctxVal, req.reqVal})
			req.done()
		case methodType_Request:
			respVal := reflect.New(respType.Elem())
			f.Call([]reflect.Value{val, ctxVal, req.reqVal, respVal})
			req.done()
		case methodType_AsyncRequest:
			respVal := reflect.New(respType.Elem())
			cbVal := reflect.ValueOf(func() {
				req.done()
			})
			f.Call([]reflect.Value{val, ctxVal, req.reqVal, respVal, cbVal})
		}
	}
	ret := &method{
		name:    m.Name,
		reqType: reqType,
		handle:  h,
	}
	return ret, nil
}
