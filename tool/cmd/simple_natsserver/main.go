package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats-server/v2/server"
)

func main() {
	s := RunSimpleNatsServer(nil)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sig
	s.Shutdown()
}

// RunSimpleNatsServer run a simple nats server
func RunSimpleNatsServer(opts *server.Options) *server.Server {
	var defaultTestOptions = server.Options{
		Host:                  "127.0.0.1",
		Port:                  4222,
		NoLog:                 false,
		NoSigs:                true,
		MaxControlLine:        4096,
		DisableShortFirstPing: true,
		Trace : true,
		Debug: true,
	}

	if opts == nil {
		opts = &defaultTestOptions
	}
	// Optionally override for individual debugging of tests
	//opts.Trace = true
	s, err := server.NewServer(opts)
	if err != nil || s == nil {
		panic(fmt.Sprintf("No NATS Server object returned: %v", err))
	}

	s.ConfigureLogger()

	// Run server in Go routine.
	go s.Start()

	// Wait for accept loop(s) to be started
	if !s.ReadyForConnections(10 * time.Second) {
		panic("Unable to start NATS Server in Go Routine")
	}
	return s
}
