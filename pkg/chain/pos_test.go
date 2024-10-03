package chain_test

import (
	"math/big"

	. "github.com/masa-finance/masa-oracle/pkg/chain"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProofOfStake", func() {
	var (
		block *Block
		stake *big.Int
		pos   *ProofOfStake
	)

	BeforeEach(func() {
		block = &Block{
			Link: []byte("previous hash"),
			Data: []byte("block data"),
		}
		stake = big.NewInt(1000)
		pos = &ProofOfStake{
			Block:  block,
			Stake:  stake,
			Target: GetProofOfStakeTarget(stake),
		}
	})

	Describe("Run", func() {
		It("finds a valid nonce and hash", func() {
			nonce, hash := pos.Run()
			Expect(nonce).To(BeNumerically(">", 0))
			Expect(hash).NotTo(BeNil())

			var hashInt big.Int
			hashInt.SetBytes(hash)
			Expect(hashInt.Cmp(pos.Target)).To(Equal(-1))
		})
	})

	Describe("IsValidPoS", func() {
		It("validates the block's proof of stake", func() {
			nonce, _ := pos.Run()
			block.Nonce = nonce

			isValid := IsValidPoS(block, stake)
			Expect(isValid).To(BeTrue())
		})
	})
})
