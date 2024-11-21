package network

import (
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

const PeerFound = "PeerFound"

type discoveryNotifee struct {
	PeerChan   chan PeerEvent
	Rendezvous string
}

// HandlePeerFound interface to be called when new  peer is found
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	pe := PeerEvent{
		AddrInfo: pi,
		// WTF: Use const
		Action:     "PeerFound",
		Source:     "mdns",
		Rendezvous: n.Rendezvous,
	}
	n.PeerChan <- pe
}

func EnableMDNS(host host.Host, rendezvous string, peerChan chan PeerEvent) error {
	notifee := &discoveryNotifee{
		PeerChan:   peerChan,
		Rendezvous: rendezvous,
	}
	mdnsService := mdns.NewMdnsService(host, rendezvous, notifee)
	if err := mdnsService.Start(); err != nil {
		return err
	}
	return nil
}
