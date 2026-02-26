package natsrpc

import (
	"github.com/nats-io/nats.go"
)

const (
	headerMethod = "_ns_method" // method
	headerUser   = "_ns_user"   // user header
	headerError  = "_ns_error"  // reply error
)

func encodeHeader(method string, header Header) nats.Header {
	ret := make(nats.Header, len(header)+1)
	for k, v := range header {
		ret[k] = v
	}
	ret[headerMethod] = []string{method}
	return ret
}

func decodeHeader(h nats.Header) (method string, header Header, err error) {
	methods := h[headerMethod]
	if len(methods) == 0 {
		return "", nil, ErrHeaderFormat
	}
	method = methods[0]
	header = Header(h)
	return method, header, nil
}

func makeErrorHeader(err error) nats.Header {
	if err != nil {
		return map[string][]string{headerError: {err.Error()}}
	}
	return nil
}

func getErrorHeader(h nats.Header) string {
	if h == nil {
		return ""
	}
	val := h[headerError]
	if len(val) == 0 {
		return ""
	}
	return val[0]
}
