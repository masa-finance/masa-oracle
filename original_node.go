package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	"github.com/libp2p/go-libp2p/p2p/transport/websocket"
	"github.com/multiformats/go-multiaddr"

	"github.com/masa-finance/masa-oracle/crypto"
)

func run() {
	// Start with the default scaling limits.
	scalingLimits := rcmgr.DefaultLimits
	concreteLimits := scalingLimits.AutoScale()
	limiter := rcmgr.NewFixedLimiter(concreteLimits)

	rm, err := rcmgr.NewResourceManager(limiter)
	if err != nil {
		panic(err)
	}

	// Get or create the private key
	privKey, err := crypto.GetOrCreatePrivateKey("")
	if err != nil {
		panic(err)
	}

	h, err := libp2p.New(
		libp2p.Transport(websocket.New),
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0/ws"),
		libp2p.ResourceManager(rm),
		libp2p.Identity(privKey),
		libp2p.Ping(false), // disable built-in ping
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
	)
	if err != nil {
		panic(err)
	}

	// Configure our own heartbeat protocol
	heartbeatSvc := ping.NewPingService(h)
	h.SetStreamHandler(ping.ID, func(s network.Stream) {
		fmt.Println("Heartbeat received!")
		heartbeatSvc.PingHandler(s)
	})

	peerInfo := peer.AddrInfo{
		ID:    h.ID(),
		Addrs: h.Addrs(),
	}
	multiaddrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("libp2p host address:", multiaddrs[0])

	dhtInstance, err := dht.New(context.Background(), h)
	if err != nil {
		panic(err)
	}

	if err = dhtInstance.Bootstrap(context.Background()); err != nil {
		panic(err)
	}

	if len(os.Args) > 1 {
		addr, err := multiaddr.NewMultiaddr(os.Args[1])
		if err != nil {
			panic(err)
		}

		peerInfo, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			panic(err)
		}

		if err := h.Connect(context.Background(), *peerInfo); err != nil {
			fmt.Println("Error connecting to peer:", err)
			return
		}

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			ch := heartbeatSvc.Ping(context.Background(), peerInfo.ID)
			res := <-ch
			fmt.Println("Heartbeat sent to", addr, "received in", res.RTT)
		}
	} else {
		// Wait for a SIGINT or SIGTERM signal
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		fmt.Println("Received signal, shutting down...")
	}

	// Shut the node down
	if err := h.Close(); err != nil {
		panic(err)
	}
}
