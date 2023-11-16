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

func WithMDNS(host host.Host, rendezvous string, peerChan chan peer.AddrInfo) {
	notifee := &discoveryNotifee{
		PeerChan: peerChan,
	}
	mdnsService := mdns.NewMdnsService(host, rendezvous, notifee)
	if err := mdnsService.Start(); err != nil {
		panic(err)
	}
}
