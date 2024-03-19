package sshx

import "time"

const (
	AuthTypePass = iota + 1 // 密码
	AuthTypePublicKey
)

type (
	Options struct {
		Host     string
		Port     uint16
		User     string
		Auth     string
		AuthType int
		TimeOut  time.Duration
	}
	Option func(o *Options)
)

func NewOptions(opts ...Option) *Options {
	options := &Options{
		TimeOut: time.Second, // 默认1秒超时
	}
	for _, o := range opts {
		o(options)
	}
	return options
}

func Host(host string) Option {
	return func(o *Options) {
		o.Host = host
	}
}

func Port(port uint16) Option {
	return func(o *Options) {
		o.Port = port
	}
}

func User(user string) Option {
	return func(o *Options) {
		o.User = user
	}
}

func Auth(auth string) Option {
	return func(o *Options) {
		o.Auth = auth
	}
}

func AuthType(tp int) Option {
	return func(o *Options) {
		o.AuthType = tp
	}
}

func TimeOut(time time.Duration) Option {
	return func(o *Options) {
		o.TimeOut = time
	}
}
