package main

import (
	"github.com/byebyebruce/natsrpc/cmd/protoc-gen-nrpc/generator"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

func main() {
	//file, err := os.Open("/Users/liuwangchen/work/go/me/natsrpc/cmd/protoc-gen-nrpc/input.txt")
	//if err != nil {
	//	panic(456)
	//}
	//pgs.Init(pgs.ProtocInput(file)).
	//	RegisterModule(generator.New()).
	//	Render()
	pgs.Init(pgs.DebugEnv("DEBUG")).
		RegisterModule(generator.New()).
		RegisterPostProcessor(pgsgo.GoFmt()).
		Render()
}
