package natsrpc

import (
	"context"
	"sync"

	"github.com/nats-io/nats.go"
)

// Reply 用手动回复消息. 当用户要延迟返回结果时，
// 可以在当前handle函数 return nil, ErrReplyLater. 然后在其他地方调用Reply函数
//
// 例如：
//
//	func XXHandle(ctx context.Context, req *XXReq) (*XXRep, error) {
//		go func() {
//			time.Sleep(time.Second)
//			Reply(ctx, &XXRep{}, nil)
//		}
//		return nil, ErrReplyLater
//	}
func Reply(ctx context.Context, rep interface{}, repErr error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

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

// MakeReplyFunc 构造一个延迟返回函数
func MakeReplyFunc[T any](ctx context.Context) (replay func(T, error) error) {
	once := sync.Once{}
	replay = func(rep T, errRep error) error {
		var err error
		once.Do(func() {
			err = Reply(ctx, errRep, err)
		})
		return err
	}
	return
}
