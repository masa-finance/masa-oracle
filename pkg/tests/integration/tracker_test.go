// Nodetracker integration test
package masa_test

import (
	"context"
	"fmt"

	. "github.com/masa-finance/masa-oracle/node"
	"github.com/masa-finance/masa-oracle/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Oracle integration tests", func() {
	Context("NodeData distribution", func() {
		It("is distributed across two nodes", func() {
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			n, err := NewOracleNode(
				ctx,
				config.WithConstantOptions(
					EnableStaked,
					EnableRandomIdentity,
				)...,
			)
			Expect(err).ToNot(HaveOccurred())

			err = n.Start()
			Expect(err).ToNot(HaveOccurred())

			addrs, err := n.GetP2PMultiAddrs()
			Expect(err).ToNot(HaveOccurred())

			var bootNodes []string
			for _, addr := range addrs {
				bootNodes = append(bootNodes, addr.String())
			}

			By(fmt.Sprintf("Generating second node with bootnodes %+v", bootNodes))
			n2, err := NewOracleNode(
				ctx,
				config.WithConstantOptions(
					EnableStaked,
					WithBootNodes(bootNodes...),
					EnableRandomIdentity,
				)...,
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(n.Host.ID()).ToNot(Equal(n2.Host.ID()))

			err = n2.Start()
			Expect(err).ToNot(HaveOccurred())

			// Wait for the nodes to see each others in their respective
			// nodeTracker
			Eventually(func() bool {
				data := n2.NodeTracker.GetAllNodeData()
				return len(data) == 2
			}, "30s").Should(BeTrue())

			Eventually(func() bool {
				data := n.NodeTracker.GetAllNodeData()
				return len(data) == 2
			}, "30s").Should(BeTrue())

			data := n.NodeTracker.GetAllNodeData()

			peerIds := []string{}
			for _, d := range data {
				peerIds = append(peerIds, d.PeerId.String())
			}

			Expect(peerIds).To(ContainElement(n.Host.ID().String()))
			Expect(peerIds).To(ContainElement(n2.Host.ID().String()))
		})
	})
})
