package network

import (
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

// HandlePeerFound interface to be called when new  peer is found
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

// StartMDNS Initializes and starts the MDNS service
func StartMDNS(host host.Host, rendezvous string) chan peer.AddrInfo {
	// register with service so that we get notified about peer discovery
	notifee := &discoveryNotifee{
		PeerChan: make(chan peer.AddrInfo),
	}

	mdnsService := mdns.NewMdnsService(host, rendezvous, notifee)
	if err := mdnsService.Start(); err != nil {
		panic(err)
	}
	return notifee.PeerChan
}
