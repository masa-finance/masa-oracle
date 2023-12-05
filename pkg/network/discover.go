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

	ticker := time.NewTicker(time.Second * 120)
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
				logrus.Errorf("Failed to advertise protocol: %v", err)
			} else {
				logrus.Debugf("Successfully advertised protocol")
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
				if availPeer.ID == host.ID() {
					logrus.Debugf("Skipping connect to self: %s", availPeer.String())
					continue
				}
				logrus.Infof("Available Peer: %s", availPeer.String())

				if host.Network().Connectedness(availPeer.ID) != network.Connected {
					//err := host.Connect(ctx, availPeer)
					_, err := host.Network().DialPeer(ctx, availPeer.ID)
					if err != nil {
						logrus.Warningf("Failed to connect to peer %s: %v", availPeer.ID.String(), err)
						continue
					}
					logrus.Infof("Connected to peer %s", availPeer.ID.String())
					//logrus.Infof("Connected to peer %s", conn.RemoteMultiaddr().String())
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
