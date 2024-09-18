package config

import (
	"context"

	"github.com/masa-finance/masa-oracle/node/types"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type AppOption struct {
	DisableCLIParse      bool
	IsStaked             bool
	Bootnodes            []string
	RandomIdentity       bool
	Services             []func(ctx context.Context, node host.Host)
	PubSubHandles        []PubSubHandlers
	ProtocolHandlers     map[protocol.ID]network.StreamHandler
	MasaProtocolHandlers map[string]network.StreamHandler
	Environment          string
	Version              string
}

type PubSubHandlers struct {
	ProtocolName string
	Handler      types.SubscriptionHandler
	IncludeSelf  bool
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

func WithEnvironment(env string) Option {
	return func(o *AppOption) {
		o.Environment = env
	}
}

func WithVersion(version string) Option {
	return func(o *AppOption) {
		o.Version = version
	}
}

func WithMasaProtocolHandler(pid string, n network.StreamHandler) Option {
	return func(o *AppOption) {
		if o.MasaProtocolHandlers == nil {
			o.MasaProtocolHandlers = make(map[string]network.StreamHandler)
		}
		o.MasaProtocolHandlers[pid] = n
	}
}

func WithPubSubHandler(protocolName string, handler types.SubscriptionHandler, includeSelf bool) Option {
	return func(o *AppOption) {
		o.PubSubHandles = append(o.PubSubHandles, PubSubHandlers{protocolName, handler, includeSelf})
	}
}
