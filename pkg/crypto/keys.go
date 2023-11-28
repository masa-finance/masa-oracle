package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/sirupsen/logrus"
)

func GetOrCreatePrivateKey(keyFile string) (privKey crypto.PrivKey, err error) {
	// Check if the private key file is set in the environment
	envKey := os.Getenv("PRIVATE_KEY")
	if envKey != "" {
		rawKey, err := hex.DecodeString(envKey)
		if err != nil {
			logrus.Errorf("Error decoding private key: %s\n", err)
			return nil, err
		}
		privKey, err = crypto.UnmarshalPrivateKey(rawKey)
		if err != nil {
			logrus.Errorf("Error unmarshalling private key: %s\n", err)
			return nil, err
		}
	} else {
		// Check if the private key file exists
		data, err := os.ReadFile(keyFile)
		if err == nil {
			// Decode the private key from the file
			rawKey, err := hex.DecodeString(string(data))
			if err != nil {
				logrus.Errorf("Error decoding private key: %s\n", err)
				return nil, err
			}
			privKey, err = crypto.UnmarshalPrivateKey(rawKey)
			if err != nil {
				logrus.Errorf("Error unmarshalling private key: %s\n", err)
				return nil, err
			}
			logrus.Infof("Loaded private key from %s", keyFile)

		} else {
			// Generate a new private key
			privKey, _, err = crypto.GenerateKeyPair(crypto.Secp256k1, 2048)
			if err != nil {
				return nil, err
			}
			// Marshal the private key to bytes
			data, err := crypto.MarshalPrivateKey(privKey)
			if err != nil {
				return nil, err
			}
			encodedKey := hex.EncodeToString(data)
			// Save the private key to the file
			if err := os.WriteFile(keyFile, []byte(encodedKey), 0600); err != nil {
				return nil, err
			}
			logrus.Infof("Generated and saved a new private key to %s: %s", keyFile, privKey)
		}
	}
	return privKey, nil
}

func GenerateSelfSignedCert(certPath, keyPath string) error {
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

func VerifyEthereumCompatibility(privKey crypto.PrivKey) (string, error) {
	// Convert the libp2p private key to an Ethereum private key
	raw, err := privKey.Raw()
	if err != nil {
		return "", err
	}
	ecdsaPrivKey, err := ethCrypto.ToECDSA(raw)
	if err != nil {
		return "", err
	}

	ethAddress := ethCrypto.PubkeyToAddress(ecdsaPrivKey.PublicKey).Hex()

	// Print the private key in hexadecimal format
	data, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		return "", err
	}
	fmt.Printf("Private key: \n%s\n", hex.EncodeToString(data))

	// Print the public key and address
	fmt.Println("Ethereum public key:", ecdsaPrivKey.PublicKey)
	// Derive the Ethereum address from the private key
	fmt.Println("Ethereum address:", ethCrypto.PubkeyToAddress(ecdsaPrivKey.PublicKey).Hex())

	return ethAddress, nil
}
