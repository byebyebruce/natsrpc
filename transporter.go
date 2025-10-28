package natsrpc

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/nats-io/nats.go"
)

const (
	Kind transport.Kind = "natsrpc"
)

var _ transport.Transporter = (*Transport)(nil)

// Transport is an HTTP transport.
type Transport struct {
	operation      string
	reqHeader      Header
	replyHeader    Header
	request        any
	requestSubject string
	replySubject   string
	replyFunc      func(any, error) error
}

// Kind returns the transport kind.
func (tr *Transport) Kind() transport.Kind {
	return Kind
}

// Endpoint returns the transport endpoint.
func (tr *Transport) Endpoint() string {
	return ""
}

// Operation returns the transport operation.
func (tr *Transport) Operation() string {
	return tr.operation
}

// Request returns the HTTP request.
func (tr *Transport) Request() any {
	return tr.request
}

// RequestHeader returns the request header.
func (tr *Transport) RequestHeader() transport.Header {
	return tr.reqHeader
}

// ReplyHeader returns the reply header.
func (tr *Transport) ReplyHeader() transport.Header {
	return tr.replyHeader
}

func (tr *Transport) ReplySubject() string {
	return tr.replySubject
}

// SetOperation sets the transport operation.
func SetOperation(ctx context.Context, op string) {
	if tr, ok := transport.FromServerContext(ctx); ok {
		if tr, ok := tr.(*Transport); ok {
			tr.operation = op
		}
	}
}

type Header nats.Header

// Get returns the value associated with the passed key.
func (hc Header) Get(key string) string {
	return http.Header(hc).Get(key)
}

// Set stores the key-value pair.
func (hc Header) Set(key string, value string) {
	http.Header(hc).Set(key, value)
}

// Add append value to key-values pair.
func (hc Header) Add(key string, value string) {
	http.Header(hc).Add(key, value)
}

// Keys lists the keys stored in this carrier.
func (hc Header) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range http.Header(hc) {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of values associated with the passed key.
func (hc Header) Values(key string) []string {
	return http.Header(hc).Values(key)
}
