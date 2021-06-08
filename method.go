package natsrpc

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
)

var (
	errorFuncType = errors.New("handler must be a function like [func(ctx context.Context,req *proto.MyUser, resp *proto.MyUser)]")
)

type method struct {
	handler func(ctx context.Context, data []byte) interface{}
	replay  bool
	name    string
}

func parseStruct(i interface{}) ([]*method, error) {
	var ret []*method
	typ := reflect.TypeOf(i)
	val := reflect.ValueOf(i)
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)

		if !isExportedOrBuiltinType(m.Type) {
			continue
		}

		if pM, err := parseMethod(val, m); nil != err {
			return ret, err
		} else {
			ret = append(ret, pM)
		}
	}
	return ret, nil
}

// cb(*req,*resp)
func parseMethod(val reflect.Value, m reflect.Method) (*method, error) {
	const paraNum = 3
	var (
		ctxType  reflect.Type
		reqType  reflect.Type
		respType reflect.Type
	)
	mType := m.Type
	numArgs := mType.NumIn()

	if paraNum == numArgs {

	} else if paraNum+1 == numArgs {
		respType = mType.In(3)
		if respType.Kind() != reflect.Ptr {
			return nil, errorFuncType
		}

		if _, ok := reflect.New(respType.Elem()).Interface().(proto.Message); !ok {
			return nil, errorFuncType
		}
	} else {
		return nil, errorFuncType
	}
	ctxType = mType.In(1)
	if ctxType.Kind() != reflect.Interface {
		return nil, errorFuncType
	}
	if ctxType.Name() != "Context" {
		return nil, errorFuncType
	}

	reqType = mType.In(2)
	if reqType.Kind() != reflect.Ptr {
		return nil, errorFuncType
	}
	// check arg type
	if _, ok := reflect.New(reqType.Elem()).Interface().(proto.Message); !ok {
		return nil, errorFuncType
	}

	f := m.Func
	fmt.Println(f.String(), f)
	h := func(ctx context.Context, data []byte) interface{} {
		ctxVal := reflect.ValueOf(ctx)
		reqVal := reflect.New(reqType.Elem())
		reqPB := reqVal.Interface().(proto.Message)
		if err := proto.Unmarshal(data, reqPB); nil != err {
			return nil
			//l4g.Error("[nats(%s)] cb proto.Unmarshal error=[%s]", s.name, err.Error())
		}

		if nil == respType {
			f.Call([]reflect.Value{val, ctxVal, reqVal})
			return nil
		} else {
			respVal := reflect.New(respType.Elem())
			f.Call([]reflect.Value{val, ctxVal, reqVal, respVal})
			return respVal.Interface()
		}

		//l4g.Debug("[nats(%s)] sync callback sub=[%s] reply=[%s]", s.name, m.Subject, m.Reply)
	}
	return &method{handler: h, replay: true, name: typeName(reqType)}, nil
}
