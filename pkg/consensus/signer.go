package consensus

// This package provides cryptographic functionalities for signing and verifying data within the consensus mechanism of the system.
// It leverages the go-libp2p core crypto library to handle cryptographic operations, ensuring secure data handling.
// - SignData: Signs arbitrary data using a private key, ensuring the data integrity and source authenticity.
// - VerifySignature: Validates the signature of the data against a public key, confirming the data's integrity and origin.

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/crypto"
	mh "github.com/multiformats/go-multihash"
	"github.com/sirupsen/logrus"
)

func GenerateJWTToken(peerId string) (string, error) {
	// Set the expiration time for the token (e.g., 1 hour from now)
	expirationTime := time.Now().Add(1 * time.Hour)

	mhHash, err := mh.Sum([]byte(peerId), mh.SHA2_256, -1)
	if err != nil {
		return "", err
	}
	apiKey := cid.NewCidV1(cid.Raw, mhHash).String()

	// Create the JWT claims
	claims := jwt.MapClaims{
		"apiKey": apiKey,
		"exp":    expirationTime.Unix(),
	}

	// Create a new JWT token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	secretKey := []byte(peerId)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// SignData signs the data using the provided private key and returns the signature.
func SignData(privKey crypto.PrivKey, data []byte) ([]byte, error) {
	if privKey == nil {
		logrus.Error("Private key is nil")
		return nil, fmt.Errorf("private key is nil")
	}

	logrus.Info("Signing data")

	signature, err := privKey.Sign(data)
	if err != nil {
		logrus.WithError(err).Error("Failed to sign data")
		return nil, err
	}

	logrus.Info("Data signed successfully")
	return signature, nil
}

// VerifySignature verifies the signature of the data using the signers public key, the data that was signed, and the signature.
func VerifySignature(pubKey crypto.PubKey, data []byte, signatureHex string) (bool, error) {
	if pubKey == nil {
		logrus.Error("Public key is nil")
		return false, fmt.Errorf("public key is nil")
	}

	// Decode the hexadecimal-encoded signature back to its original byte format
	signatureBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		logrus.WithError(err).Error("Failed to decode signature from hexadecimal")
		return false, err
	}

	logrus.Info("Verifying signature")

	verified, err := pubKey.Verify(data, signatureBytes)
	if err != nil {
		logrus.WithError(err).Error("Failed to verify signature")
		return false, err
	}

	if verified {
		logrus.Info("Signature verified successfully")
	} else {
		logrus.Info("Signature verification failed")
	}
	return verified, nil
}
