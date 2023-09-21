package network

import (
	"context"
	"fmt"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

func Discover(ctx context.Context, host host.Host, dht *dht.IpfsDHT, protocol protocol.ID, address multiaddr.Multiaddr) {
	protocolString := string(protocol)
	logrus.Infof("Discovering peers for protocol: %s", protocolString)
	routingDiscovery := routing.NewRoutingDiscovery(dht)

	// Advertise this node
	logrus.Infof("Attempting to advertise protocol: %s", protocolString)
	_, err := routingDiscovery.Advertise(ctx, protocolString)
	if err != nil {
		logrus.Errorf("Failed to advertise protocol: %v", err)
	} else {
		logrus.Infof("Successfully advertised protocol")
	}

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			logrus.Debug("Searching for other peers...")
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
				if availPeer.ID == host.ID() {
					logrus.Debugf("Skipping connect to self: %s", availPeer.String())
					continue
				}
				logrus.Infof("Found availPeer: %s", availPeer.String())

				if host.Network().Connectedness(availPeer.ID) != network.Connected {
					_, err = host.Network().DialPeer(ctx, availPeer.ID)
					fmt.Printf("Connected to peer %s\n", availPeer.ID.String())
					if err != nil {
						continue
					}
				}
				// Send a message with this node's multi address string to each availPeer that is found
				stream, err := host.NewStream(ctx, availPeer.ID, protocol)
				if err != nil {
					logrus.Error("Error opening stream:", err)
					continue
				}
				_, err = stream.Write([]byte(fmt.Sprintf("Discovery Hello from %s", address.String())))
				if err != nil {
					logrus.Error("Error writing to stream:", err)
					continue
				}
			case <-ctx.Done():
				logrus.Info("Stopping peer discovery")
				return
			}
		}
	}
	logrus.Infof("found %d peers", len(host.Network().Peers()))
}
