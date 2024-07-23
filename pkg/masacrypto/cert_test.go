package masacrypto

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"
)

func TestGenerateSelfSignedCert(t *testing.T) {
	certPath := "testCert.pem"
	keyPath := "testKey.pem"

	err := GenerateSelfSignedCert(certPath, keyPath)
	if err != nil {
		t.Fatal("[-] Failed to generate self-signed cert:", err)
	}

	// Validate cert file was created
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		t.Fatal("[-] Cert file was not created")
	}

	// Validate key file was created
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		t.Fatal("[-] Key file was not created")
	}

	// Validate cert can be parsed
	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		t.Fatal("[-] Failed to read generated cert:", err)
	}
	block, _ := pem.Decode(certBytes)
	if block == nil {
		t.Fatal("[-] Failed to decode PEM block")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatal("[-] Failed to parse certificate:", err)
	}

	// Validate key can be parsed
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		t.Fatal("[-] Failed to read generated key:", err)
	}
	block, _ = pem.Decode(keyBytes)
	if block == nil {
		t.Fatal("[-] Failed to decode PEM block")
	}
	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		t.Fatal("[-] Failed to parse key:", err)
	}

	// Validate cert and key match
	if cert.PublicKey == key.Public() {
		t.Fatal("[-] Certificate and key do not match")
	}

	// Cleanup test files
	os.Remove(certPath)
	os.Remove(keyPath)
}

func TestGenerateSelfSignedCertErrors(t *testing.T) {
	// Invalid cert path
	err := GenerateSelfSignedCert("/invalid/cert/path", "key.pem")
	if err == nil {
		t.Fatal("[-] Expected error with invalid cert path")
	}

	// Invalid key path
	err = GenerateSelfSignedCert("cert.pem", "/invalid/key/path")
	if err == nil {
		t.Fatal("[-] Expected error with invalid key path")
	}

	err = GenerateSelfSignedCert("cert.pem", "key.pem")
	if err != nil {
		t.Fatal("[-] Expected error when ECDSA key generation fails")
	}
}
