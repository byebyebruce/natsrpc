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
	"github.com/byebyebruce/natsrpc/example"
	"github.com/nats-io/nats.go"
)

var (
	url       = flag.String("url", "nats://127.0.0.1:4222", "nats url")
	sn        = flag.Int("s", 0, "natsURL count,0:cpu num")
	cn        = flag.Int("c", 10, "client count,0:cpu num")
	totalTime = flag.Int("t", 10, "total time")
)

var n int32

type BenchService struct {
}

func (a *BenchService) HelloToAll(ctx context.Context, req *example.HelloRequest) (*example.Empty, error) {
	atomic.AddInt32(&n, 1)
	return nil, nil
}

func main() {
	flag.Parse()
	if 0 == *sn {
		*sn = runtime.NumCPU()
	}
	if 0 == *cn {
		*cn = runtime.NumCPU()
	}

	for i := 0; i < *sn; i++ {
		conn, err := nats.Connect(*url)
		if err != nil {
			panic(err)
		}
		example.IfNotNilPanic(err)
		defer conn.Close()

		server, err := natsrpc.NewServer(conn)
		example.IfNotNilPanic(err)
		defer server.Close(context.Background())

		_, err = example.RegisterGreetingToAllNATSRPCServer(server, &BenchService{})
		example.IfNotNilPanic(err)
	}

	var totalFailed uint32
	var totalSuccess uint32

	fmt.Println("start...")
	wg := sync.WaitGroup{}

	for i := 0; i <= *cn; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			conn, err := nats.Connect(*url)
			if err != nil {
				panic(err)
			}
			example.IfNotNilPanic(err)
			defer conn.Close()

			client := example.NewGreetingToAllNATSRPCClient(conn)
			example.IfNotNilPanic(err)
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*totalTime)*time.Second)
			defer cancel()

			req := &example.HelloRequest{}

			for {
				select {
				case <-ctx.Done():
					return
				default:
				}
				if err := client.HelloToAll(req); nil != err {
					atomic.AddUint32(&totalFailed, 1)
					continue
				}
				atomic.AddUint32(&totalSuccess, 1)
			}

		}(i)
	}

	wg.Wait()
	fmt.Println("elapse:", *totalTime,
		"suber", *sn,
		"pub", totalSuccess,
		"pub failed", totalFailed,
		"receive", n)
}
