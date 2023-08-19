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
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	multiaddr "github.com/multiformats/go-multiaddr"
)

const keyFilePath = "private.key"

func getOrCreatePrivateKey() (crypto.PrivKey, error) {
	// Check if the private key file exists
	data, err := os.ReadFile(keyFilePath)
	if err == nil {
		fmt.Printf("Raw data from %s: %s\n", keyFilePath, string(data))
		// Decode the private key from the file
		privKey, err := crypto.UnmarshalPrivateKey(data)
		if err != nil {
			fmt.Printf("Error unmarshalling private key: %s\n", err)
			return nil, err
		}
		fmt.Printf("Loaded private key from %s: %s\n", keyFilePath, privKey)
		return privKey, nil
	} else {
		// Generate a new private key
		privKey, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
		if err != nil {
			return nil, err
		}
		// Marshal the private key to bytes
		data, err := crypto.MarshalPrivateKey(privKey)
		if err != nil {
			return nil, err
		}
		// Save the private key to the file
		if err := os.WriteFile(keyFilePath, data, 0600); err != nil {
			return nil, err
		}
		fmt.Printf("Generated and saved a new private key to %s: %s\n", keyFilePath, privKey)
		return privKey, nil
	}
}

func main() {
	// Start with the default scaling limits.
	scalingLimits := rcmgr.DefaultLimits
	concreteLimits := scalingLimits.AutoScale()
	limiter := rcmgr.NewFixedLimiter(concreteLimits)

	rm, err := rcmgr.NewResourceManager(limiter)
	if err != nil {
		panic(err)
	}

	// Get or create the private key
	privKey, err := getOrCreatePrivateKey()
	if err != nil {
		panic(err)
	}

	h, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
		libp2p.ResourceManager(rm),
		libp2p.Identity(privKey),
		libp2p.Ping(false), // disable built-in ping
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
