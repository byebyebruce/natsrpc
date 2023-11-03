package natsrpc

import (
	"encoding/json"
)

func encodeHeader(method string, header map[string]string) (map[string][]string, error) {
	val := []string{method}
	ret := map[string][]string{headerNATSRPC: []string{method}}
	if header != nil {
		b, err := json.Marshal(header)
		if err != nil {
			return nil, err
		}
		val = append(val, string(b))
	}
	ret[headerNATSRPC] = val
	return ret, nil
}

func decodeHeader(h map[string][]string) (method string, header map[string]string, err error) {
	val := h[headerNATSRPC]
	if len(val) == 0 {
		return "", nil, ErrHeaderFormat
	}
	method = val[0]
	if len(val) > 1 {
		err = json.Unmarshal([]byte(val[1]), &header)
		if err != nil {
			return "", nil, err
		}
	}
	return
}

func makeErrorHeader(err error) map[string][]string {
	if err != nil {
		return map[string][]string{headerError: {err.Error()}}
	}
	return nil
}

func getErrorHeader(h map[string][]string) string {
	if h == nil {
		return ""
	}
	val := h[headerError]
	if len(val) == 0 {
		return ""
	}
	return val[0]
}
