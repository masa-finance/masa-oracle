package staking_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStaking(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Staking test suite")
}
