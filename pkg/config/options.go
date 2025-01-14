package config

import (
	"github.com/masa-finance/masa-oracle/node"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/masa-finance/masa-oracle/pkg/workers"
)

var constantOptions = []node.Option{
	node.WithOracleProtocol(OracleProtocol),
	node.WithNodeDataSyncProtocol(NodeDataSyncProtocol),
	node.WithNodeGossipTopic(NodeGossipTopic),
	node.WithRendezvous(Rendezvous),
	node.WithPageSize(PageSize),
}

// WithConstantOptions adds options that are set to constant values. We need to add them to
// the node to avoid a dependency loop.
func WithConstantOptions(nodes ...node.Option) []node.Option {
	return append(nodes, constantOptions...)
}

func InitOptions(cfg *AppConfig) ([]node.Option, *workers.WorkHandlerManager, *pubsub.PublicKeySubscriptionHandler) {
	// WorkerManager configuration
	workerManagerOptions := []workers.WorkerOptionFunc{
		workers.WithMasaDir(cfg.MasaDir),
	}

	cachePath := cfg.CachePath
	if cachePath == "" {
		cachePath = cfg.MasaDir + "/cache"
	}

	masaNodeOptions := WithConstantOptions(
		node.EnableStaked,
		//	WithService(),
		node.WithEnvironment(cfg.Environment),
		node.WithVersion(cfg.Version),
		node.WithPort(cfg.PortNbr),
		node.WithBootNodes(cfg.Bootnodes...),
		node.WithMasaDir(cfg.MasaDir),
		node.WithCachePath(cachePath),
		node.WithKeyManager(cfg.KeyManager),
		node.WithWorkerProtocol(WorkerProtocol),
	)

	if cfg.TwitterScraper {
		workerManagerOptions = append(workerManagerOptions, workers.EnableTwitterWorker)
		masaNodeOptions = append(masaNodeOptions, node.IsTwitterScraper)
	}

	if cfg.TelegramScraper {
		// TODO: Telegram scraper is not implemented yet in the worker (?)
		masaNodeOptions = append(masaNodeOptions, node.IsTelegramScraper)
	}

	if cfg.DiscordScraper {
		workerManagerOptions = append(workerManagerOptions, workers.EnableDiscordScraperWorker)
		masaNodeOptions = append(masaNodeOptions, node.IsDiscordScraper)
	}

	if cfg.WebScraper {
		workerManagerOptions = append(workerManagerOptions, workers.EnableWebScraperWorker)
		masaNodeOptions = append(masaNodeOptions, node.IsWebScraper)
	}

	if cfg.TunnelEnabled {
		masaNodeOptions = append(masaNodeOptions, node.EnableTunnel)
	}

	workHandlerManager := workers.NewWorkHandlerManager(workerManagerOptions...)
	blockChainEventTracker := node.NewBlockChain()
	pubKeySub := &pubsub.PublicKeySubscriptionHandler{}

	masaNodeOptions = append(masaNodeOptions, []node.Option{
		// Register the worker manager
		node.WithMasaProtocolHandler(
			WorkerProtocol,
			workHandlerManager.HandleWorkerStream,
		),
		node.WithPubSubHandler(PublicKeyTopic, pubKeySub, false),
		node.WithPubSubHandler(BlockTopic, blockChainEventTracker, true),
	}...)

	if cfg.Validator {
		// Subscribe and if actor start monitoring actor workers
		// considering all that matters is if the node is staked
		// and other peers can do work we only need to check this here
		// if this peer can or cannot scrape or write that is checked in other places
		masaNodeOptions = append(masaNodeOptions,
			node.WithService(blockChainEventTracker.Start(cfg.MasaDir)),
		)
	}

	if cfg.UDP {
		masaNodeOptions = append(masaNodeOptions, node.EnableUDP)
	}

	if cfg.TCP {
		masaNodeOptions = append(masaNodeOptions, node.EnableTCP)
	}

	if cfg.Validator {
		masaNodeOptions = append(masaNodeOptions, node.IsValidator)
	}

	return masaNodeOptions, workHandlerManager, pubKeySub
}
