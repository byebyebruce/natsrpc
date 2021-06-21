package natsrpc

import (
	"context"
	"errors"
	"reflect"
	"sync"

	"github.com/golang/protobuf/proto"
)

var (
	errorFuncType = errors.New(`method must be a function likes:
func (s *MyService)Notify(ctx context.Context,req *proto.Request)
func (s *MyService)Request(ctx context.Context,req *proto.Request)(*proto.Reply, error))
func (s *MyService)Request(ctx context.Context,req *proto.Request, resp *proto.Reply, done func())`)
)

type request struct {
	reqVal reflect.Value
	over   chan struct{}
	reply  interface{}
	err    error
	once   sync.Once
}

func (s *request) done(reply interface{}, err error) {
	s.once.Do(func() {
		s.err = err
		s.reply = reply
		close(s.over)
	})
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
func parseMethod(typ reflect.Type) (map[string]*method, error) {
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
		ctxType  reflect.Type
		reqType  reflect.Type
		respType reflect.Type
	)
	mType := m.Type

	// 检查参数
	numArgs := mType.NumIn()
	numRets := mType.NumOut()

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

	// 第一个返回值必须是pb类型
	mt := methodType_None
	// 检查参数
	switch {
	case numArgs == paraNum: // notify
		if numRets == 0 {
			mt = methodType_Publish
		} else if mType.NumOut() == 2 { // request
			if numArgs > paraNum {
				return nil, errorFuncType
			}
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
			mt = methodType_Request
		} else {
			return nil, errorFuncType
		}
	case numArgs == paraNum+1: // async reply
		if numRets > 0 {
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
		case methodType_Publish, methodType_Request:
			repVal := f.Call([]reflect.Value{val, ctxVal, req.reqVal})
			if methodType_Publish == mt {
				req.done(nil, nil)
			} else {
				var err error
				if errInter := repVal[1].Interface(); errInter != nil {
					err = errInter.(error)
				}
				req.done(repVal[0].Interface(), err)
			}
		case methodType_AsyncRequest:
			cbVal := func(reply interface{}, err error) {
				req.done(reply, err)
			}
			f.Call([]reflect.Value{val, ctxVal, req.reqVal, reflect.ValueOf(cbVal)})
		}
	}
	ret := &method{
		name:    m.Name,
		reqType: reqType,
		handle:  h,
	}
	return ret, nil
}
