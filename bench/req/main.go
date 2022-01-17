package main

import (
	"context"
	"flag"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"gitlab.uuzu.com/war/natsrpc"
)

var (
	natsURL   = flag.String("url", "nats://127.0.0.1:4222", "nats server")
	sn        = flag.Int("s", 128, "server number")
	cn        = flag.Int("c", 128, "client number")
	totalTime = flag.Int("t", 10, "total time")
)

type BenchReqService struct {
}

func (a *BenchReqService) Request(ctx context.Context, req *natsrpc.Empty) (*natsrpc.Empty, error) {
	repl := &natsrpc.Empty{}
	return repl, nil
}

func main() {
	flag.Parse()

	groupOpt := natsrpc.WithServiceGroup("mygroup")

	var serviceName = "bench"
	enc, err := natsrpc.NewPBEnc(*natsURL)
	if err != nil {
		panic(err)
	}

	for i := 0; i < *sn; i++ {
		server, err := natsrpc.NewServer(enc)
		if nil != err {
			panic(err)
		}
		defer server.Close(context.Background())
		_, err = server.Register(serviceName, &BenchReqService{}, groupOpt)
		if nil != err {
			panic(err)
		}
	}

	var totalSuccess uint32
	var totalFailed uint32

	fmt.Println("start...")
	wg := sync.WaitGroup{}
	req := &natsrpc.Empty{}
	for i := 0; i <= *cn; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			enc, err := natsrpc.NewPBEnc(*natsURL)
			if err != nil {
				panic(err)
			}
			client, err := natsrpc.NewClient(enc, serviceName)
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
				resp := &natsrpc.Empty{}
				if err := client.Request(context.Background(), "Request", req, resp); nil != err {
					atomic.AddUint32(&totalFailed, 1)
					continue
				}
				atomic.AddUint32(&totalSuccess, 1)
			}

		}(i)
	}

	wg.Wait()
	fmt.Println("elapse:", *totalTime, "qps", totalSuccess/uint32(*totalTime), "req", totalSuccess+totalFailed, "success", totalSuccess, "failed", totalFailed)
}
