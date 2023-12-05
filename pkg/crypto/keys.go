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

func GetOrCreatePrivateKey(keyFile string) (privKey crypto.PrivKey, ecdsaPrivKey *ecdsa.PrivateKey, err error) {
	// Check if the private key file is set in the environment
	envKey := os.Getenv("PRIVATE_KEY")
	if envKey != "" {
		rawKey, err := hex.DecodeString(envKey)
		if err != nil {
			logrus.Errorf("Error decoding private key: %s\n", err)
			return nil, nil, err
		}
		privKey, err = crypto.UnmarshalPrivateKey(rawKey)
		if err != nil {
			logrus.Errorf("Error unmarshalling private key: %s\n", err)
			return nil, nil, err
		}
	} else {
		// Check if the private key file exists
		data, err := os.ReadFile(keyFile)
		if err == nil {
			// Decode the private key from the file
			rawKey, err := hex.DecodeString(string(data))
			if err != nil {
				logrus.Errorf("Error decoding private key: %s\n", err)
				return nil, nil, err
			}
			privKey, err = crypto.UnmarshalPrivateKey(rawKey)
			if err != nil {
				logrus.Errorf("Error unmarshalling private key: %s\n", err)
				return nil, nil, err
			}
			logrus.Infof("Loaded private key from %s", keyFile)

		} else {
			// Generate a new private key
			privKey, _, err = crypto.GenerateKeyPair(crypto.Secp256k1, 2048)
			if err != nil {
				return nil, nil, err
			}
			// Marshal the private key to bytes
			data, err := crypto.MarshalPrivateKey(privKey)
			if err != nil {
				return nil, nil, err
			}
			encodedKey := hex.EncodeToString(data)
			// Save the private key to the file
			if err := os.WriteFile(keyFile, []byte(encodedKey), 0600); err != nil {
				return nil, nil, err
			}
			logrus.Infof("Generated and saved a new private key to %s: %s", keyFile, privKey)
		}
	}
	// After obtaining the libp2p privKey, convert it to an ECDSA private key
	ecdsaPrivKey, err = Libp2pPrivateKeyToEcdsa(privKey)
	if err != nil {
		return nil, nil, err
	}
	// Save the ECDSA private key in the same directory as the libp2p key
	ecdsaKeyFilePath := keyFile + ".ecdsa"
	ecdsaKeyBytes := ethCrypto.FromECDSA(ecdsaPrivKey)
	ecdsaKeyHex := hex.EncodeToString(ecdsaKeyBytes)
	if err := os.WriteFile(ecdsaKeyFilePath, []byte(ecdsaKeyHex), 0600); err != nil {
		logrus.Errorf("Error saving ECDSA private key to file: %s\n", err)
		return privKey, nil, err
	}
	logrus.Infof("Saved ECDSA private key to %s", ecdsaKeyFilePath)

	return privKey, ecdsaPrivKey, nil
}

func Libp2pPrivateKeyToEcdsa(privKey crypto.PrivKey) (*ecdsa.PrivateKey, error) {
	// After obtaining the libp2p privKey, convert it to an ECDSA private key
	raw, err := privKey.Raw()
	if err != nil {
		logrus.Errorf("Error getting raw private key: %s\n", err)
		return nil, err
	}

	ecdsaPrivKey, err := ethCrypto.ToECDSA(raw)
	if err != nil {
		logrus.Errorf("Error converting to ECDSA private key: %s\n", err)
		return nil, err
	}

	return ecdsaPrivKey, nil
}

func Libp2pPubKeyToEcdsaHex(pubKey crypto.PubKey) (string, error) {
	ecdsaPublic, err := Libp2pPubKeyToEcdsa(pubKey)
	if err != nil {
		return "", err
	}
	raw, err := x509.MarshalPKIXPublicKey(ecdsaPublic)
	if err != nil {
		return "", err
	}
	ecdsaPubKeyHex := hex.EncodeToString(raw)
	return ecdsaPubKeyHex, nil
}

func Libp2pPubKeyToEcdsa(pubKey crypto.PubKey) (*ecdsa.PublicKey, error) {
	raw, err := pubKey.Raw()
	if err != nil {
		return nil, err
	}

	unmarshalledPub, err := x509.ParsePKIXPublicKey(raw)
	if err != nil {
		return nil, err
	}

	ecdsaPub, ok := unmarshalledPub.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not of type ecdsa")
	}
	return ecdsaPub, nil
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
	ecdsaPrivKey, err := Libp2pPrivateKeyToEcdsa(privKey)
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
