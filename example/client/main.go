package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/testdata"
	"github.com/byebyebruce/natsrpc/testdata/pb"
	"github.com/nats-io/nats.go"
)

var (
	server       = flag.String("server", "nats://127.0.0.1:4222", "nats server")
	namespace    = flag.String("ns", "testsapce", "namespace")
	group        = flag.String("group", "", "subscribe group")
	id           = flag.String("id", "", "service id")
	count        = flag.Int("c", 10, "request count")
	thread       = flag.Int("t", 0, "thread count")
	singleThread = flag.Bool("st", false, "single thread handle")
)

func main() {
	flag.Parse()

	cfg := natsrpc.Config{
		Server: *server,
	}
	enc, err := natsrpc.NewNATSConn(cfg, nats.Name("example_client"+*id))
	if nil != err {
		panic(err)
	}
	defer enc.Close()

	opt := []natsrpc.Option{natsrpc.WithNamespace(*namespace),
		natsrpc.WithGroup(*group),
		//natsrpc.WithID(*id),
		natsrpc.WithTimeout(time.Second)}

	if *singleThread {
		singleThreadChan := make(chan func())
		opt = append(opt, natsrpc.WithSingleThreadCallback(singleThreadChan))
		go func() {
			for f := range singleThreadChan {
				f()
			}
		}()
	}
	client, err := testdata.NewGreeterClient(enc, opt...)
	if nil != err {
		panic(err)
	}
	var c int32
	wg := sync.WaitGroup{}

	if 0 == *thread {
		*thread = runtime.NumCPU()
	}
	for i := 0; i < *thread; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				next := atomic.AddInt32(&c, 1)
				if next > int32(*count) {
					break
				}

				req := &pb.HelloRequest{
					Name: fmt.Sprintf("hello %d", next),
				}

				if *singleThread {
					wg.Add(1)
					client.ID(*id).AsyncRequestAreYouOK(req, func(reply *pb.HelloReply, err error) {
						defer wg.Done()
						fmt.Println("begin AsyncRequest", reply.Message)
						if nil != err {
							panic(err)
						}
						fmt.Println("end AsyncRequest", reply.Message)
					})

				} else {
					if reply, err := client.ID(*id).RequestAreYouOK(nil, req); nil != err {
						panic(err)
					} else {
						fmt.Println("reply", reply.Message)
					}
				}

			}
		}()
	}
	wg.Wait()
}
