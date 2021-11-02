//go:generate protoc --proto_path=. --example_out=paths=source_relative:. pb/example.proto
package main

import (
	"github.com/byebyebruce/natsrpc/tool/generator"
	"github.com/byebyebruce/natsrpc/tool/generator/protoc-gen-example/myplugin"
)

func main() {
	generator.Main("example", &myplugin.MyPlugin{})
}
