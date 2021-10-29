//go:generate protoc --proto_path=. --go_out=plugins=grpc:. natsrpc.proto
package natsrpc

const (
	headerError = "error"
)

type marshaller interface {
	Unmarshal(b []byte, i interface{}) error
	Marshal(i interface{}) ([]byte, error)
}
