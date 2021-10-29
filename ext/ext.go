package ext

// WithSingleThreadCallback 服务单线程处理
func WithSingleThreadCallback(singleThreadCbChan chan func()) Option {
	return func(options *Options) {
		if nil == singleThreadCbChan {
			panic("singleThreadCbChan is nil")
		}
		options.singleThreadCbChan = singleThreadCbChan
	}
}
