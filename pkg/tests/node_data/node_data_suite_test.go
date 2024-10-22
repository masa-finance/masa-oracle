package node_data

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestNodeData(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NodeData Test Suite")
}
