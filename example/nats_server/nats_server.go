package nats_server

import (
	"fmt"
	"time"

	"github.com/nats-io/nats-server/v2/server"
)

var DefaultTestOptions = server.Options{
	Host:                  "127.0.0.1",
	Port:                  4222,
	NoLog:                 false,
	NoSigs:                true,
	MaxControlLine:        4096,
	DisableShortFirstPing: true,
}

// RunServer run nats server
func RunServer(opts *server.Options) *server.Server {
	if opts == nil {
		opts = &DefaultTestOptions
	}
	// Optionally override for individual debugging of tests
	opts.Trace = true
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
