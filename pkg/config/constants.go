package config

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
	PortNbr     = "PORT_NBR"
	UDP         = "UDP"
	TCP         = "TCP"
	PrivateKey  = "PRIVATE_KEY"
	StakeAmount = "STAKE_AMOUNT"
	LogLevel    = "LOG_LEVEL"
	LogFilePath = "LOG_FILEPATH"
	Environment = "ENV"
	AllowedPeer = "allowedPeer"
	Signature   = "signature"
	Debug       = "debug"
	Version     = "v0.0.9-alpha"
	DbPath      = "DB_PATH"
	FilePath    = "FILE_PATH"
	WriterNode  = "WRITER_NODE"

	MasaPrefix           = "/masa"
	OracleProtocol       = "oracle_protocol"
	NodeDataSyncProtocol = "nodeDataSync"
	NodeGossipTopic      = "gossip"
	AdTopic              = "ad"
	PublicKeyTopic       = "bootNodePublicKey"
	Rendezvous           = "masa-mdns"
	PageSize             = 25
)

func ProtocolWithVersion(protocolName string) protocol.ID {
	if GetInstance().Environment == "" {
		return protocol.ID(fmt.Sprintf("%s/%s/%s", MasaPrefix, protocolName, Version))
	}
	return protocol.ID(fmt.Sprintf("%s/%s/%s-%s", MasaPrefix, protocolName, Version, viper.GetString(Environment)))
}

func TopicWithVersion(protocolName string) string {
	if GetInstance().Environment == "" {
		return fmt.Sprintf("%s/%s/%s", MasaPrefix, protocolName, Version)
	}
	return fmt.Sprintf("%s/%s/%s-%s", MasaPrefix, protocolName, Version, viper.GetString(Environment))
}
