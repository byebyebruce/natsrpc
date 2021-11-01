package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/byebyebruce/natsrpc/tool/cmd/protoc-gen-natsrpc/generator"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

func main() {
	if out, ok := os.LookupEnv("OUTPUT"); ok {
		s, _ := ioutil.ReadAll(os.Stdin)
		fmt.Println(ioutil.WriteFile(out, s, os.ModePerm))
		os.Exit(0)
	}

	opt := []pgs.InitOption{pgs.DebugEnv("DEBUG")}

	if in, ok := os.LookupEnv("INPUT"); ok {
		f, err := os.Open(in)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		opt = append(opt, pgs.ProtocInput(f))
	}

	pgs.Init(opt...).
		RegisterModule(generator.New()).
		RegisterPostProcessor(pgsgo.GoFmt()).
		Render()
}
