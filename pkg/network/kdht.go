package network

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

const (
	maxRetries  = 3
	retryDelay  = time.Second * 5
	PeerAdded   = "PeerAdded"
	PeerRemoved = "PeerRemoved"
)

func WithDht(ctx context.Context, host host.Host, bootstrapPeers []multiaddr.Multiaddr,
	pId protocol.ID, peerChan chan PeerEvent) (*dht.IpfsDHT, error) {
	options := make([]dht.Option, 0)
	options = append(options, dht.Mode(dht.ModeAutoServer))
	options = append(options, dht.ProtocolPrefix(pId))

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
			defer wg.Done()
			for i := 0; i < maxRetries; i++ {
				if err := host.Connect(ctx, *peerinfo); err != nil {
					logrus.Errorf("Failed to connect to bootstrap peer %s: %v", peerinfo.ID, err)
					time.Sleep(retryDelay)
				} else {
					logrus.Info("Connection established with node:", *peerinfo)
					stream, err := host.NewStream(ctx, peerinfo.ID, pId)
					if err != nil {
						logrus.Error("Error opening stream:", err)
						return
					}
					_, err = stream.Write([]byte(fmt.Sprintf("Initial Hello from %s\n", peerAddr.String())))
					if err != nil {
						logrus.Error("Error writing to stream:", err)
						return
					}
					break
				}
			}
		}()
	}
	wg.Wait()
	return kademliaDHT, nil
}

func NewDHT(ctx context.Context, host host.Host, protocol protocol.ID, bootstrapPeers []multiaddr.Multiaddr) (*dht.IpfsDHT, error) {
	var options []dht.Option

	// if no bootstrap peers give this peer act as a bootstraping node
	// other peers can use this peers ipfs address for peer discovery via dht
	if len(bootstrapPeers) == 0 {
		log.Info("running as server")
		options = append(options, dht.Mode(dht.ModeAutoServer))
	}
	options = append(options, dht.ProtocolPrefix(protocol))

	kdht, err := dht.New(ctx, host, options...)
	if err != nil {
		return nil, err
	}
	kdht.RoutingTable().PeerAdded = func(p peer.ID) {
		logrus.Infof("Peer added to DHT: %s", p)
	}

	kdht.RoutingTable().PeerRemoved = func(p peer.ID) {
		logrus.Infof("Peer removed from DHT: %s", p)
	}
	go monitorRoutingTable(ctx, kdht, time.Minute)

	if err = kdht.Bootstrap(ctx); err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for _, peerAddr := range bootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(ctx, *peerinfo); err != nil {
				logrus.Error("Error while connecting to node %q: %-v", peerinfo, err)
			} else {
				logrus.Info("Connection established with bootstrap node: %q", *peerinfo)
			}
		}()
	}
	wg.Wait()

	return kdht, nil
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
			logrus.Infof("Routing table size: %d", routingTable.Size())
			// Log the peer IDs in the routing table
			for _, p := range routingTable.ListPeers() {
				logrus.Infof("Peer in routing table: %s", p.String())
			}
		case <-ctx.Done():
			// If the context is cancelled, stop the goroutine
			return
		}
	}
}
