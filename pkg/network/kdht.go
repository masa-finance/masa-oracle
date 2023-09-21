package network

import (
	"context"
	"fmt"
	"sync"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

func NewDht(ctx context.Context, host host.Host, bootstrapPeers []multiaddr.Multiaddr,
	protocol protocol.ID, address multiaddr.Multiaddr) (*dht.IpfsDHT, error) {
	options := make([]dht.Option, 0)
	if len(bootstrapPeers) == 0 {
		options = append(options, dht.Mode(dht.ModeServer))
	}
	kademliaDHT, err := dht.New(ctx, host, options...)
	if err != nil {
		return nil, err
	}
	kademliaDHT.RoutingTable().PeerAdded = func(p peer.ID) {
		logrus.Infof("Peer added to DHT: %s", p)
	}

	kademliaDHT.RoutingTable().PeerRemoved = func(p peer.ID) {
		logrus.Infof("Peer removed from DHT: %s", p)
	}

	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for _, peerAddr := range bootstrapPeers {
		peerinfo, err := peer.AddrInfoFromP2pAddr(peerAddr)
		if err != nil {
			logrus.Error(err)
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
			defer wg.Done()
			if err := host.Connect(ctx, *peerinfo); err != nil {
				logrus.Errorf("Failed to connect to bootstrap peer %s: %v", peerinfo.ID, err)
			} else {
				logrus.Info("Connection established with node:", *peerinfo)
				stream, err := host.NewStream(ctx, peerinfo.ID, protocol)
				if err != nil {
					logrus.Error("Error opening stream:", err)
				}
				_, err = stream.Write([]byte(fmt.Sprintf("Initial Hello from %s\n", address.String())))
				if err != nil {
					logrus.Error("Error writing to stream:", err)
				}
			}
		}()
	}
	wg.Wait()
	return kademliaDHT, nil
}
