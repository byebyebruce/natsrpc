package natsrpc

import (
	"strings"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

// joinSubject 组合字符串成subject
func joinSubject(s ...string) string {
	// 跳过前导空字符串
	start := 0
	for start < len(s) && s[start] == "" {
		start++
	}
	s = s[start:]

	switch len(s) {
	case 0:
		return ""
	case 1:
		return s[0]
	case 2:
		if s[1] == "" {
			return s[0]
		}
		return s[0] + "." + s[1]
	}

	// 3个或更多元素使用 bufPool
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
			bf.WriteByte('.')
		}
		bf.WriteString(v)
	}
	return bf.String()
}
