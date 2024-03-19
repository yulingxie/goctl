package session

type Options struct {
	clientAddr string
	serverAddr string
	protocol   string
	filter     []string
	heartbeat  bool
}

type Option func(o *Options)

func NewOptions(options ...Option) Options {
	opts := Options{
		heartbeat: false,
	}
	for _, opt := range options {
		opt(&opts)
	}
	return opts
}

func WithClientAddr(addr string) Option {
	return func(o *Options) {
		o.clientAddr = addr
	}
}

func WithServerAddr(addr string) Option {
	return func(o *Options) {
		o.serverAddr = addr
	}
}

func WithProtocol(proto string) Option {
	return func(o *Options) {
		o.protocol = proto
	}
}

func WithFilter(filter []string) Option {
	return func(o *Options) {
		o.filter = filter
	}
}

func WithHeartbeat(h bool) Option {
	return func(o *Options) {
		o.heartbeat = h
	}
}
