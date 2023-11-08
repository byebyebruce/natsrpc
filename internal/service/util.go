package service

import (
	"go/ast"
	"reflect"

	"google.golang.org/protobuf/proto"
)

// isExportedOrBuiltinType 是导出或内置类型
func isExportedOrBuiltinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}

// isProtoPtrType 是否是proto指针类型
func isProtoPtrType(t reflect.Type) bool {
	if t.Kind() != reflect.Ptr {
		return false
	}
	_, ok := reflect.New(t.Elem()).Interface().(proto.Message)
	return ok
}

// isErrorType 是否是error类型
func isErrorType(t reflect.Type) bool {
	return t == reflect.TypeOf((*error)(nil)).Elem()
}

// isContextType 是否是context类型
func isContextType(t reflect.Type) bool {
	if t.Kind() != reflect.Interface {
		return false
	}
	if t.String() != "context.Context" {
		return false
	}
	return true
}
