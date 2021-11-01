package main

import (
	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/extension/simpleserver"
)

func main() {
	s := simpleserver.Run(nil)
	defer s.Shutdown()

	enc, err := natsrpc.NewPBEnc(s.ClientURL())
	natsrpc.IfNotNilPanic(err)
	defer enc.Close()

}
