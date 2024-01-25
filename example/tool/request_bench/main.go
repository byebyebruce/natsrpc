package main

import (
	"context"
	"flag"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/example"
	"github.com/nats-io/nats.go"
)

var (
	natsURL   = flag.String("url", "nats://127.0.0.1:4222", "nats server")
	cn        = flag.Int("c", 64, "client number")
	totalTime = flag.Int("t", 10, "total time")
)

type BenchService struct {
	total int32
}

func (a *BenchService) Hello(ctx context.Context, req *example.HelloRequest) (*example.HelloReply, error) {
	atomic.AddInt32(&a.total, 1)
	return &example.HelloReply{}, nil
}

func main() {
	flag.Parse()

	conn, err := nats.Connect(*natsURL)
	if err != nil {
		panic(err)
	}

	server, err := natsrpc.NewServer(conn)
	if nil != err {
		panic(err)
	}
	defer server.Close(context.Background())

	bs := &BenchService{}
	svc, err := example.RegisterGreetingNRServer(server, bs)
	example.IfNotNilPanic(err)
	defer svc.Close()

	var totalSuccess uint32
	var totalFailed uint32

	fmt.Println("start...")
	wg := sync.WaitGroup{}
	req := &example.HelloRequest{}
	for i := 0; i <= *cn; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			connClient, err := nats.Connect(*natsURL)
			if err != nil {
				panic(err)
			}
			client := natsrpc.NewClient(connClient)
			greetingNRClient := example.NewGreetingNRClient(client)
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
				if _, err := greetingNRClient.Hello(context.Background(), req); nil != err {
					atomic.AddUint32(&totalFailed, 1)
					continue
				}
				atomic.AddUint32(&totalSuccess, 1)
			}

		}(i)
	}

	wg.Wait()
	fmt.Println("elapse:", *totalTime,
		"qps", totalSuccess/uint32(*totalTime),
		"req", totalSuccess+totalFailed,
		"success", totalSuccess,
		"reply", bs.total,
		"failed", totalFailed)
}
