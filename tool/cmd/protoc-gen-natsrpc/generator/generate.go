package generator

import (
	"fmt"
	"path"
	"strings"

	"github.com/byebyebruce/natsrpc"
	"github.com/byebyebruce/natsrpc/tool/codegen"
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
	return "natsrpc"
}

func (m *NatsRpcModule) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	tmp := codegen.Template()
	// 遍历所有文件
	for _, f := range targets {
		m.Push(f.Name().String())

		// 解析出file语法树
		fileSpec := m.ExtraFile(f)
		base := strings.Split(f.Name().String(), ".")[0]
		m.OverwriteCustomTemplateFile(path.Join(m.OutputPath(), fmt.Sprintf("%s.pb.natsrpc.go", base)), tmp, fileSpec, 0644)

		m.Pop()
	}

	return m.Artifacts()
}

func (m *NatsRpcModule) ExtraFile(f pgs.File) codegen.FileSpec {
	serviceSpecs := make([]codegen.ServiceSpec, 0, len(f.Services()))
	// 遍历所有service
	for _, service := range f.Services() {
		// 通过pb的语法书解析出结构
		serviceSpecs = append(serviceSpecs, m.ExtraService(service))
	}
	return codegen.FileSpec{
		PackageName:   f.Package().ProtoName().String(),
		GoPackageName: f.Descriptor().GetOptions().GetGoPackage(),
		ServiceList:   serviceSpecs,
	}
}

// 根据语法树提取出结构
func (m *NatsRpcModule) ExtraService(service pgs.Service) codegen.ServiceSpec {
	serviceData := codegen.ServiceSpec{
		ServiceName: service.Name().String(),
	}
	sa, ca := false, false
	svcOpts := service.Descriptor().GetOptions()
	svcDescs, _ := proto.ExtensionDescs(svcOpts)
	for _, desc := range svcDescs {
		// 找到对应field
		if desc.Field == natsrpc.E_ServiceAsync.Field {
			ext, _ := proto.GetExtension(svcOpts, desc)
			// 解析出methodType
			if value, ok := ext.(*bool); ok {
				sa = *value
			}
		} else if desc.Field == natsrpc.E_ClientAsync.Field {
			ext, _ := proto.GetExtension(svcOpts, desc)
			// 解析出methodType
			if value, ok := ext.(*bool); ok {
				ca = *value
			}
		}
	}
	serviceData.ServiceAsync = sa
	serviceData.ClientAsync = ca

	for _, method := range service.Methods() {
		methodSpec := codegen.ServiceMethodSpec{
			MethodName:     method.Name().String(),
			InputTypeName:  method.Input().Name().String(),
			OutputTypeName: method.Output().Name().String(),
		}
		// 获取method的option
		opts := method.Descriptor().GetOptions()
		descs, _ := proto.ExtensionDescs(opts)
		for _, desc := range descs {
			// 找到对应field
			if desc.Field == natsrpc.E_Publish.Field {
				ext, _ := proto.GetExtension(opts, desc)
				// 解析出methodType
				if value, ok := ext.(*bool); ok {
					methodSpec.Publish = *value
					break
				}
			}
		}
		serviceData.MethodList = append(serviceData.MethodList, methodSpec)
	}
	return serviceData
}
