package staking_test

import (
	. "github.com/masa-finance/masa-oracle/pkg/staking"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Staking tests", func() {
	Context("LoadContractAddresses", func() {
		It("Returns the contract addresses", func() {
			cont, err := LoadContractAddresses()
			Expect(err).ToNot(HaveOccurred())
			Expect(cont.Sepolia.MasaToken).ToNot(BeEmpty())
			Expect(cont.Sepolia.MasaFaucet).ToNot(BeEmpty())
			Expect(cont.Sepolia.ProtocolStaking).ToNot(BeEmpty())
		})
	})
})
