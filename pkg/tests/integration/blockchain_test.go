// Blockchain integration test
package masa_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	. "github.com/masa-finance/masa-oracle/node"
	"github.com/masa-finance/masa-oracle/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Blockchain tests", func() {
	Context("blockchain events", func() {
		It("contains data published by nodes and propagates correctly", func() {
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			tempDir, err := os.MkdirTemp("", "blockchain")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tempDir)

			blockChainEventTracker := NewBlockChain()

			// We create two nodes in this test - one publishes data over the blockchain
			// and the other one should be able to receive the data
			n, err := NewOracleNode(
				ctx,
				EnableStaked,
				EnableRandomIdentity,
				WithPubSubHandler(config.BlockTopic, blockChainEventTracker, true),
				WithService(blockChainEventTracker.Start(tempDir)),
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

			tempDir2, err := os.MkdirTemp("", "blockchain")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tempDir)

			// We start the second node (receives the data)
			// And we also set the first node as bootstrap node
			blockChainEventTracker2 := NewBlockChain()
			n2, err := NewOracleNode(ctx,
				EnableStaked,
				WithBootNodes(bootNodes...),
				EnableRandomIdentity,
				WithPubSubHandler(config.BlockTopic, blockChainEventTracker2, true),
				WithService(blockChainEventTracker2.Start(tempDir2)),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(n.Host.ID()).ToNot(Equal(n2.Host.ID()))

			err = n2.Start()
			Expect(err).ToNot(HaveOccurred())

			// Initially, there shouldn't be any events in the blockchain
			Eventually(func() int {
				return len(blockChainEventTracker2.BlockEvents)
			}, "30s").Should(Equal(0))

			publishedData := map[string]interface{}{
				"foo": interface{}("bar"),
			}

			publishBytes, err := json.Marshal(publishedData)
			Expect(err).ToNot(HaveOccurred())

			// Publish data with the first node to kick off the first event
			err = n.PublishTopic(config.BlockTopic, publishBytes)
			Expect(err).ToNot(HaveOccurred())

			// Eventually we should have at least one event
			Eventually(func() int {
				return len(blockChainEventTracker2.BlockEvents)
			}, "30s").ShouldNot(Equal(0))

			// Check that the event has the data that we published
			Expect(blockChainEventTracker2.BlockEvents[0]).To(Equal(BlockEvents(publishedData)))
		})
	})
})
