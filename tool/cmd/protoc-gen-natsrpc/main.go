package main

import (
	"gitlab.uuzu.com/sanguox/natsrpc/tool/codegen_plugin"
	protoc_gen_base "gitlab.uuzu.com/sanguox/natsrpc/tool/protoc-gen-base"
)

func main() {
	protoc_gen_base.Main("natsrpc", &codegen_plugin.MyPlugin{})
}
