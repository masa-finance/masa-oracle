package main

import (
	"context"
	"fmt"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
)

func main() {
	ctx := context.Background()

	// Generate a new random identity
	privateKey, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		panic(err)
	}

	// Create a new libp2p node with the generated identity
	node, err := libp2p.New(ctx, libp2p.Identity(privateKey))
	if err != nil {
		panic(err)
	}

	fmt.Println("Node created with ID:", node.ID())
}
