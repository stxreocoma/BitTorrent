package dialer

import (
	"context"
	"net"
	"time"
)

type DialerType int8

const (
	Default = iota
	HTTPProxy
	SOCKS5Proxy
)

type Config struct {
	Timeout  time.Duration
	Type     DialerType
	ProxyURL string
}

type DialFunc func(ctx context.Context, address string) (net.Conn, error)

func New(cfg Config) DialFunc {
	return func(ctx context.Context, address string) (net.Conn, error) {
		return net.DialTimeout("tcp", address, cfg.Timeout)
	}
}

func NewCustomDialer(dialFunc DialFunc) DialFunc {
	return dialFunc
}
