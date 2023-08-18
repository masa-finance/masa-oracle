// /infrastructure/libp2p/node_config.go

package libp2p

import (
	"context"

	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	host "github.com/libp2p/go-libp2p-core/host"
)

// SetupNode initializes a libp2p node with an RSA identity.
// Returns the created node and any error that might occur.
func SetupNode(ctx context.Context) (host.Host, error) {
	// Generate an RSA key pair for this host.
	privKey, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		return nil, err
	}

	// Create a new instance of libp2p with the given private key.
	return libp2p.New(libp2p.Identity(privKey))
}
