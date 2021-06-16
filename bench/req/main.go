package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/byebyebruce/natsrpc/testdata/pb"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/byebyebruce/natsrpc"
)

var (
	server    = flag.String("server", "nats://127.0.0.1:4222", "nats server")
	sn        = flag.Int("s", 0, "server count,0:cpu num")
	cn        = flag.Int("c", 0, "client count,0:cpu num")
	totalTime = flag.Int("t", 30, "total time")
)

type BenchReqService struct {
}

func (a *BenchReqService) Func2(ctx context.Context, req *pb.HelloRequest, repl *pb.HelloReply) {
	repl.Message = req.Name
}

func main() {
	flag.Parse()
	if 0 == *sn {
		*sn = runtime.NumCPU()
	}
	if 0 == *cn {
		*cn = runtime.NumCPU()
	}

	cfg := natsrpc.Config{
		Server: *server,
	}

	op := []natsrpc.Option{natsrpc.WithNamespace("bench_req"),
		natsrpc.WithGroup("mygroup")}

	for i := 0; i < *sn; i++ {
		server, err := natsrpc.NewNatsRPCWithConfig(cfg, nats.Name(fmt.Sprintf("bench_req_server_%d", i)))
		if nil != err {
			panic(err)
		}
		defer server.Close()
		_, err = server.Register(&BenchReqService{}, op...)
		if nil != err {
			panic(err)
		}
	}

	var totalReq uint32
	var totalSuccess uint32
	var totalFailed uint32

	fmt.Println("start...")
	wg := sync.WaitGroup{}
	req := &pb.HelloRequest{}
	for i := 0; i <= *cn; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			client, err := natsrpc.NewClientWithConfig(cfg, fmt.Sprintf("bench_req_client_%d", idx), &BenchReqService{}, op...)
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
				if err := client.Request(req, &pb.HelloReply{}); nil != err {
					atomic.AddUint32(&totalFailed, 1)
					continue
				}
				atomic.AddUint32(&totalSuccess, 1)
			}

		}(i)
	}

	wg.Wait()
	fmt.Println("elapse:", *totalTime, "qps", totalReq/uint32(*totalTime), "req", totalReq, "success", totalSuccess, "failed", totalFailed)
}
