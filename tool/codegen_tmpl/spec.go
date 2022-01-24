package codegen_tmpl

// FileSpec 文件描述
type FileSpec struct {
	PackageName   string
	GoPackageName string
	ServiceList   []ServiceSpec
}

// ServiceSpec 服务描述
type ServiceSpec struct {
	ServiceName  string
	Comment      string
	MethodList   []ServiceMethodSpec
	ServiceAsync bool // service 异步handler
	ClientAsync  bool // client 异步handler
}

// ServiceMethodSpec 方法描述
type ServiceMethodSpec struct {
	MethodName     string
	Comment        string
	InputTypeName  string
	OutputTypeName string
	Publish        bool // false表示request(需要返回值)，true表示广播(不需要返回值)
}
