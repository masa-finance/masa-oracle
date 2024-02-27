package masa

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/spf13/viper"
)

const (
	PrivKeyFile = "MASA_PRIV_KEY_FILE"
	BootNodes   = "BOOTNODES"
	MasaDir     = "MASA_DIR"
	RpcUrl      = "RPC_URL"
	Environment = "ENV"
	Version     = "v0.0.9-alpha"

	masaPrefix           = "/masa"
	oracleProtocol       = "oracle_protocol"
	NodeDataSyncProtocol = "nodeDataSync"
	NodeGossipTopic      = "gossip"
	AdTopic              = "ad"
	rendezvous           = "masa-mdns"
	PageSize             = 25
)

func ProtocolWithVersion(protocolName string) protocol.ID {
	if viper.GetString(Environment) == "" {
		return protocol.ID(fmt.Sprintf("%s/%s/%s", masaPrefix, protocolName, Version))
	}
	return protocol.ID(fmt.Sprintf("%s/%s/%s-%s", masaPrefix, protocolName, Version, viper.GetString(Environment)))
}

func TopicWithVersion(protocolName string) string {
	if viper.GetString(Environment) == "" {
		return fmt.Sprintf("%s/%s/%s", masaPrefix, protocolName, Version)
	}
	return fmt.Sprintf("%s/%s/%s-%s", masaPrefix, protocolName, Version, viper.GetString(Environment))
}
