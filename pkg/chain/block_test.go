package chain_test

import (
	"math/big"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/masa-finance/masa-oracle/pkg/chain"
)

var _ = Describe("Block", func() {
	var (
		block       *Block
		stake       *big.Int
		data        []byte
		link        []byte
		blockNumber uint64
	)

	BeforeEach(func() {
		block = &Block{}
		stake = big.NewInt(1)
		data = []byte("sample block data")
		link = []byte("previous block hash")
		blockNumber = 0
	})

	Describe("Build", func() {
		It("should build the block with correct data and hash", func() {
			block.Build(data, link, stake, blockNumber)

			Expect(block.Block).To(Equal(blockNumber))
			Expect(block.Data).To(Equal(data))
			Expect(block.Link).To(Equal(link))
			Expect(block.Hash).NotTo(BeNil())
			Expect(block.Nonce).NotTo(BeZero())
		})
	})

	Describe("Serialize", func() {
		It("should serialize the block without error", func() {
			block.Build(data, link, stake, blockNumber)

			serializedData, err := block.Serialize()
			Expect(err).To(BeNil())
			Expect(serializedData).NotTo(BeNil())

			// Ensure serialized data can be deserialized correctly
			newBlock := &Block{}
			err = newBlock.Deserialize(serializedData)
			Expect(err).To(BeNil())
			Expect(newBlock.Block).To(Equal(block.Block))
			Expect(newBlock.Data).To(Equal(block.Data))
			Expect(newBlock.Link).To(Equal(block.Link))
			Expect(newBlock.Hash).To(Equal(block.Hash))
			Expect(newBlock.Nonce).To(Equal(block.Nonce))
		})
	})

	Describe("Deserialize", func() {
		It("should deserialize correctly from serialized data", func() {
			block.Build(data, link, stake, blockNumber)
			serializedData, err := block.Serialize()
			Expect(err).To(BeNil())

			// Create a new block and deserialize from serialized data
			newBlock := &Block{}
			err = newBlock.Deserialize(serializedData)
			Expect(err).To(BeNil())

			// Compare fields between original and deserialized block
			Expect(newBlock.Block).To(Equal(block.Block))
			Expect(newBlock.Data).To(Equal(block.Data))
			Expect(newBlock.Link).To(Equal(block.Link))
			Expect(newBlock.Hash).To(Equal(block.Hash))
			Expect(newBlock.Nonce).To(Equal(block.Nonce))
		})

		It("should return an error when deserializing invalid data", func() {
			invalidData := []byte("invalid serialized data")
			err := block.Deserialize(invalidData)
			Expect(err).NotTo(BeNil())
		})
	})
})
