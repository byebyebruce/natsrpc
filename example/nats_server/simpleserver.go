package nats_server

import (
	"fmt"
	"time"

	"github.com/nats-io/nats-server/v2/server"
)

// Run run a simple nats server
func Run(opts *server.Options) *server.Server {
	var defaultTestOptions = server.Options{
		Host:                  "127.0.0.1",
		Port:                  4222,
		NoLog:                 false,
		NoSigs:                true,
		MaxControlLine:        4096,
		DisableShortFirstPing: true,
		//Trace : true,
		//Debug: true,
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
