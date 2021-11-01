package async

import (
	"context"
	"fmt"
	"time"
)

type AsyncFunc struct {
	fc chan func()
}

func New(fc chan func()) *AsyncFunc {
	return &AsyncFunc{
		fc: fc,
	}
}

func (c *AsyncFunc) Do(ctx context.Context, f func() (interface{}, error)) (ret interface{}, err error) {
	over := make(chan struct{})
	cb := func() {
		defer close(over)
		ret, err = f()
	}
	select {
	case <-ctx.Done():
		err = ctx.Err()
		return
	case c.fc <- cb:
	}
	select {
	case <-over:
	case <-ctx.Done():
		err = ctx.Err()
	}
	return
}

func (c *AsyncFunc) DoNoneBlock(task func() (interface{}, error), timeout time.Duration, cb func(interface{}, error)) {
	go func() {
		var (
			ret interface{}
			err error
		)
		over := make(chan struct{})
		t := time.After(timeout)

		go func() {
			defer close(over)
			ret, err = task()
		}()
		select {
		case <-t:
			err = fmt.Errorf("error timeout")
			return
		case <-over:
		}
		f := func() {
			cb(ret, err)
		}
		select {
		case <-t:
			return
		case c.fc <- f:
		}
	}()
}
