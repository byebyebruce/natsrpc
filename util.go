package natsrpc

import (
	"bytes"
	"go/ast"
	"reflect"
	"sync"

	"github.com/golang/protobuf/proto"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer([]byte{})
	},
}

// CombineSubject 组合字符串成subject
func CombineSubject(prefix string, s ...string) string {
	if len(s) == 0 {
		return prefix
	}
	bf := bufPool.Get().(*bytes.Buffer)
	bf.WriteString(prefix)
	for _, v := range s {
		if v == "" {
			continue
		}
		bf.WriteString(".")
		bf.WriteString(v)
	}
	subject := bf.String()
	bf.Reset()
	bufPool.Put(bf)
	return subject
}

// isExportedOrBuiltinType 是导出或内置类型
func isExportedOrBuiltinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}

// IsProtoPtrType 是否是proto指针类型
func IsProtoPtrType(t reflect.Type) bool {
	if t.Kind() != reflect.Ptr {
		return false
	}
	_, ok := reflect.New(t.Elem()).Interface().(proto.Message)
	return ok
}

// IsErrorType 是否是error类型
func IsErrorType(t reflect.Type) bool {
	return t == reflect.TypeOf((*error)(nil)).Elem()
}

// IsContextType 是否是context类型
func IsContextType(t reflect.Type) bool {
	if t.Kind() != reflect.Interface {
		return false
	}
	if t.String() != "context.Context" {
		return false
	}
	return true
}
