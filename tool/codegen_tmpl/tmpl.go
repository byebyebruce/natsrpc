package codegen_tmpl

const serviceTmpl = `{{- range .ServiceList}}
{{$serviceName := .ServiceName}}
{{$serviceAsync := .ServiceAsync}}
{{$serviceWrapperName := print .ServiceName "Wrapper"}}
{{$clientAsync := .ClientAsync}}
{{$clientInterface := print .ServiceName "Client"}}
{{$clientWrapperName := print "_" .ServiceName "Client"}}

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
	ss := &{{ $serviceWrapperName }}{
		doer: doer,
		s:    s,
	}
	return server.Register("{{ $.GoPackageName }}.{{ $serviceName }}", ss, opts...)
}

// {{ $serviceWrapperName }} DO NOT USE
type {{ $serviceWrapperName }} struct {
	doer natsrpc.AsyncDoer
	s    {{ $serviceName }}
}
{{- range .MethodList }}
// {{ .MethodName }} DO NOT USE
	{{- if eq .Publish true }}
		func (s *{{ $serviceWrapperName }}){{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }}) {
			s.doer.Do(ctx, func() {
				s.s.{{ .MethodName }}(ctx , req)
			})
		}
	{{- else }}
		func (s *{{ $serviceWrapperName }}){{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }})(rep *{{ .OutputTypeName }}, err error) {
			done := make(chan struct{})
			cb := func(r *{{ .OutputTypeName }}, e error) {
				select {
				case <-ctx.Done():
					return
				default:
					rep, err = r, e
					select {
					case done <- struct{}{}:
					default:
					}
				}
			}

			s.doer.Do(ctx, func() {
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



// {{ $clientInterface }}
type {{ $clientInterface }} interface {
{{- range .MethodList }}
// {{ .MethodName }}
	{{- if eq .Publish false }}
		{{- if eq $clientAsync true }}
			{{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }}, cb func(*{{ .OutputTypeName }}, error), opt ...natsrpc.CallOption)
		{{- else }}
			{{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }}, opt ...natsrpc.CallOption)(*{{ .OutputTypeName }}, error)
		{{- end }}
	{{- else }}
		{{ .MethodName }}(notify *{{ .InputTypeName }}) error
	{{- end }}
{{- end }}
}

{{- if eq $clientAsync true }}
type {{ $clientWrapperName }} struct {
	c *natsrpc.Client
	doer natsrpc.AsyncDoer
}

// New{{ $clientInterface }}
func New{{ $clientInterface }}(enc *nats.EncodedConn,doer natsrpc.AsyncDoer, opts ...natsrpc.ClientOption) ({{ $clientInterface }}, error) {
	c, err := natsrpc.NewClient(enc, "{{ $.GoPackageName }}.{{ $serviceName }}", opts...)
	if err != nil {
		return nil, err
	}
	ret := &{{ $clientWrapperName }}{
		c:c,
		doer: doer,
	}
	return ret, nil
}
{{- else }}
type {{ $clientWrapperName }} struct {
	c *natsrpc.Client
}

// New{{ $clientInterface }}
func New{{ $clientInterface }}(enc *nats.EncodedConn, opts ...natsrpc.ClientOption) ({{ $clientInterface }}, error) {
	c, err := natsrpc.NewClient(enc, "{{ $.GoPackageName }}.{{ $serviceName }}", opts...)
	if err != nil {
		return nil, err
	}
	ret := &{{ $clientWrapperName }}{
		c:c,
	}
	return ret, nil
}
{{- end }}

{{- range .MethodList }}
	{{- if eq .Publish false }}
		{{- if eq $clientAsync true }}
			func (c *{{ $clientWrapperName }}) {{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }}, cb func(*{{ .OutputTypeName }}, error), opt ...natsrpc.CallOption) {
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
			func (c *{{ $clientWrapperName }}) {{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }}, opt ...natsrpc.CallOption)(*{{ .OutputTypeName }}, error) {
				rep := &{{ .OutputTypeName }}{}
				err := c.c.Request(ctx, "{{ .MethodName }}", req, rep, opt...)
				return rep, err 
			}
		{{- end }}
	{{- else }}
		func (c *{{ $clientWrapperName }}) {{ .MethodName }}(notify *{{ .InputTypeName }}) error {
			return c.c.Publish("{{ .MethodName }}", notify)
		}
	{{- end }}
{{- end }}

{{- end }}
`
