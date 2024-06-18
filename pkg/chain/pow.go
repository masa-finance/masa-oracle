package chain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"

	"github.com/sirupsen/logrus"
)

/*
Implementation of the difficulty rate
*/
const Difficulty = int64(21)

type ProofOfWork struct {
	Block  *Block   // a block from the blockchain
	Target *big.Int // number that represents the requirements. it's derived form the difficulty
}

func getProofOfWorkTarget() *big.Int {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	return target
}

func (pow *ProofOfWork) joinData(nonce int64) []byte {
	return bytes.Join(
		[][]byte{
			pow.Block.Link,
			pow.Block.Data,
			big.NewInt(nonce).Bytes(),
			big.NewInt(Difficulty).Bytes(),
		},
		[]byte{},
	)
}

func (pow *ProofOfWork) Run() (int64, []byte) {
	var hash [32]byte
	var hashIntegerRep big.Int
	nonce := int64(0)

	logrus.WithFields(logrus.Fields{"block_content": string(pow.Block.Data)}).Info("Running Proof of Work...")
	for nonce < math.MaxInt64 {

		data := pow.joinData(nonce)
		hash = sha256.Sum256(data)
		hashIntegerRep.SetBytes(hash[:])
		fmt.Printf("\r%x", hash)

		if hashIntegerRep.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()
	return nonce, hash[:]
}

func IsValid(block *Block) bool {
	var hashIntegerRep big.Int
	pow := &ProofOfWork{Block: block, Target: getProofOfWorkTarget()}
	data := pow.joinData(block.Nonce)
	hash := sha256.Sum256(data)
	hashIntegerRep.SetBytes(hash[:])
	return hashIntegerRep.Cmp(pow.Target) == -1
}
