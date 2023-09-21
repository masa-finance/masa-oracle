package ethereum

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	libp2pCrypto "github.com/libp2p/go-libp2p/core/crypto"
)

func LibP2pToEcdsa(key libp2pCrypto.PrivKey) (*ecdsa.PrivateKey, error) {
	raw, err := key.Raw()
	if err != nil {
		return nil, err
	}
	ecdsaPrivKey, err := crypto.ToECDSA(raw)
	if err != nil {
		return nil, err
	}
	return ecdsaPrivKey, nil
}
