package monitor

import (
	"time"
)

type Options struct {
	snapshotLen int32
	promiscuous bool
	timeout     time.Duration
	protocol    string
	serverAddr  string
	heartbeat   bool
	packets     []string
}

type Option func(*Options)

func NewOptions(options ...Option) Options {
	opts := Options{
		snapshotLen: 10240,
		promiscuous: true,
		timeout:     10,
		heartbeat:   false,
	}
	for _, opt := range options {
		opt(&opts)
	}
	return opts
}

func WithSnapshotLen(len int32) Option {
	return func(o *Options) {
		o.snapshotLen = len
	}
}

func WithPromiscuous(p bool) Option {
	return func(o *Options) {
		o.promiscuous = p
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.timeout = timeout
	}
}

func WithProtocol(proto string) Option {
	return func(o *Options) {
		o.protocol = proto
	}
}

func WithServerAddr(addr string) Option {
	return func(o *Options) {
		o.serverAddr = addr
	}
}

func WithHeartbeat(filter bool) Option {
	return func(o *Options) {
		o.heartbeat = filter
	}
}

func WithPackets(pkts []string) Option {
	return func(o *Options) {
		o.packets = pkts
	}
}
