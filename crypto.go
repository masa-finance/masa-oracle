package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/sirupsen/logrus"
)

func getOrCreatePrivateKey(keyFile string) (crypto.PrivKey, error) {
	// Check if the private key file exists
	data, err := os.ReadFile(keyFile)
	if err == nil {
		// Decode the private key from the file
		privKey, err := crypto.UnmarshalPrivateKey(data)
		if err != nil {
			logrus.Errorf("Error unmarshalling private key: %s\n", err)
			return nil, err
		}
		logrus.Infof("Loaded private key from %s", keyFile)
		return privKey, nil
	} else {
		// Generate a new private key
		privKey, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
		if err != nil {
			return nil, err
		}
		// Marshal the private key to bytes
		data, err := crypto.MarshalPrivateKey(privKey)
		if err != nil {
			return nil, err
		}
		// Save the private key to the file
		if err := os.WriteFile(keyFile, data, 0600); err != nil {
			return nil, err
		}
		logrus.Infof("Generated and saved a new private key to %s: %s", keyFile, privKey)
		return privKey, nil
	}
}

func generateSelfSignedCert(certPath, keyPath string) error {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	certOut, err := os.Create(certPath)
	if err != nil {
		return err
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return err
	}

	keyOut, err := os.Create(keyPath)
	if err != nil {
		return err
	}
	defer keyOut.Close()

	b, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return err
	}

	if err := pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}); err != nil {
		return err
	}
	return nil
}
