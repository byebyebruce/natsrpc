package main

import (
	"github.com/byebyebruce/natsrpc/tool/cmd/protoc-gen-natsrpc/plugin"
	"github.com/byebyebruce/natsrpc/tool/generator"
)

func main() {
	generator.Main("natsrpc", &plugin.MyPlugin{})
}
