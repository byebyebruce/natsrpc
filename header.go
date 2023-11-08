package natsrpc

import (
	"context"

	"github.com/nats-io/nats.go"
)

type headerKey struct{}

// setHeader 填充Header
func setHeader(ctx context.Context, header map[string]string) context.Context {
	newCtx := context.WithValue(ctx, headerKey{}, header)
	return newCtx
}

// GetHeader 获得Header
func GetHeader(ctx context.Context) map[string]string {
	if ctx == nil {
		return nil
	}
	val := ctx.Value(headerKey{})
	if val != nil {
		return val.(map[string]string)
	}
	return nil
}

func encodeHeader(method string, header map[string]string) (nats.Header, error) {
	ret := map[string][]string{headerMethod: {method}}
	if len(header) > 0 {
		ret[headerUser] = make([]string, 0, len(header)*2)
		for k, v := range header {
			ret[headerUser] = append(ret[headerUser], k, v)
		}
	}
	return ret, nil
}

func decodeHeader(h nats.Header) (method string, header map[string]string, err error) {
	val := h[headerMethod]
	if len(val) == 0 {
		return "", nil, ErrHeaderFormat
	}
	method = val[0]

	if kv := h[headerUser]; len(kv) > 0 && len(kv)%2 == 0 {
		header = make(map[string]string)
		for i := 0; i < len(kv); i += 2 {
			header[kv[i]] = kv[i+1]
		}
	}
	return
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
