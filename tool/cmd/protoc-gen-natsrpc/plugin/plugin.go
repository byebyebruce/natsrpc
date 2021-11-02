package plugin

import (
	"strings"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/tool/codegen"
	"github.com/byebyebruce/natsrpc/tool/generator"
	"github.com/golang/protobuf/proto"
)

// MyPlugin is an implementation of the Go protocol buffer compiler's
// plugin architecture.  It generates bindings for MyPlugin support.
type MyPlugin struct {
	gen *generator.Generator
}

// Name returns the name of this plugin, "MyPlugin".
func (g *MyPlugin) Name() string {
	return "myplugin"
}

// The names for packages imported in the generated code.
// They may vary from the final path component of the import path
// if the name is used by other packages.
var (
	contextPkg string
)

// Init initializes the plugin.
func (g *MyPlugin) Init(gen *generator.Generator) {
	g.gen = gen
}

// Given a type name defined in a .proto, return its object.
// Also record that we're using it, to guarantee the associated import.
func (g *MyPlugin) objectNamed(name string) generator.Object {
	g.gen.RecordTypeUse(name)
	return g.gen.ObjectNamed(name)
}

// Given a type name defined in a .proto, return its name as we will print it.
func (g *MyPlugin) typeName(str string) string {
	return g.gen.TypeName(g.objectNamed(str))
}

// P forwards to g.gen.P.
func (g *MyPlugin) P(args ...interface{}) { g.gen.P(args...) }

// Generate generates code for the services in the given file.
func (g *MyPlugin) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}

	contextPkg = string(g.gen.AddImport("context"))
	g.gen.AddImport("github.com/byebyebruce/natsrpc")
	g.gen.AddImport("github.com/nats-io/nats.go")

	goPkg := strings.Replace(strings.Split(file.GetOptions().GetGoPackage(), ";")[0], "/", ".", -1)
	f := codegen.FileSpec{
		GoPackageName: goPkg,
	}

	for _, service := range file.FileDescriptorProto.Service {
		s := codegen.ServiceSpec{}
		s.ServiceName = service.GetName()
		if v, err := proto.GetExtension(service.GetOptions(), natsrpc.E_ServiceAsync); err == nil {
			s.ServiceAsync = *(v.(*bool))
		}
		if v, err := proto.GetExtension(service.GetOptions(), natsrpc.E_ClientAsync); err == nil {
			s.ClientAsync = *(v.(*bool))
		}
		for _, m := range service.Method {
			ms := codegen.ServiceMethodSpec{}
			if v, err := proto.GetExtension(m.GetOptions(), natsrpc.E_Publish); err == nil {
				ms.Publish = *(v.(*bool))
			}
			ms.MethodName = m.GetName()
			ms.InputTypeName = g.typeName(m.GetInputType())
			ms.OutputTypeName = g.typeName(m.GetOutputType())
			s.MethodList = append(s.MethodList, ms)
		}
		f.ServiceList = append(f.ServiceList, s)
	}
	b, err := codegen.GenText(codegen.ServiceTemplate(), f)
	if err != nil {
		g.gen.Error(err)
	}
	g.P(string(b))
}

// GenerateImports generates the import declaration for this file.
func (g *MyPlugin) GenerateImports(file *generator.FileDescriptor) {
}
