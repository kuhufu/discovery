package discovery

import "time"

type Option func(n *Options)

type Options struct {
	ServiceName    string        //服务名
	ServiceVersion string        //服务版本
	Addrs          []string      //注册中心地址（etcd）
	DialTimeout    time.Duration //拨号超时时间
}

func newOptions(opts ...Option) Options {
	options := Options{
		DialTimeout: time.Second * 5,
	}

	for _, opt := range opts {
		opt(&options)
	}

	return options
}

func Name(name string) Option {
	return func(n *Options) {
		n.ServiceName = name
	}
}

func Version(version string) Option {
	return func(n *Options) {
		n.ServiceVersion = version
	}
}

func RegisterAddrs(addr ...string) Option {
	return func(n *Options) {
		n.Addrs = addr
	}
}

func DialTimeout(timeout time.Duration) Option {
	return func(n *Options) {
		n.DialTimeout = timeout
	}
}
