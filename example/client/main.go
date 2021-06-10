package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/golang/protobuf/proto"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/example"
	helloworld "github.com/byebyebruce/natsrpc/testdata"
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
	conn, err := natsrpc.NewNATSConn(cfg, nats.Name("example_client"+*id))
	if nil != err {
		panic(err)
	}
	defer conn.Close()

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
	client, err := natsrpc.NewClient(conn, &example.ExampleService{}, opt...)
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

				req := &helloworld.HelloRequest{
					Name: fmt.Sprintf("hello %d", next),
				}

				reply := &helloworld.HelloReply{}
				if *singleThread {
					wg.Add(1)
					client.ID(*id).AsyncRequest(req, reply, func(message proto.Message, err error) {
						defer wg.Done()
						fmt.Println("begin AsyncRequest", reply.Message)
						if nil != err {
							panic(err)
						}
						if reply.Message != req.Name {
							panic("reply.Message")
						}
						fmt.Println("end AsyncRequest", reply.Message)
					})

				} else {
					if err := client.ID(*id).Request(req, reply); nil != err {
						panic(err)
					}
					if reply.Message != req.Name {
						panic("reply.Message")
					}
					fmt.Println("reply", reply.Message)
				}

			}
		}()
	}
	wg.Wait()
}
