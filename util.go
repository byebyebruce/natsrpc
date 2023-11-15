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
	if len(s) == 0 {
		return ""
	} else if len(s) == 1 {
		return s[0]
	} else if len(s) == 2 {
		if s[0] == "" {
			return s[1]
		} else if s[1] == "" {
			return s[0]
		}
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

	return bf.String()
}
