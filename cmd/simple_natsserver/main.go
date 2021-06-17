package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/byebyebruce/natsrpc"
)

func main() {
	s := natsrpc.RunSimpleNatsServer(nil)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sig
	s.Shutdown()
}
