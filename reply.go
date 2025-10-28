package natsrpc

import (
	"context"
	"errors"
	"sync"

	"github.com/go-kratos/kratos/v2/transport"
)

var ErrNotTransport = errors.New("not natsrpc transport")

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
func Reply(ctx context.Context, rep any, repErr error) error {
	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		return ErrNotTransport
	}
	st := tr.(*Transport)
	return st.replyFunc(rep, repErr)
}

// MakeReplyFunc 构造一个延迟返回函数
func MakeReplyFunc(ctx context.Context) (replay func(any, error) error) {
	once := sync.Once{}
	replay = func(rep any, errRep error) error {
		var err error
		once.Do(func() {
			err = Reply(ctx, rep, errRep)
		})
		return err
	}
	return
}
