package masa_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOracle(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Oracle integration test suite")
}
