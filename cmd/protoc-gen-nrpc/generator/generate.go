package generator

import (
	"fmt"
	"path"
	"strings"
	"text/template"

	"github.com/byebyebruce/natsrpc/annotation"
	"github.com/golang/protobuf/proto"
	pgs "github.com/lyft/protoc-gen-star"
)

func New() pgs.Module {
	return &NatsRpcModule{
		ModuleBase: &pgs.ModuleBase{},
	}
}

type NatsRpcModule struct {
	*pgs.ModuleBase
}

func (m *NatsRpcModule) Name() string {
	return "nrpc"
}

func (m *NatsRpcModule) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	// 遍历所有文件
	for _, f := range targets {
		m.Push(f.Name().String())
		serviceSpecs := make([]ServiceSpec, 0, len(f.Services()))
		// 遍历所有service
		for _, service := range f.Services() {
			// 通过pb的语法书解析出结构
			serviceSpecs = append(serviceSpecs, m.ExtraService(service))
		}

		base := strings.Split(f.Name().String(), ".")[0]
		tmp, err := template.New("tem").Parse(tempService)
		if err != nil {
			fmt.Println(err)
			continue
		}
		m.OverwriteCustomTemplateFile(path.Join(m.OutputPath(), fmt.Sprintf("%s.nrpc.pb.go", base)), tmp, serviceSpecs, 0644)

		m.Pop()
	}

	return m.Artifacts()
}

// 根据语法树提取出结构
func (m *NatsRpcModule) ExtraService(service pgs.Service) ServiceSpec {
	serviceData := ServiceSpec{
		ServiceName: service.Name().String(),
	}
	for _, method := range service.Methods() {
		methodSpec := ServiceMethodSpec{
			MethodName:     method.Name().String(),
			InputTypeName:  method.Input().Name().String(),
			OutputTypeName: method.Output().Name().String(),
		}
		// 获取method的option
		opts := method.Descriptor().GetOptions()
		descs, _ := proto.ExtensionDescs(opts)
		for _, desc := range descs {
			// 找到对应field
			if desc.Field == 2360 {
				ext, _ := proto.GetExtension(opts, desc)
				// 解析出methodType
				if value, ok := ext.(*annotation.MethodType); ok {
					methodSpec.MethodType = *value
					break
				}
			}
		}
		serviceData.MethodList = append(serviceData.MethodList, methodSpec)
	}
	return serviceData
}

type ServiceSpec struct {
	ServiceName string
	MethodList  []ServiceMethodSpec
}

type ServiceMethodSpec struct {
	MethodName     string
	InputTypeName  string
	OutputTypeName string
	MethodType     annotation.MethodType
}

const tempService = `
{{- range .}}

// {{ .ServiceName }}
type {{ .ServiceName }} interface {
{{- range .MethodList -}}
	{{- if eq .MethodType 0 }}
	// {{ .MethodName }}Async
	{{ .MethodName }}Async(ctx context.Context, req *{{ .InputTypeName }}, reply func(*{{ .OutputTypeName }}, error))	
	{{- end }}

	{{- if eq .MethodType 1 }}
	// {{ .MethodName }}
	{{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }}) (*{{ .OutputTypeName }}, error)
	{{- end}}

	{{- if eq .MethodType 2 }}
	// Publish{{ .MethodName }}
	Publish{{ .MethodName }}(ctx context.Context, notify *{{ .InputTypeName }})
	{{- end}}
{{- end }}
}


// Register{{ .ServiceName }}
func Register{{ .ServiceName }}(server *natsrpc.Server, s {{ .ServiceName }}, opts ...natsrpc.Option) (natsrpc.Service, error) {
	return server.Register("xxx.{{ .ServiceName }}", s, opts...)
}

{{- $clientName := .ServiceName}}


// {{ $clientName }}Client
type {{ $clientName }}Client struct {
	c *natsrpc.Client
}

// New{{$clientName}}Client
func New{{ $clientName }}Client(enc *nats.EncodedConn, opts ...natsrpc.Option) (*{{ $clientName }}Client, error) {
	c, err := natsrpc.NewClient(enc, "xxx.{{ $clientName }}", opts...)
	if err != nil {
		return nil, err
	}
	ret := &{{ $clientName }}Client{
		c:c,
	}
	return ret, nil
}

// ID 根据ID获得client
func (c *{{ $clientName }}Client) ID(id interface{}) *{{ $clientName }}Client {
	return &{{ $clientName }}Client{
		c : c.c.ID(id),
	}
}

{{ range .MethodList}}
	{{- if eq .MethodType 0 -}}
// {{ .MethodName }}Async
func (c *{{ $clientName }}Client) {{ .MethodName }}Async(req *{{ .InputTypeName }}, cb func(*{{ .OutputTypeName }}, error)){
	rep := &{{ .OutputTypeName }}{}
	f := func(_ proto.Message, err error) {
		cb(rep, err)
	}
	c.c.AsyncRequest("{{ .MethodName }}", req, rep, f)
}
	{{- end}}

	{{- if eq .MethodType 1 -}}
// {{ .MethodName }}
func (c *{{ $clientName }}Client) {{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }}) (*{{ .OutputTypeName }}, error) {
	rep := &{{ .OutputTypeName }}{}
	err := c.c.Request(ctx, "{{ .MethodName }}", req, rep)
	return rep, err
}
	{{- end}}

	{{- if eq .MethodType 2 -}}
// Publish{{ .MethodName }}
func (c *{{ $clientName }}Client) Publish{{ .MethodName }}(notify *{{ .MethodName }}) error {
	return c.c.Publish("{{ .MethodName }}", notify)
}
	{{- end}}
{{ end }}
	
{{- end}}
`
