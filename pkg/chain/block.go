package chain

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"math/big"

	"github.com/sirupsen/logrus"
)

type Block struct {
	Data  []byte //	this block's data
	Hash  []byte //	this block's hash
	Link  []byte //	the hash of the last block in the chain
	Nonce int64  //	the nonce used to sign the block for verification
}

func (b *Block) Build(data []byte, link []byte, stake *big.Int) {
	b.Data = data
	b.Link = link

	pos := &ProofOfStake{Block: b, Target: getProofOfStakeTarget(stake), Stake: stake}
	b.Nonce, b.Hash = pos.Run()
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
	fmt.Printf("\t Input Data:    \t%s\n", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", b.Data))))
	fmt.Printf("\t Transaction Hash:\t%x\n", b.Hash)
	fmt.Printf("\t Previous Hash:  \t%x\n", b.Link)
	fmt.Printf("\t Transaction Nonce:\t%d\n", b.Nonce)
}
