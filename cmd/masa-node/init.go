package main

import (
	"encoding/hex"
	"log"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/config"
	masaCrypto "github.com/masa-finance/masa-oracle/pkg/crypto"
)

func init() {

	cfg := config.GetInstance()
	// Load or generate the private key
	privKey, _, _, err := masaCrypto.GetOrCreatePrivateKey(cfg.PrivateKeyFile)
	if err != nil {
		log.Fatalf("Failed to load or create private key: %v", err)
	}

	// Correctly handle the public key
	pubKey := privKey.GetPublic()

	// Assuming you implement or correct the implementation of MarshalPublicKey
	pubKeyBytes, err := crypto.MarshalPublicKey(pubKey)
	if err != nil {
		log.Fatalf("Failed to marshal public key: %v", err)
	}

	allowedPeerPubKeyHex := hex.EncodeToString(pubKeyBytes)

	// Assuming you implement or correct the implementation of Libp2pPubKeyToPeerID
	peerID, err := peer.IDFromPublicKey(pubKey)
	if err != nil {
		log.Fatalf("Failed to convert public key to peer ID: %v", err)
	}

	if cfg.AllowedPeer {
		cfg.AllowedPeerId = peerID.String()
		cfg.AllowedPeerPublicKey = allowedPeerPubKeyHex
		logrus.Infof("This node is set as the allowed peer with ID: %s and PubKey: %s", peerID.String(), allowedPeerPubKeyHex)
	} else {
		logrus.Info("This node is not set as the allowed peer")
	}
	err = cfg.SetupLogging()
	if err != nil {
		logrus.Fatalf("Failed to setup logging: %v", err)
	}
}
