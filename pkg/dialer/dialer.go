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
	switch cfg.Type {
	case HTTPProxy:
		return BuildHTTPProxyDialer(cfg)
	case SOCKS5Proxy:
		return BuildSOCKS5ProxyDialer(cfg)
	default:
		return buildDefaultDialer(cfg)
	}
}

func NewCustomDialer(dialFunc DialFunc) DialFunc {
	return dialFunc
}

func buildDefaultDialer(cfg Config) DialFunc {
	return func(ctx context.Context, address string) (net.Conn, error) {
		return net.DialTimeout("tcp", address, cfg.Timeout)
	}
}

func BuildHTTPProxyDialer(cfg Config) DialFunc {
	panic("HTTP proxy support not implemented yet!")
}

func BuildSOCKS5ProxyDialer(cfg Config) DialFunc {
	panic("SOCKS proxy support not implemented yet!")
}
