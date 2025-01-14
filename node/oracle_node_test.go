package node

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestNodeSignature(t *testing.T) {
	// Generate a new private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	// Create a data string
	data := "This is some test data"

	// Hash the data
	hash := crypto.Keccak256Hash([]byte(data))

	// Sign the hash with the private key
	sig, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new OracleNode
	node := &OracleNode{
		Signature: hex.EncodeToString(sig),
	}

	// Check if the node is a publisher
	if !node.IsPublisher() {
		t.Errorf("Expected node to be a publisher, but it's not")
	}
}

func TestOracleNode_GetNodeData(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		opts    []Option
		checkFn func(*testing.T, *pubsub.NodeData)
	}{
		{
			name: "random identity node",
			opts: []Option{
				EnableRandomIdentity,
				EnableUDP,
				WithPort(0),
			},
			checkFn: func(t *testing.T, data *pubsub.NodeData) {
				assert.NotEmpty(t, data.PeerId)
				assert.True(t, strings.HasPrefix(data.EthAddress, "0x"))
				assert.False(t, data.IsStaked)
			},
		},
		{
			name: "node with all capabilities",
			opts: []Option{
				EnableRandomIdentity,
				EnableUDP,
				WithPort(0),
				EnableStaked,
				IsTwitterScraper,
				IsWebScraper,
				IsValidator,
			},
			checkFn: func(t *testing.T, data *pubsub.NodeData) {
				assert.NotEmpty(t, data.PeerId)
				assert.True(t, data.IsStaked)
				assert.True(t, data.IsTwitterScraper)
				assert.True(t, data.IsWebScraper)
				assert.True(t, data.IsValidator)
				assert.True(t, data.IsActive)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := NewOracleNode(ctx, tt.opts...)
			require.NoError(t, err)
			defer node.Host.Close()

			nodeData := node.getNodeData()
			require.NotNil(t, nodeData)
			tt.checkFn(t, nodeData)
		})
	}
}
