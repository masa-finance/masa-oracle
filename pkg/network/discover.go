package network

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/sirupsen/logrus"
)

// Discover searches for and connects to peers that support the given protocol ID.
// It initializes discovery via the DHT and advertises this node.
// It runs discovery in a loop with a ticker, re-advertising and finding new peers.
// For each discovered peer, it checks if already connected, and if not, dials them.
func Discover(ctx context.Context, host host.Host, dht *dht.IpfsDHT, protocol protocol.ID) {
	protocolString := string(protocol)
	logrus.Infof("Discovering peers for protocol: %s", protocolString)
	routingDiscovery := routing.NewRoutingDiscovery(dht)

	// Advertise this node right away, then it will re-advertise with each ticker interval
	logrus.Infof("Attempting to advertise protocol: %s", protocolString)
	_, err := routingDiscovery.Advertise(ctx, protocolString)
	if err != nil {
		logrus.Warnf("Failed to advertise protocol: %v", err)
	} else {
		logrus.Infof("Successfully advertised protocol")
	}

	ticker := time.NewTicker(time.Second * 10) // 120
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if ctx.Err() != nil {
				logrus.Errorf("Context error in discovery loop: %v", ctx.Err())
			}
			logrus.Info("Stopping peer discovery")
			return
		case <-ticker.C:
			logrus.Debug("Searching for other peers...")

			routingDiscovery := routing.NewRoutingDiscovery(dht)

			// Advertise this node
			logrus.Debugf("Attempting to advertise protocol: %s", protocolString)
			_, err := routingDiscovery.Advertise(ctx, protocolString)
			if err != nil {
				logrus.Warnf("Failed to advertise protocol: %v", err)
			} else {
				logrus.Infof("Successfully advertised protocol")
			}

			// Use the routing discovery to find peers.
			peerChan, err := routingDiscovery.FindPeers(ctx, protocolString)
			if err != nil {
				logrus.Errorf("Failed to find peers: %v", err)
			} else {
				logrus.Debug("Successfully started finding peers")
			}
			select {
			case availPeer, ok := <-peerChan:
				if !ok {
					logrus.Info("Peer channel closed, restarting discovery")
					break
				}
				availPeerAddrInfo := peer.AddrInfo{
					ID:    availPeer.ID,
					Addrs: []multiaddr.Multiaddr{},
				}
				hostAddrInfo := peer.AddrInfo{
					ID:    host.ID(),
					Addrs: []multiaddr.Multiaddr{},
				}
				if availPeerAddrInfo.ID.String() == hostAddrInfo.ID.String() {
					logrus.Debugf("Skipping connect to self: %s", availPeer.String())
					continue
				}
				logrus.Infof("Available Peer: %s", availPeer.String())

				if host.Network().Connectedness(availPeer.ID) != network.Connected {
					_, err := host.Network().DialPeer(ctx, availPeer.ID)
					if err != nil {
						logrus.Warningf("Failed to connect to peer %s: %v", availPeer.ID.String(), err)
						continue
					}
					logrus.Infof("Connected to peer %s", availPeer.ID.String())
				}
			case <-ctx.Done():
				logrus.Info("Stopping peer discovery")
				return
			}
		}
	}
}
