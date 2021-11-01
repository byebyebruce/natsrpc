package codegen

type FileSpec struct {
	PackageName   string
	GoPackageName string
	ServiceList   []ServiceSpec
}

type ServiceSpec struct {
	ServiceName  string
	MethodList   []ServiceMethodSpec
	ServiceAsync bool // service 异步handler
	ClientAsync  bool // client 异步handler
}

type ServiceMethodSpec struct {
	MethodName     string
	InputTypeName  string
	OutputTypeName string
	Publish        bool // false表示request(需要返回值)，true表示广播(不需要返回值)
}
