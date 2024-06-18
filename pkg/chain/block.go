package chain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/big"

	"github.com/sirupsen/logrus"
)

type Block struct {
	Data      []byte //	this block's data
	Hash      []byte //	this block's hash
	Link      []byte //	the hash of the last block in the chain. this is the key part that links the blocks together
	Nonce     int64  //	the nonce used to sing the block. useful for verification
	Consensus string // consensus mechanism used: "PoW" or "PoS"
}

func (b *Block) Build(data []byte, link []byte, consensus string, stake *big.Int) {
	b.Data = data
	b.Link = link
	b.Consensus = consensus

	if consensus == "PoW" {
		pow := &ProofOfWork{Block: b, Target: getProofOfWorkTarget()}
		b.Nonce, b.Hash = pow.Run()
	} else if consensus == "PoS" {
		pos := &ProofOfStake{Block: b, Target: getProofOfStakeTarget(stake), Stake: stake}
		b.Nonce, b.Hash = pos.Run()
	}
}

func (b *Block) Serialize() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(b)
	if err != nil {
		logrus.Error("Failed to serialize block: ", b, err)
	}
	return buffer.Bytes(), err
}

func (b *Block) Deserialize(data []byte) error {
	buffer := bytes.Buffer{}
	buffer.Write(data)
	decoder := gob.NewDecoder(&buffer)
	err := decoder.Decode(&b)
	if err != nil {
		logrus.Error("Failed to deserialize data into block: ", data, err)
	}
	return err
}

func (b *Block) Print() {
	fmt.Printf("\t Data:\t%s\n", b.Data)
	fmt.Printf("\t Hash:\t%x\n", b.Hash)
	fmt.Printf("\t Link:\t%x\n", b.Link)
	fmt.Printf("\t Nonce:\t%d\n", b.Nonce)
}
