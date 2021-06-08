package natsrpc

type serviceOptions struct {
	group     string
	namespace string
	id        string
	timeout   int64
}

func newDefaultOption() serviceOptions {
	return serviceOptions{
		namespace: "default",
		id:        "",
		timeout:   3,
	}
}

// ServiceOption ServiceOption
type ServiceOption func(options *serviceOptions)

// WithGrouped
func WithGroup(group string) ServiceOption {
	return func(options *serviceOptions) {
		options.group = group
	}
}

// WithNamespace
func WithNamespace(namespace string) ServiceOption {
	return func(options *serviceOptions) {
		options.namespace = namespace
	}
}

// WithID
func WithID(id string) ServiceOption {
	return func(options *serviceOptions) {
		options.id = id
	}
}

// WithTimeout
func WithTimeout(timeout int64) ServiceOption {
	return func(options *serviceOptions) {
		options.timeout = timeout
	}
}
