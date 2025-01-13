package network

import (
	"testing"

	"github.com/libp2p/go-libp2p"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMultiAddressesForHost(t *testing.T) {
	tests := []struct {
		name        string
		listenAddrs []string
		expectEmpty bool
		expectError bool
	}{
		{
			name: "valid non-local addresses",
			listenAddrs: []string{
				"/ip4/0.0.0.0/tcp/0", // Use 0.0.0.0 and port 0 for testing
			},
			expectEmpty: false,
			expectError: false,
		},
		{
			name: "only localhost addresses",
			listenAddrs: []string{
				"/ip4/127.0.0.1/tcp/0",
			},
			expectEmpty: true,
			expectError: false,
		},
		{
			name: "mixed addresses",
			listenAddrs: []string{
				"/ip4/127.0.0.1/tcp/0",
				"/ip4/0.0.0.0/tcp/0",
			},
			expectEmpty: false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test host with specified listen addresses
			opts := []libp2p.Option{
				libp2p.ListenAddrStrings(tt.listenAddrs...),
			}
			h, err := libp2p.New(opts...)
			require.NoError(t, err)
			defer h.Close()

			// Get multiaddresses
			addrs, err := GetMultiAddressesForHost(h)

			// Check error
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Check if empty when expected
			if tt.expectEmpty {
				assert.Empty(t, addrs)
				return
			}

			// Verify no localhost addresses
			for _, addr := range addrs {
				assert.NotContains(t, addr.String(), "127.0.0.1")
			}
		})
	}
}

func TestGetBootNodesMultiAddress(t *testing.T) {
	tests := []struct {
		name           string
		bootstrapNodes []string
		expectedLen    int
		expectError    bool
	}{
		{
			name: "valid addresses",
			bootstrapNodes: []string{
				"/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
				"/ip4/104.131.131.83/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuK",
			},
			expectedLen: 2,
			expectError: false,
		},
		{
			name:           "empty list",
			bootstrapNodes: []string{},
			expectedLen:    0,
			expectError:    false,
		},
		{
			name:           "list with empty string",
			bootstrapNodes: []string{""},
			expectedLen:    0,
			expectError:    false,
		},
		{
			name: "invalid address",
			bootstrapNodes: []string{
				"invalid-address",
			},
			expectedLen: 0,
			expectError: true,
		},
		{
			name: "mixed valid and empty",
			bootstrapNodes: []string{
				"",
				"/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
			},
			expectedLen: 1,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addrs, err := GetBootNodesMultiAddress(tt.bootstrapNodes)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, addrs, tt.expectedLen)

			// Verify each address is valid
			for _, addr := range addrs {
				assert.NotNil(t, addr)
				assert.Implements(t, (*multiaddr.Multiaddr)(nil), addr)
			}
		})
	}
}
