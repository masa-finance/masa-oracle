package network

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/crypto"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
)

const (
	maxRetries  = 3
	retryDelay  = time.Second * 5
	PeerAdded   = "PeerAdded"
	PeerRemoved = "PeerRemoved"
)

func WithDht(ctx context.Context, host host.Host, bootstrapPeers []multiaddr.Multiaddr,
	protocolId, prefix protocol.ID, peerChan chan PeerEvent, isStaked bool) (*dht.IpfsDHT, error) {
	options := make([]dht.Option, 0)
	options = append(options, dht.Mode(dht.ModeAutoServer))
	options = append(options, dht.ProtocolPrefix(prefix))

	kademliaDHT, err := dht.New(ctx, host, options...)
	if err != nil {
		return nil, err
	}
	go monitorRoutingTable(ctx, kademliaDHT, time.Minute)

	kademliaDHT.RoutingTable().PeerAdded = func(p peer.ID) {
		logrus.Infof("Peer added to DHT: %s", p)
		pe := PeerEvent{
			AddrInfo: peer.AddrInfo{ID: p},
			Action:   PeerAdded,
			Source:   "kdht",
		}
		peerChan <- pe
	}

	kademliaDHT.RoutingTable().PeerRemoved = func(p peer.ID) {
		logrus.Infof("Peer removed from DHT: %s", p)
		pe := PeerEvent{
			AddrInfo: peer.AddrInfo{ID: p},
			Action:   PeerRemoved,
			Source:   "kdht",
		}
		peerChan <- pe
	}

	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for _, peerAddr := range bootstrapPeers {
		peerinfo, err := peer.AddrInfoFromP2pAddr(peerAddr)
		if err != nil {
			logrus.Errorf("kdht: %s", err.Error())
		}
		if peerinfo.ID == host.ID() {
			logrus.Info("DHT Skipping connect to self")
			continue
		}
		// Add the bootstrap node to the DHT
		added, err := kademliaDHT.RoutingTable().TryAddPeer(peerinfo.ID, true, false)
		if err != nil {
			logrus.Warningf("Failed to add bootstrap peer %s to DHT: %v", peerinfo.ID, err)
		} else if !added {
			logrus.Warningf("Bootstrap peer %s was not added to DHT", peerinfo.ID)
		} else {
			logrus.Infof("Successfully added bootstrap peer %s to DHT", peerinfo.ID)
		}

		wg.Add(1)
		go func() {
			ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel() // Cancel the context when done to release resources

			defer wg.Done()
			if err := host.Connect(ctxWithTimeout, *peerinfo); err != nil {
				logrus.Errorf("Failed to connect to bootstrap peer %s: %v", peerinfo.ID, err)
				time.Sleep(retryDelay)
			} else {
				logrus.Info("Connection established with node:", *peerinfo)
				stream, err := host.NewStream(ctxWithTimeout, peerinfo.ID, protocolId)
				if err != nil {
					logrus.Error("Error opening stream:", err)
					return
				}
				defer func(stream network.Stream) {
					err := stream.Close()
					if err != nil {
						logrus.Error("Error closing stream:", err)
					}
				}(stream) // Close the stream when done

				_, err = stream.Write(getSelfNodeDataJson(host, isStaked))
				if err != nil {
					logrus.Error("Error writing to stream:", err)
					return
				}
			}
		}()
	}
	wg.Wait()
	return kademliaDHT, nil
}

func getSelfNodeDataJson(host host.Host, isStaked bool) []byte {
	var publicKeyHex string
	var err error
	// Get the public key of the host node
	pubKey := host.Peerstore().PubKey(host.ID())
	if pubKey == nil {
		logrus.WithFields(logrus.Fields{
			"Peer": host.ID().String(),
		}).Warn("No public key found for peer")
	} else {
		publicKeyHex, err = crypto.Libp2pPubKeyToEthAddress(pubKey)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Peer": host.ID().String(),
			}).Warnf("Error getting public key %v", err)
		}
	}

	// Create and populate NodeData
	nodeData := pubsub.NodeData{
		PeerId:     host.ID(),
		IsStaked:   isStaked,
		EthAddress: publicKeyHex,
	}

	// Convert NodeData to JSON
	jsonData, err := json.Marshal(nodeData)
	if err != nil {
		logrus.Error("Error marshalling NodeData:", err)
		return nil
	}
	return jsonData
}

func monitorRoutingTable(ctx context.Context, dht *dht.IpfsDHT, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// This block will be executed every 'interval' duration
			routingTable := dht.RoutingTable()
			// Log the size of the routing table
			logrus.Debugf("Routing table size: %d", routingTable.Size())
			// Log the peer IDs in the routing table
			for _, p := range routingTable.ListPeers() {
				logrus.Debugf("Peer in routing table: %s", p.String())
			}
		case <-ctx.Done():
			// If the context is cancelled, stop the goroutine
			return
		}
	}
}
