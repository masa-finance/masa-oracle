package network

import (
	"context"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/masa-finance/masa-oracle/pkg/config"
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
	var routingDiscovery *routing.RoutingDiscovery
	protocolString := string(protocol)
	logrus.Infof("Discovering peers for protocol: %s", protocolString)

	routingDiscovery = routing.NewRoutingDiscovery(dht)

	// Advertise node right away, then it will re-advertise with each ticker interval
	logrus.Infof("Attempting to advertise protocol: %s", protocolString)
	_, err := routingDiscovery.Advertise(ctx, protocolString)
	if err != nil {
		logrus.Debugf("Failed to advertise protocol: %v", err)
	} else {
		logrus.Infof("Successfully advertised protocol %s", protocolString)
	}

	ticker := time.NewTicker(time.Minute * 1)
	defer ticker.Stop()

	var peerChan <-chan peer.AddrInfo

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
			routingDiscovery = routing.NewRoutingDiscovery(dht)

			// Advertise this node
			logrus.Debugf("Attempting to advertise protocol: %s", protocolString)
			_, err := routingDiscovery.Advertise(ctx, protocolString)
			if err != nil {
				logrus.Debugf("Failed to advertise protocol with error %v", err)

				// Network retry when connectivity is temporarily lost using NewExponentialBackOff
				expBackOff := backoff.NewExponentialBackOff()
				expBackOff.MaxElapsedTime = time.Second * 10
				err := backoff.Retry(func() error {
					peerChan, err = routingDiscovery.FindPeers(ctx, protocolString)
					return err
				}, expBackOff)
				if err != nil {
					logrus.Warningf("Retry failed to find peers: %v", err)
				}

			} else {
				logrus.Infof("Successfully advertised protocol: %s", protocolString)
			}

			// Use the routing discovery to find peers.
			peerChan, err = routingDiscovery.FindPeers(ctx, protocolString)
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
				// validating proper peers to connect to
				availPeerAddrInfo := peer.AddrInfo{
					ID:    availPeer.ID,
					Addrs: availPeer.Addrs,
				}
				hostAddrInfo := peer.AddrInfo{
					ID:    host.ID(),
					Addrs: host.Addrs(),
				}
				if availPeerAddrInfo.ID.String() == hostAddrInfo.ID.String() {
					logrus.Debugf("Skipping connect to self: %s", availPeerAddrInfo.ID.String())
					continue
				}
				if len(availPeerAddrInfo.Addrs) == 0 {
					for _, bn := range config.GetInstance().Bootnodes {
						bootNode := strings.Split(bn, "/")[len(strings.Split(bn, "/"))-1]
						if availPeerAddrInfo.ID.String() != bootNode {
							logrus.Warningf("Skipping connect to non bootnode peer with no multiaddress: %s", availPeerAddrInfo.ID.String())
							continue
						}
					}
				}
				logrus.Infof("Available Peer: %s", availPeer.String())

				if host.Network().Connectedness(availPeer.ID) != network.Connected {
					if isConnectedToBootnode(host, config.GetInstance().Bootnodes) {
						_, err := host.Network().DialPeer(ctx, availPeer.ID)
						if err != nil {
							logrus.Warningf("Failed to connect to peer %s, will retry...", availPeer.ID.String())
							continue
						} else {
							logrus.Infof("Connected to peer %s", availPeer.ID.String())
						}
					} else {
						for _, bn := range config.GetInstance().Bootnodes {
							if len(bn) > 0 {
								logrus.Info("Not connected to any bootnode. Attempting to reconnect...")
								reconnectToBootnodes(ctx, host, config.GetInstance().Bootnodes)
							}
						}
					}
				}
			case <-ctx.Done():
				logrus.Info("Stopping peer discovery")
				return
			}
		}
	}
}

// Function to check if connected to at least one bootnode
func isConnectedToBootnode(host host.Host, bootnodes []string) bool {
	for _, bn := range bootnodes {
		peerID, _ := peer.Decode(strings.Split(bn, "/")[len(strings.Split(bn, "/"))-1])
		if host.Network().Connectedness(peerID) == network.Connected {
			return true
		}
	}
	return false
}

// Function to attempt reconnection to bootnodes
func reconnectToBootnodes(ctx context.Context, host host.Host, bootnodes []string) {
	for _, bn := range bootnodes {
		ma, _ := multiaddr.NewMultiaddr(bn)
		peerInfo, _ := peer.AddrInfoFromP2pAddr(ma)
		if err := host.Connect(ctx, *peerInfo); err != nil {
			logrus.Errorf("Failed to reconnect to bootnode %s: %v", bn, err)
		} else {
			logrus.Infof("Reconnected to bootnode %s", bn)
			break // Exit after successful reconnection to one bootnode
		}
	}
}
