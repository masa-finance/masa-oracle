package tests

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	libp2pquic "github.com/libp2p/go-libp2p/p2p/transport/quic"
)

type MasaOracle struct {
	Node host.Host
}

func (m *MasaOracle) start() {
	m.Node.SetStreamHandler("/masa/1.0.0", m.handleStream)
	fmt.Println("MasaOracle node started and listening on:", m.Node.Addrs())
}

func (m *MasaOracle) handleStream(s network.Stream) {
	log.Println("Got a new stream!")
	// Handle your stream here
}

func NewMasaOracle(privKey crypto.PrivKey) (*MasaOracle, error) {
	ctx := context.Background()

	// Create a QUIC transport
	quicTransport, err := libp2pquic.NewTransport(privKey)
	if err != nil {
		return nil, err
	}

	// Create a libp2p host with the QUIC transport
	node, err := libp2p.New(ctx,
		libp2p.Identity(privKey),
		libp2p.Transport(quicTransport),
	)
	if err != nil {
		return nil, err
	}

	return &MasaOracle{Node: node}, nil
}

func main() {
	// Generate a new key pair
	priv, _, err := crypto.GenerateECDSAKeyPair(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new MasaOracle
	masaOracle, err := NewMasaOracle(priv)
	if err != nil {
		log.Fatal(err)
	}

	// Start the MasaOracle
	masaOracle.start()
}
