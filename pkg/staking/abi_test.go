package staking_test

import (
	. "github.com/masa-finance/masa-oracle/pkg/staking"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ABI tests", func() {
	Context("ABI", func() {
		It("Gets the ABI's methods", func() {
			abi, err := GetABI(MasaTokenABIPath)
			Expect(err).ToNot(HaveOccurred())
			Expect(abi.Methods).To(HaveKey("SEND_AND_CALL"))
		})
	})
})
