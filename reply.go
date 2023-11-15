package natsrpc

import (
	"context"
	"sync"

	"github.com/nats-io/nats.go"
)

// Reply 用手动回复消息，一般用于延迟回复
func Reply(ctx context.Context, rep interface{}, repErr error) error {
	meta := getMeta(ctx)
	if meta == nil {
		return ErrNoMeta
	}
	if meta.reply == "" {
		return ErrEmptyReply
	}

	respMsg := &nats.Msg{
		Subject: meta.reply,
		//Data:    b,
		Header: makeErrorHeader(repErr),
	}

	b, err := meta.server.opt.encoder.Encode(rep)
	if err != nil {
		return err
	}
	respMsg.Data = b
	return meta.server.conn.PublishMsg(respMsg)
}

// MakeReplyFunc 用手动回复消息，一般用于延迟回复
func MakeReplyFunc[T any](ctx context.Context) func(T, error) error {
	once := sync.Once{}
	return func(rep T, errRep error) error {
		var err error
		once.Do(func() {
			err = Reply(ctx, errRep, err)
		})
		return err
	}
}
