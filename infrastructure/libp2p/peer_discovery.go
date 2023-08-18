// /infrastructure/libp2p/peer_discovery.go

package libp2p

import (
	"context"

	datastore "github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p-core/host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

// SetupDHT initializes a libp2p DHT for the given host.
// Returns the DHT instance.
func SetupDHT(ctx context.Context, node host.Host) *dht.IpfsDHT {
	// Use a simple map-based datastore for the DHT.
	dstore := datastore.NewMapDatastore()

	// Initialize the DHT with the datastore.
	dhtInstance, _ := dht.New(ctx, node, dht.Datastore(dstore))

	// Bootstrap the DHT to ensure it's connected to the network.
	dhtInstance.Bootstrap(ctx)

	return dhtInstance
}
