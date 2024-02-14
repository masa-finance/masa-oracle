package masa

import (
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p/core/protocol"
)

const (
	KeyFileKey           = "private.key"
	CertPem              = "cert.pem"
	Cert                 = "cert"
	Peers                = "peerList"
	masaPrefix           = "/masa"
	oracleProtocol       = "oracle_protocol"
	NodeDataSyncProtocol = "nodeDataSync"
	NodeGossipTopic      = "gossip"
	AdTopic              = "ad"
	rendezvous           = "masa-mdns"
	PortNbr              = "portNbr"
	PageSize             = 25
	NodeBackupFileName   = "nodeBackup.json"
	NodeBackupPath       = "nodeBackupPath"
	Version              = "v0.0.6-alpha"
	DefaultRPCURL        = "https://ethereum-sepolia.publicnode.com"
	Environment          = "ENV"
)

var env string

func ProtocolWithVersion(protocolName string) protocol.ID {
	if getEnv() == "" {
		return protocol.ID(fmt.Sprintf("%s/%s/%s", masaPrefix, protocolName, Version))
	}
	return protocol.ID(fmt.Sprintf("%s/%s/%s-%s", masaPrefix, protocolName, Version, getEnv()))
}

func TopicWithVersion(protocolName string) string {
	if getEnv() == "" {
		return fmt.Sprintf("%s/%s/%s", masaPrefix, protocolName, Version)
	}
	return fmt.Sprintf("%s/%s/%s-%s", masaPrefix, protocolName, Version, getEnv())
}

func getEnv() string {
	if env == "" {
		env = os.Getenv(Environment)
	}
	return env
}
