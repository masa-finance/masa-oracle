package network

import (
	"context"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReconnectToBootnodes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	bootnode, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"),
	)
	require.NoError(t, err)
	defer bootnode.Close()

	regularNode, err := libp2p.New()
	require.NoError(t, err)
	defer regularNode.Close()

	bootNodeAddr := bootnode.Addrs()[0].String() + "/p2p/" + bootnode.ID().String()

	tests := []struct {
		name            string
		bootnodes       []string
		expectConnected bool // New field to explicitly state connection expectation
	}{
		{
			name:            "successful connection to valid bootnode",
			bootnodes:       []string{bootNodeAddr},
			expectConnected: true,
		},
		{
			name:            "invalid bootnode address",
			bootnodes:       []string{"/ip4/256.256.256.256/tcp/1234/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN"},
			expectConnected: false,
		},
		{
			name:            "empty bootnode list",
			bootnodes:       []string{},
			expectConnected: false, // Changed to false since no connection is expected
		},
		{
			name:            "unreachable bootnode",
			bootnodes:       []string{"/ip4/127.0.0.1/tcp/1234/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN"},
			expectConnected: false,
		},
		{
			name: "multiple bootnodes with one valid",
			bootnodes: []string{
				"/ip4/256.256.256.256/tcp/1234/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
				bootNodeAddr,
			},
			expectConnected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Disconnect from any existing connections
			for _, conn := range regularNode.Network().Conns() {
				_ = conn.Close()
			}

			// Test reconnection
			reconnectToBootnodes(ctx, regularNode, tt.bootnodes)

			// Verify connection status matches expectation
			connected := isConnectedToAnyBootnode(regularNode, tt.bootnodes)
			assert.Equal(t, tt.expectConnected, connected,
				"Connection status mismatch: expected connected=%v, got connected=%v",
				tt.expectConnected, connected)

			// Add small delay to allow for connection cleanup
			time.Sleep(100 * time.Millisecond)
		})
	}
}

func isConnectedToAnyBootnode(h host.Host, bootnodes []string) bool {
	for _, bn := range bootnodes {
		if bn == "" {
			continue
		}
		ma, err := multiaddr.NewMultiaddr(bn)
		if err != nil {
			continue
		}
		pinfo, err := peer.AddrInfoFromP2pAddr(ma)
		if err != nil {
			continue
		}
		if h.Network().Connectedness(pinfo.ID) == network.Connected {
			return true
		}
	}
	return false
}
