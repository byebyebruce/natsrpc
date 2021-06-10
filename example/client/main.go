package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/example/service"
	helloworld "github.com/byebyebruce/natsrpc/testdata"
)

var (
	server    = flag.String("server", "nats://127.0.0.1:4222", "nats server")
	namespace = flag.String("ns", "testsapce", "namespace")
	group     = flag.String("group", "", "subscribe group")
	id        = flag.String("id", "", "service id")
	count     = flag.Int("c", 10, "request count")
	thread    = flag.Int("t", 0, "thread count")
)

func main() {
	flag.Parse()

	cfg := natsrpc.Config{
		Server: *server,
	}
	conn, err := natsrpc.NewNATSConn(cfg, "example_client"+*id)
	if nil != err {
		panic(err)
	}
	defer conn.Close()

	client, err := natsrpc.NewClient(conn, &service.ExampleService{},
		natsrpc.WithNamespace(*namespace),
		natsrpc.WithGroup(*group),
		//natsrpc.WithID(*id),
		natsrpc.WithTimeout(time.Second))
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

				if err := client.ID(*id).Request(req, reply); nil != err {
					panic(err)
				}
				if reply.Message != req.Name {
					panic("reply.Message")
				}
				fmt.Println("reply", reply.Message)
			}

		}()
	}
	wg.Wait()
}
