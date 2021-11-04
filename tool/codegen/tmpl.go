package codegen

const serviceTmpl = `{{- range .ServiceList}}
{{$serviceName := .ServiceName}}
{{$serviceAsync := .ServiceAsync}}
{{$clientAsync := .ClientAsync}}

// {{ $serviceName }}
type {{ $serviceName }} interface {
{{- range .MethodList }}
// {{ .MethodName }}
	{{- if eq .Publish false }}
		{{- if eq $serviceAsync true }}
			{{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }}, cb func(*{{ .OutputTypeName }}, error))	
		{{- else }}
			{{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }}) (*{{ .OutputTypeName }}, error)
		{{- end }}
	{{- else }}
		{{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }})
	{{- end }}
{{- end }}
}

// Register{{ $serviceName }}
{{- if eq $serviceAsync false }}
func Register{{ $serviceName }}(server *natsrpc.Server, s {{ $serviceName }}, opts ...natsrpc.ServiceOption) (natsrpc.IService, error) {
	return server.Register("{{ $.GoPackageName }}.{{ $serviceName }}", s, opts...)
}
{{- else }}
func Register{{ $serviceName }}(server *natsrpc.Server, s {{ $serviceName }}, doer natsrpc.AsyncDoer, opts ...natsrpc.ServiceOption) (natsrpc.IService, error) {
	ss := &{{ $serviceName }}Wrapper{
		doer: doer,
		s:    s,
	}
	return server.Register("{{ $.GoPackageName }}.{{ $serviceName }}", ss, opts...)
}

// {{ $serviceName }}Wrapper DO NOT USE
type {{ $serviceName }}Wrapper struct {
	doer natsrpc.AsyncDoer
	s    {{ $serviceName }}
}
{{- range .MethodList }}
// {{ .MethodName }} DO NOT USE
	{{- if eq .Publish true }}
		func (s *{{ $serviceName }}Wrapper){{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }}) {
			s.doer.Do(ctx, func() {
				s.s.{{ .MethodName }}(ctx , req)
			})
		}
	{{- else }}
		func (s *{{ $serviceName }}Wrapper){{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }})(rep *{{ .OutputTypeName }}, err error) {
			done := make(chan struct{})
			s.doer.Do(ctx, func() {
				cb := func(r *{{ .OutputTypeName }}, e error) {
					rep, err = r, e
					select {
					case done <- struct{}{}:
					default:
					}
				}
				s.s.Hello(ctx, req, cb)
			})
			select {
			case <-ctx.Done():
				rep, err = nil, ctx.Err()
			case <-done:
			}
			return
		}
	{{- end }}
{{- end }}

{{- end }}





{{- $clientName := .ServiceName}}

// {{ $clientName }}Client
{{- if eq $clientAsync true }}
type {{ $clientName }}Client struct {
	c *natsrpc.Client
	doer natsrpc.AsyncDoer
}

// New{{$clientName}}Client
func New{{ $clientName }}Client(enc *nats.EncodedConn,doer natsrpc.AsyncDoer, opts ...natsrpc.ClientOption) (*{{ $clientName }}Client, error) {
	c, err := natsrpc.NewClient(enc, "{{ $.GoPackageName }}.{{ $clientName }}", opts...)
	if err != nil {
		return nil, err
	}
	ret := &{{ $clientName }}Client{
		c:c,
		doer: doer,
	}
	return ret, nil
}
{{- else }}
type {{ $clientName }}Client struct {
	c *natsrpc.Client
}

// New{{$clientName}}Client
func New{{ $clientName }}Client(enc *nats.EncodedConn, opts ...natsrpc.ClientOption) (*{{ $clientName }}Client, error) {
	c, err := natsrpc.NewClient(enc, "{{ $.GoPackageName }}.{{ $clientName }}", opts...)
	if err != nil {
		return nil, err
	}
	ret := &{{ $clientName }}Client{
		c:c,
	}
	return ret, nil
}
{{- end }}

{{- range .MethodList }}
// {{ .MethodName }}
	{{- if eq .Publish false }}
		{{- if eq $clientAsync true }}
			func (c *{{ $clientName }}Client) {{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }}, cb func(*{{ .OutputTypeName }}, error), opt ...natsrpc.CallOption) {
				go func() {
					rep := &{{ .OutputTypeName }}{}
					err := c.c.Request(ctx, "{{ .MethodName }}", req, rep, opt...)
					newCb := func() {
						cb(rep, err)
					}
					c.doer.Do(ctx, newCb)
				}()
			}
		{{- else }}
			func (c *{{ $clientName }}Client) {{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }}, opt ...natsrpc.CallOption)(*{{ .OutputTypeName }}, error) {
				rep := &{{ .OutputTypeName }}{}
				err := c.c.Request(ctx, "{{ .MethodName }}", req, rep, opt...)
				return rep, err 
			}
		{{- end }}
	{{- else }}
		func (c *{{ $clientName }}Client) {{ .MethodName }}(notify *{{ .InputTypeName }}) error {
			return c.c.Publish("{{ .MethodName }}", notify)
		}
	{{- end }}
{{- end }}

{{- end }}
`
