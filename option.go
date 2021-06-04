package xnats

type serviceOptions struct {
	grouped   bool
	namespace string
	sync      bool
	id        string
}

func newDefaultOption() serviceOptions {
	return serviceOptions{
		grouped:   false,
		namespace: "default",
		sync:      false,
		id:        "0",
	}
}

// ServiceOption ServiceOption
type ServiceOption func(options *serviceOptions)

// WithGrouped
func WithGrouped(grouped bool) ServiceOption {
	return func(options *serviceOptions) {
		options.grouped = grouped
	}
}

// WithNamespace
func WithNamespace(namespace string) ServiceOption {
	return func(options *serviceOptions) {
		options.namespace = namespace
	}
}

// WithSyncHandler
func WithSyncHandler(sync bool) ServiceOption {
	return func(options *serviceOptions) {
		options.sync = sync
	}
}
