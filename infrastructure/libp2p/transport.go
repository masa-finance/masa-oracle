// /infrastructure/libp2p/transport.go

package libp2p

import (
	"context"

	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	host "github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"
)

// SetupTransport initializes a libp2p node with TCP and WebSocket transports.
// Returns the created node and any error that might occur.
func SetupTransport(ctx context.Context, privKey crypto.PrivKey) (host.Host, error) {
	return libp2p.New(
		libp2p.Identity(privKey),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(ws.New),
	)
}
