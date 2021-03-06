package main

import (
	"context"
	"flag"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/testdata"
)

var (
	url       = flag.String("url", "nats://127.0.0.1:4222", "nats url")
	sn        = flag.Int("s", 0, "natsURL count,0:cpu num")
	cn        = flag.Int("c", 1, "client count,0:cpu num")
	totalTime = flag.Int("t", 10, "total time")
)

var n int32

type BenchNotifyService struct {
}

func (a *BenchNotifyService) Notify(ctx context.Context, req *testdata.Empty) {
	atomic.AddInt32(&n, 1)
}

func main() {
	flag.Parse()
	if 0 == *sn {
		*sn = runtime.NumCPU()
	}
	if 0 == *cn {
		*cn = runtime.NumCPU()
	}

	var serviceName = fmt.Sprintf("%d", time.Now().UnixNano())

	op := []natsrpc.ServiceOption{natsrpc.WithServiceNoGroup()}

	for i := 0; i < *sn; i++ {
		enc, err := natsrpc.NewPBEnc(*url)
		natsrpc.IfNotNilPanic(err)
		defer enc.Close()

		server, err := natsrpc.NewServer(enc)
		natsrpc.IfNotNilPanic(err)
		defer server.Close(context.Background())

		_, err = server.Register(serviceName, &BenchNotifyService{}, op...)
		if nil != err {
			panic(err)
		}
	}

	var totalFailed uint32
	var totalSuccess uint32

	fmt.Println("start...")
	wg := sync.WaitGroup{}
	req := &testdata.Empty{}

	for i := 0; i <= *cn; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			enc, err := natsrpc.NewPBEnc(*url)
			natsrpc.IfNotNilPanic(err)
			defer enc.Close()

			client, err := natsrpc.NewClient(enc, serviceName)
			natsrpc.IfNotNilPanic(err)
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*totalTime)*time.Second)
			defer cancel()
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}
				if err := client.Publish("Notify", req); nil != err {
					atomic.AddUint32(&totalFailed, 1)
					continue
				}
				atomic.AddUint32(&totalSuccess, 1)
			}

		}(i)
	}

	wg.Wait()
	fmt.Println("elapse:", *totalTime, "suber", *sn, "pub", totalSuccess, "pub failed", totalFailed, "receive", n, "/", uint32(*sn)*totalSuccess)
}
