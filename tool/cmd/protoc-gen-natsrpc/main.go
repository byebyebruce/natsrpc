package main

import (
	"github.com/byebyebruce/natsrpc/tool/codegen_plugin"
	protoc_gen_base "github.com/byebyebruce/natsrpc/tool/protoc-gen-base"
)

func main() {
	protoc_gen_base.Main("natsrpc", &codegen_plugin.MyPlugin{})
}
