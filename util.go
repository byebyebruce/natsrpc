package natsrpc

import (
	"context"
	"strings"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

// JoinSubject 组合字符串成subject
func JoinSubject(s ...string) string {
	if len(s) == 0 {
		return ""
	}
	bf := bufPool.Get().(*strings.Builder)
	defer func() {
		bf.Reset()
		bufPool.Put(bf)
	}()
	first := true
	for _, v := range s {
		if v == "" {
			continue
		}
		if first {
			first = false
		} else {
			bf.WriteString(".")
		}
		bf.WriteString(v)
	}
	subject := bf.String()

	return subject
}

func IfNotNilPanic(err error) {
	if err != nil {
		panic(err)
	}
}

type AsyncDoer interface {
	AsyncDo(context.Context, func(func(interface{}, error))) (interface{}, error)
}
