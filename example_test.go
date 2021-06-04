package xnats

import (
	"fmt"
	"testing"

	
)

func init() {
	s := Server{}
	fmt.Println(s)
}

type a struct {}
func (a a)(testdata.Test){

}
func Test_parseMethod(t *testing.T) {

}

/*
import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/proto/testdata"
	"war/log4go"
	"reflect"
	"sync"
	"testing"
	"time"
)

func init() {
	log4go.Close()
}

func createTestDispatcher() *Client {
	cfg := Config{
		Server:        "nats://172.25.156.5:4222",
		ReconnectWait: 5,
		MaxReconnects: 1000,
	}
	dis, err := NewClient(&cfg, "Test_Dispatcher", 10)
	if nil != err {
		panic(err)
	}
	return dis
}

func TestNatsDispatcher_Notify(t *testing.T) {

	postfix := "fdasfdas"

	cb1 := func(pb *testdata.InnerMessage, reply string, err string) {
		fmt.Println("1 pb", pb, "reply", reply, "err", err)
	}
	dis := createTestDispatcher()
	defer dis.Close(false)
	dis.RegisterHandler(cb1, false, postfix)

	cb2 := func(pb *testdata.InnerMessage, reply string, err string) {
		fmt.Println("2 pb", pb, "reply", reply, "err", err)
	}
	dis1 := createTestDispatcher()
	defer dis1.Close(false)
	dis1.RegisterHandler(cb2, false, postfix)

	dis3 := createTestDispatcher()
	defer dis3.Close(false)
	dis3.Notify(&testdata.InnerMessage{
		Host: proto.String("127.0.0.1"),
	}, postfix)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for {
			m := <-dis.MsgChan()
			dis.Process(m)
			wg.Done()
		}
	}()
	go func() {
		for {
			m2 := <-dis1.MsgChan()
			dis1.Process(m2)
			wg.Done()
		}
	}()
	wg.Wait()
}

func TestNatsDispatcher_GroupNotify(t *testing.T) {
	wg := sync.WaitGroup{}

	postfix := "fdasfdas"

	cb1 := func(pb *testdata.InnerMessage, reply string, err string) {
		fmt.Println("1 pb", pb, "reply", reply, "err", err)
		wg.Done()
	}
	dis := createTestDispatcher()
	defer dis.Close(false)
	dis.RegisterHandler(cb1, true, postfix)
	go func() {
		for {
			m := <-dis.MsgChan()
			dis.Process(m)
		}

	}()
	cb2 := func(pb *testdata.InnerMessage, reply string, err string) {
		fmt.Println("2 pb", pb, "reply", reply, "err", err)
		wg.Done()
	}
	dis1 := createTestDispatcher()
	defer dis1.Close(false)
	dis1.RegisterHandler(cb2, true, postfix)
	go func() {
		for {
			m := <-dis1.MsgChan()
			dis1.Process(m)
		}
	}()

	wg.Add(1)
	dis.Notify(&testdata.InnerMessage{
		Host: proto.String("127.0.0.1"),
	}, postfix)

	wg.Wait()
}

func TestNatsDispatcher_Request(t *testing.T) {
	dis := createTestDispatcher()
	defer dis.Close(false)

	dis1 := createTestDispatcher()
	defer dis1.Close(false)

	postfix := "fdasfdas"

	h := func(pb *testdata.InnerMessage, reply string, err string) {
		fmt.Println("req pb", pb, "reply", reply, "err", err)
		dis.Reply(reply, &testdata.OtherMessage{Key: proto.Int64(1002)})
	}
	dis.RegisterHandler(h, true, postfix)
	go func() {
		for {
			m := <-dis.MsgChan()
			dis.Process(m)
		}
	}()

	wg := sync.WaitGroup{}
	cb := func(pb *testdata.OtherMessage, reply string, err string) {
		fmt.Println("resp pb=", pb, "reply=", reply, "err=", err)
		wg.Done()
	}
	wg.Add(1)
	dis1.Request(&testdata.InnerMessage{
		Host: proto.String("127.0.0.2"),
	}, cb, postfix)
	go func() {
		for {
			m := <-dis1.MsgChan()
			dis1.Process(m)
		}
	}()

	wg.Wait()
}

func TestNatsDispatcher_RegisterHandler(t *testing.T) {
	dis := createTestDispatcher()
	defer dis.Close(false)
	dis1 := createTestDispatcher()
	defer dis1.Close(false)
	wg := sync.WaitGroup{}
	dis.RegisterHandler(func(pb *testdata.InnerMessage, reply string, error string) {
		fmt.Println("dis", error)
		wg.Done()
	}, false)
	dis1.RegisterHandler(func(pb *testdata.InnerMessage, reply string, error string) {
		fmt.Println("dis1", error)
		wg.Done()
	}, false)
	for i := 0; i < 10; i++ {
		dis.Notify(&testdata.InnerMessage{
			Host: proto.String("127.0.0.1"),
		})
		wg.Add(2)
	}

	go func() {
		for m := range dis.MsgChan() {
			dis.Process(m)
		}
	}()
	go func() {
		for m := range dis1.MsgChan() {
			dis1.Process(m)
		}
	}()
	wg.Wait()
}

func TestNatsDispatcher_RegisterSyncHandler(t *testing.T) {
	dis := createTestDispatcher()
	defer dis.Close(false)
	wg := sync.WaitGroup{}
	dis.RegisterSyncHandler(func(pb *testdata.InnerMessage, reply string, error string) {
		fmt.Println("dis", error)
		wg.Done()
	}, false)

	for i := 0; i < 10; i++ {
		dis.Notify(&testdata.InnerMessage{
			Host: proto.String("127.0.0.1"),
		})
		wg.Add(1)
	}

	go func() {
		for m := range dis.MsgChan() {
			dis.Process(m)
		}
	}()
	wg.Wait()
}

func TestNatsDispatcher_RegisterGroupedSyncHandler(t *testing.T) {
	dis := createTestDispatcher()
	defer dis.Close(false)

	wg := sync.WaitGroup{}
	dis.RegisterSyncHandler(func(pb *testdata.InnerMessage, reply string, error string) {
		fmt.Println("dis", error)
		wg.Done()
	}, true)
	go func() {
		for m := range dis.MsgChan() {
			dis.Process(m)
		}
	}()

	dis1 := createTestDispatcher()
	defer dis1.Close(false)
	dis1.RegisterSyncHandler(func(pb *testdata.InnerMessage, reply string, error string) {
		fmt.Println("dis1", error)
		wg.Done()
	}, true)
	go func() {
		for m := range dis1.MsgChan() {
			dis1.Process(m)
		}
	}()

	for i := 0; i < 10; i++ {
		dis.Notify(&testdata.InnerMessage{
			Host: proto.String("127.0.0.1"),
		})
		wg.Add(1)
	}

	wg.Wait()
}

func TestNatsDispatcher_RegisterGroupedSyncHandlerMultiGo(t *testing.T) {
	dis := createTestDispatcher()
	defer dis.Close(false)

	wg := sync.WaitGroup{}
	//for i := 0; i < 10; i++ {

	dis.RegisterSyncHandler(func(pb *testdata.InnerMessage, reply string, error string) {
		go func() {
			fmt.Println("handle", *pb.Port)
			n := 10 - int(*pb.Port)
			if n < 0 {
				n = 0
			}
			time.Sleep(time.Second * time.Duration(n))
			ret := &testdata.InnerMessage{Host: proto.String(""), Port: pb.Port}
			dis.Reply(reply, ret)
			wg.Done()
		}()

	}, true)
	//}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			fmt.Println("req", idx)
			ret := &testdata.InnerMessage{}
			e := dis.RequestSync(&testdata.InnerMessage{
				Host: proto.String(""),
				Port: proto.Int(idx),
			}, ret)
			fmt.Println("resp", ret.GetPort(), e)
		}(i)

	}

	wg.Wait()
}

func BenchmarkDispatcher_Notify(b *testing.B) {

	dis := createTestDispatcher()
	defer dis.Close(false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dis.Notify(&testdata.InnerMessage{
			Host: proto.String("127.0.0.1"),
		}, 1000410001)
	}
}

func BenchmarkDispatcher_Request(b *testing.B) {

	dis := createTestDispatcher()
	defer dis.Close(false)
	h := func(pb *testdata.InnerMessage, reply string, err string) {
		//fmt.Println("req pb", pb, "reply", reply, "err", err)
		dis.Reply(reply, &testdata.OtherMessage{Key: proto.Int64(1002)})
	}
	dis.RegisterHandler(h, true)
	go func() {
		for m := range dis.MsgChan() {
			dis.Process(m)
		}
	}()

	dis2 := createTestDispatcher()
	defer dis2.Close(false)
	cb := func(pb *testdata.OtherMessage, reply string, err string) {
		//fmt.Println("resp pb=", pb, "reply=", reply, "err=", err)
	}
	go func() {
		for m := range dis2.MsgChan() {
			dis2.Process(m)
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dis.Request(&testdata.InnerMessage{
			Host: proto.String("127.0.0.1"),
		}, cb)
	}
}

func BenchmarkDispatcher_Process(b *testing.B) {
	dis := createTestDispatcher()

	cb := func(pb *testdata.OtherMessage, reply string, err string) {
	}
	m := &msg{
		handler: reflect.ValueOf(cb),
		arg:     reflect.ValueOf(&testdata.OtherMessage{}),
		reply:   valEmptyString,
		err:     valEmptyString,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dis.Process(m)
	}
}
*/
