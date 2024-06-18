package chain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"

	"github.com/sirupsen/logrus"
)

type ProofOfStake struct {
	Block  *Block
	Target *big.Int
	Stake  *big.Int
}

func getProofOfStakeTarget(stake *big.Int) *big.Int {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	// adjustment := new(big.Int).Div(target, stake)
	// target.Sub(target, adjustment)
	return target
}

func (pos *ProofOfStake) joinData(timestamp int64) []byte {
	return bytes.Join(
		[][]byte{
			pos.Block.Link,
			pos.Block.Data,
			big.NewInt(timestamp).Bytes(),
			big.NewInt(Difficulty).Bytes(),
			pos.Stake.Bytes(),
		},
		[]byte{},
	)
}

func (pos *ProofOfStake) Run() (int64, []byte) {
	var hash [32]byte
	var hashInt big.Int
	currentTime := time.Now().Unix()

	logrus.WithFields(logrus.Fields{"block_content": string(pos.Block.Data)}).Info("Running Proof of Stake...")
	spinner := []string{"|", "/", "-", "\\"}
	i := 0
	for {
		data := pos.joinData(currentTime)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		fmt.Printf("\r%s %x", spinner[i%len(spinner)], hash)
		i++
		if hashInt.Cmp(pos.Target) == -1 {
			break
		} else {
			currentTime++
		}
	}
	fmt.Println()
	return currentTime, hash[:]
}

func IsValidPoS(block *Block, stake *big.Int) bool {
	var hashIntegerRep big.Int
	pos := &ProofOfStake{Block: block, Target: getProofOfStakeTarget(stake), Stake: stake}
	data := pos.joinData(block.Nonce)
	hash := sha256.Sum256(data)
	hashIntegerRep.SetBytes(hash[:])
	return hashIntegerRep.Cmp(pos.Target) == -1
}
