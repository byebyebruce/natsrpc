package natsrpc

// ServiceDesc 服务描述
type ServiceDesc struct {
	ServiceName string       // 服务名
	Methods     []MethodDesc // 方法列表
	Metadata    string       // 元数据
}

func (s ServiceDesc) hasPublishMethod() bool {
	for _, v := range s.Methods {
		if v.IsPublish {
			return true
		}
	}
	return false
}

// MethodDesc 方法描述
type MethodDesc struct {
	MethodName string  // 方法名
	Handler    Handler // 方法处理函数
	IsPublish  bool    // 是否发布
	//RequestType reflect.Type // 请求类型
}

// NewRequest 创建请求
// 后面优化成静态的
/*
func (md MethodDesc) NewRequest() any {
	return reflect.New(md.RequestType).Interface()
}
*/
