package main

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
)

// fileSpec 文件描述
type fileSpec struct {
	Services []*serviceDesc
}

// serviceDesc 服务描述
type serviceDesc struct {
	ServiceType string // Greeter
	ServiceName string // helloworld.Greeter
	Metadata    string // api/helloworld/helloworld.proto
	Comment     string
	Methods     []*methodDesc
}

// methodDesc 方法描述
type methodDesc struct {
	Name         string
	OriginalName string // The parsed original name
	Comment      string
	Request      string
	Reply        string
	Publish      bool // false表示request(需要返回值)，true表示广播(不需要返回值)
}

//go:embed tmpl.gohtml
var serviceTmpl string

func (f *fileSpec) execute() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("tmpl").Parse(strings.TrimSpace(serviceTmpl))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	return strings.Trim(buf.String(), "\r\n")
}
