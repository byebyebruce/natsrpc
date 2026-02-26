package gogopb

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
)

type Encoder struct {
}

func (s Encoder) Decode(b []byte, i interface{}) error {
	msg, ok := i.(proto.Message)
	if !ok {
		return fmt.Errorf("gogopb: decode target is not proto.Message, got %T", i)
	}
	return proto.Unmarshal(b, msg)
}

func (s Encoder) Encode(i interface{}) ([]byte, error) {
	msg, ok := i.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("gogopb: encode target is not proto.Message, got %T", i)
	}
	return proto.Marshal(msg)
}
