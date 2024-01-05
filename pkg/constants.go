package masa

import (
	"fmt"

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
	Version              = "v.0.0.4-alpha"
)

func ProtocolWithVersion(protocolName string) protocol.ID {
	return protocol.ID(fmt.Sprintf("%s/%s/%s", masaPrefix, protocolName, Version))
}

func TopiclWithVersion(protocolName string) string {
	return fmt.Sprintf("%s/%s/%s", masaPrefix, protocolName, Version)
}
