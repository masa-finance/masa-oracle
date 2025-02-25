package node

import (
	"context"

	"github.com/masa-finance/masa-oracle/node/types"
	"github.com/masa-finance/masa-oracle/pkg/masacrypto"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type NodeOption struct {
	IsStaked    bool
	UDP         bool
	TCP         bool
	IsValidator bool
	PortNbr     int

	IsTwitterScraper  bool
	IsDiscordScraper  bool
	IsTelegramScraper bool
	IsWebScraper      bool

	Bootnodes            []string
	RandomIdentity       bool
	Services             []func(ctx context.Context, node *OracleNode)
	PubSubHandles        []PubSubHandlers
	ProtocolHandlers     map[protocol.ID]network.StreamHandler
	MasaProtocolHandlers map[string]network.StreamHandler
	Environment          string
	Version              string
	MasaDir              string
	CachePath            string

	OracleProtocol       string
	NodeDataSyncProtocol string
	NodeGossipTopic      string
	Rendezvous           string
	WorkerProtocol       string
	PageSize             int

	KeyManager *masacrypto.KeyManager
}

type PubSubHandlers struct {
	ProtocolName string
	Handler      types.SubscriptionHandler
	IncludeSelf  bool
}

type Option func(*NodeOption)

var EnableStaked = func(o *NodeOption) {
	o.IsStaked = true
}

var EnableRandomIdentity = func(o *NodeOption) {
	o.RandomIdentity = true
}

var EnableTCP = func(o *NodeOption) {
	o.TCP = true
}

var EnableUDP = func(o *NodeOption) {
	o.UDP = true
}

var IsValidator = func(o *NodeOption) {
	o.IsValidator = true
}

var IsTwitterScraper = func(o *NodeOption) {
	o.IsTwitterScraper = true
}

var IsDiscordScraper = func(o *NodeOption) {
	o.IsDiscordScraper = true
}

var IsTelegramScraper = func(o *NodeOption) {
	o.IsTelegramScraper = true
}

var IsWebScraper = func(o *NodeOption) {
	o.IsWebScraper = true
}

func (a *NodeOption) Apply(opts ...Option) {
	for _, opt := range opts {
		opt(a)
	}
}

// HasBootnodes checks if the AppConfig has any bootnodes configured.
// It returns true if there is at least one bootnode in the Bootnodes slice and it is not an empty string.
// Otherwise, it returns false, indicating that no bootnodes are configured.
func (a *NodeOption) HasBootnodes() bool {
	if len(a.Bootnodes) == 0 {
		return false
	}

	return a.Bootnodes[0] != ""
}

func WithBootNodes(bootnodes ...string) Option {
	return func(o *NodeOption) {
		o.Bootnodes = append(o.Bootnodes, bootnodes...)
	}
}

func WithService(plugins ...func(ctx context.Context, node *OracleNode)) Option {
	return func(o *NodeOption) {
		o.Services = append(o.Services, plugins...)
	}
}

func WithProtocolHandler(pid protocol.ID, n network.StreamHandler) Option {
	return func(o *NodeOption) {
		if o.ProtocolHandlers == nil {
			o.ProtocolHandlers = make(map[protocol.ID]network.StreamHandler)
		}
		o.ProtocolHandlers[pid] = n
	}
}

func WithEnvironment(env string) Option {
	return func(o *NodeOption) {
		o.Environment = env
	}
}

func WithVersion(version string) Option {
	return func(o *NodeOption) {
		o.Version = version
	}
}

func WithMasaProtocolHandler(pid string, n network.StreamHandler) Option {
	return func(o *NodeOption) {
		if o.MasaProtocolHandlers == nil {
			o.MasaProtocolHandlers = make(map[string]network.StreamHandler)
		}
		o.MasaProtocolHandlers[pid] = n
	}
}

func WithPubSubHandler(protocolName string, handler types.SubscriptionHandler, includeSelf bool) Option {
	return func(o *NodeOption) {
		o.PubSubHandles = append(o.PubSubHandles, PubSubHandlers{protocolName, handler, includeSelf})
	}
}

func WithPort(port int) Option {
	return func(o *NodeOption) {
		o.PortNbr = port
	}
}

func WithMasaDir(directory string) Option {
	return func(o *NodeOption) {
		o.MasaDir = directory
	}
}

func WithCachePath(path string) Option {
	return func(o *NodeOption) {
		o.CachePath = path
	}
}

func WithKeyManager(km *masacrypto.KeyManager) Option {
	return func(o *NodeOption) {
		o.KeyManager = km
	}
}

func WithOracleProtocol(s string) Option {
	return func(o *NodeOption) {
		o.OracleProtocol = s
	}
}

func WithNodeDataSyncProtocol(s string) Option {
	return func(o *NodeOption) {
		o.NodeDataSyncProtocol = s
	}
}

func WithNodeGossipTopic(s string) Option {
	return func(o *NodeOption) {
		o.NodeGossipTopic = s
	}
}

func WithRendezvous(s string) Option {
	return func(o *NodeOption) {
		o.Rendezvous = s
	}
}

func WithWorkerProtocol(s string) Option {
	return func(o *NodeOption) {
		o.WorkerProtocol = s
	}
}

func WithPageSize(size int) Option {
	return func(o *NodeOption) {
		o.PageSize = size
	}
}
