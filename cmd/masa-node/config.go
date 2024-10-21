package main

import (
	"fmt"

	"github.com/masa-finance/masa-oracle/node"
	"github.com/masa-finance/masa-oracle/pkg/config"
	teekeys "github.com/masa-finance/masa-oracle/tee"

	pubsub "github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/masa-finance/masa-oracle/pkg/workers"
)

func initOptions(cfg *config.AppConfig) ([]node.Option, *workers.WorkHandlerManager, *pubsub.PublicKeySubscriptionHandler) {
	// WorkerManager configuration
	// XXX: this needs to be moved under config, but now it's here as there are import cycles given singletons
	workerManagerOptions := []workers.WorkerOptionFunc{}

	masaNodeOptions := []node.Option{
		node.EnableStaked,
		//	config.WithService(),
		node.WithEnvironment(cfg.Environment),
		node.WithVersion(cfg.Version),
		node.WithPort(cfg.PortNbr),
		node.WithBootNodes(cfg.Bootnodes...),
	}

	if cfg.TwitterScraper {
		workerManagerOptions = append(workerManagerOptions, workers.EnableTwitterWorker)
		masaNodeOptions = append(masaNodeOptions, node.IsTwitterScraper)
	}

	if cfg.RemoteAttestation {
		masaNodeOptions = append(masaNodeOptions, node.EnableRemoteAttestationChallenge)
		signer, err := teekeys.EmbeddedCertificates.ReadFile("public.pem")
		if err != nil {
			panic(fmt.Sprintf("Remote attestation enabled, but we failed to read public key (should be embedded in the binary during build time..): %v", err))
		}
		masaNodeOptions = append(masaNodeOptions, node.WithSigner(signer))
	}

	if !cfg.Debug {
		masaNodeOptions = append(masaNodeOptions, node.EnableProductionMode)
	}

	if cfg.TelegramScraper {
		// XXX: Telegram scraper is not implemented yet in the worker (?)
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

	if cfg.LlmServer {
		workerManagerOptions = append(workerManagerOptions, workers.EnableLLMServerWorker)
		masaNodeOptions = append(masaNodeOptions, node.IsLlmServer)
	}

	workHandlerManager := workers.NewWorkHandlerManager(workerManagerOptions...)
	blockChainEventTracker := node.NewBlockChain()
	pubKeySub := &pubsub.PublicKeySubscriptionHandler{}

	// TODO: Where the config is involved, move to the config the generation of
	// Node options
	masaNodeOptions = append(masaNodeOptions, []node.Option{
		// Register the worker manager
		node.WithMasaProtocolHandler(
			config.WorkerProtocol,
			workHandlerManager.HandleWorkerStream,
		),
		node.WithPubSubHandler(config.PublicKeyTopic, pubKeySub, false),
		node.WithPubSubHandler(config.BlockTopic, blockChainEventTracker, true),
	}...)

	if cfg.Validator {
		// Subscribe and if actor start monitoring actor workers
		// considering all that matters is if the node is staked
		// and other peers can do work we only need to check this here
		// if this peer can or cannot scrape or write that is checked in other places
		masaNodeOptions = append(masaNodeOptions,
			node.WithService(blockChainEventTracker.Start(config.GetInstance().MasaDir)),
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
