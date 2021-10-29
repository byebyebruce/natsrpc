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
	"github.com/nats-io/nats.go"
)

var (
	natsURL   = flag.String("url", "nats://127.0.0.1:4222", "nats url")
	sn        = flag.Int("s", 0, "natsURL count,0:cpu num")
	cn        = flag.Int("c", 0, "client count,0:cpu num")
	totalTime = flag.Int("t", 10, "total time")
)

var n int32

type BenchNotifyService struct {
}

func (a *BenchNotifyService) Notify(ctx context.Context, req *natsrpc.Empty) {
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

	var serviceName = fmt.Sprintf("å%d", time.Now().UnixNano())

	op := []natsrpc.Option{natsrpc.WithNamespace("bench_pub")}

	for i := 0; i < *sn; i++ {
		server, err := natsrpc.NewPBServer(*natsURL, nats.Name(fmt.Sprintf("bench_pub_server_%d", i)))
		if nil != err {
			panic(err)
		}
		defer server.Close(time.Second)
		_, err = server.Register(serviceName, &BenchNotifyService{}, op...)
		if nil != err {
			panic(err)
		}
	}

	var totalReq uint32
	var totalSuccess uint32

	fmt.Println("start...")
	wg := sync.WaitGroup{}
	req := &natsrpc.Empty{}
	for i := 0; i <= *cn; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			client, err := natsrpc.NewPBClient(*natsURL)
			if nil != err {
				panic(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*totalTime)*time.Second)
			defer cancel()
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}
				atomic.AddUint32(&totalReq, 1)
				if err := client.Publish("Notify", req); nil != err {
					continue
				}
				atomic.AddUint32(&totalSuccess, 1)
			}

		}(i)
	}

	wg.Wait()
	fmt.Println("elapse:", *totalTime, "suber", *sn, "pub", totalReq, "success", totalSuccess, "process", n)
}
