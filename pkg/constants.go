package masa

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/spf13/viper"
)

const (
	Cert                 = "cert"
	CertPem              = "cert.pem"
	Peers                = "peerList"
	masaPrefix           = "/masa"
	oracleProtocol       = "oracle_protocol"
	NodeDataSyncProtocol = "nodeDataSync"
	NodeGossipTopic      = "gossip"
	AdTopic              = "ad"
	rendezvous           = "masa-mdns"
	PageSize             = 25
	NodeBackupFileName   = "nodeBackup.json"
	NodeBackupPath       = "nodeBackupPath"
	Version              = "v0.0.7-alpha"
	DefaultRPCURL        = "https://ethereum-sepolia.publicnode.com"
	Environment          = "ENV"
)

func ProtocolWithVersion(protocolName string) protocol.ID {
	if viper.GetString("ENV") == "" {
		return protocol.ID(fmt.Sprintf("%s/%s/%s", masaPrefix, protocolName, Version))
	}
	return protocol.ID(fmt.Sprintf("%s/%s/%s-%s", masaPrefix, protocolName, Version, viper.GetString("ENV")))
}

func TopicWithVersion(protocolName string) string {
	if viper.GetString("ENV") == "" {
		return fmt.Sprintf("%s/%s/%s", masaPrefix, protocolName, Version)
	}
	return fmt.Sprintf("%s/%s/%s-%s", masaPrefix, protocolName, Version, viper.GetString("ENV"))
}
