package config

import (
	"github.com/masa-finance/masa-oracle/node"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("InitOptions", func() {
	It("parses the AppConfig into NodeOptions", func() {
		conf := AppConfig{
			Version:         "1.0",
			PortNbr:         42,
			UDP:             true,
			TCP:             true,
			Bootnodes:       []string{"boot1", "boot2"},
			Environment:     "test",
			MasaDir:         "dir",
			Validator:       true,
			CachePath:       "cache",
			TwitterScraper:  true,
			DiscordScraper:  true,
			TelegramScraper: true,
			WebScraper:      true,
		}

		opts, _, _ := InitOptions(&conf)

		actual := &node.NodeOption{}
		actual.Apply(opts...)

		// These have pointers to functions which we can't really compare.r
		// Check that they aren't nil, then set them to nil for the rest of the test.
		Expect(actual.Services).To(Not(BeNil()))
		Expect(actual.PubSubHandles).To(Not(BeNil()))
		Expect(actual.MasaProtocolHandlers).To(Not(BeNil()))
		actual.Services = nil
		actual.PubSubHandles = nil
		actual.MasaProtocolHandlers = nil

		expected := node.NodeOption{
			IsStaked:             true,
			UDP:                  conf.UDP,
			TCP:                  conf.TCP,
			IsValidator:          conf.Validator,
			PortNbr:              conf.PortNbr,
			IsTwitterScraper:     conf.TwitterScraper,
			IsDiscordScraper:     conf.DiscordScraper,
			IsTelegramScraper:    conf.TelegramScraper,
			IsWebScraper:         conf.WebScraper,
			Bootnodes:            conf.Bootnodes,
			RandomIdentity:       false,
			ProtocolHandlers:     nil,
			Environment:          conf.Environment,
			Version:              conf.Version,
			MasaDir:              conf.MasaDir,
			CachePath:            conf.CachePath,
			OracleProtocol:       OracleProtocol,
			NodeDataSyncProtocol: NodeDataSyncProtocol,
			NodeGossipTopic:      NodeGossipTopic,
			Rendezvous:           Rendezvous,
			WorkerProtocol:       WorkerProtocol,
			PageSize:             PageSize,

			// Set these to the same values we set above
			Services:             actual.Services,
			PubSubHandles:        actual.PubSubHandles,
			MasaProtocolHandlers: actual.MasaProtocolHandlers,
		}

		Expect(*actual).To(Equal(expected))
	})
})
