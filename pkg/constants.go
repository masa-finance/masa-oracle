package masa

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/spf13/viper"
)

const (
	Cert                 = "cert"
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
	Version              = "v0.0.6-alpha"
	Environment          = "ENV"
)

CertPem = viper.getString("CERT_PEM")

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
