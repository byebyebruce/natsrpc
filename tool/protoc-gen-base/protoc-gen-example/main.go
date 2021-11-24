//go:generate protoc --proto_path=. --example_out=paths=source_relative:. pb/example.proto
package main

import (
	protoc_gen_base "github.com/byebyebruce/natsrpc/tool/protoc-gen-base"
	"github.com/byebyebruce/natsrpc/tool/protoc-gen-base/protoc-gen-example/myplugin"
)

func main() {
	protoc_gen_base.Main("example", &myplugin.MyPlugin{})
}
