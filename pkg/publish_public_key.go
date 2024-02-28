package masa

import (
	"github.com/masa-finance/masa-oracle/pkg/consensus"
	pubsub2 "github.com/masa-finance/masa-oracle/pkg/pubsub"
)

// PublishNodePublicKey publishes the node's public key and signed data to the designated topic.
func (node *OracleNode) PublishPublicKey() error {
	keyManager := KeyManager{}

	// Load the private key directly using the KeyManager from keys package
	privKey, err := keyManager.LoadPrivKey()
	if err != nil {
		return err
	}

	// Load the public key directly using the KeyManager from keys package
	pubKey, err := keyManager.LoadPubKey()
	if err != nil {
		return err
	}

	// Convert the public key to a string representation
	pubKeyString, err := PubKeyToString(pubKey)
	if err != nil {
		return err
	}

	// Set the data to be signed as the signers Peer ID
	data := []byte(node.Host.ID().String())

	// Sign the data using the private key
	signature, err := consensus.SignData(privKey, data)
	if err != nil {
		return err
	}

	// Create a new PublicKeyPublisher instance
	publisher := pubsub2.NewPublicKeyPublisher(node.PubSubManager, pubKey)

	// Publish the public key using its string representation, data, and signature
	return publisher.PublishNodePublicKey(pubKeyString, data, signature)
}
