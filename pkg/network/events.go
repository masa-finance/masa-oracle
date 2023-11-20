package network

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

type PeerEvent struct {
	AddrInfo   peer.AddrInfo
	Action     string
	Source     string
	Rendezvous string
}
