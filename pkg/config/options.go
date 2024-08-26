package config

import (
	"context"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type AppOption struct {
	DisableCLIParse  bool
	IsStaked         bool
	Bootnodes        []string
	RandomIdentity   bool
	Services         []func(ctx context.Context, node host.Host)
	ProtocolHandlers map[protocol.ID]network.StreamHandler
}

type Option func(*AppOption)

var DisableCLIParse = func(o *AppOption) {
	o.DisableCLIParse = true
}

var EnableStaked = func(o *AppOption) {
	o.IsStaked = true
}

var EnableRandomIdentity = func(o *AppOption) {
	o.RandomIdentity = true
}

func (a *AppOption) Apply(opts ...Option) {
	for _, opt := range opts {
		opt(a)
	}
}

func WithBootNodes(bootnodes ...string) Option {
	return func(o *AppOption) {
		o.Bootnodes = append(o.Bootnodes, bootnodes...)
	}
}

func WithService(plugins ...func(ctx context.Context, node host.Host)) Option {
	return func(o *AppOption) {
		o.Services = append(o.Services, plugins...)
	}
}

func WithProtocolHandler(pid protocol.ID, n network.StreamHandler) Option {
	return func(o *AppOption) {
		if o.ProtocolHandlers == nil {
			o.ProtocolHandlers = make(map[protocol.ID]network.StreamHandler)
		}
		o.ProtocolHandlers[pid] = n
	}
}
