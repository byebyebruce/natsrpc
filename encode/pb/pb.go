package pb

import (
	"google.golang.org/protobuf/proto"
)

type Encoder struct {
}

func (s Encoder) Decode(b []byte, i interface{}) error {
	return proto.Unmarshal(b, i.(proto.Message))
}
func (s Encoder) Encode(i interface{}) ([]byte, error) {
	return proto.Marshal(i.(proto.Message))
}
