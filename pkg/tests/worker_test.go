package tests

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/masa-oracle/node"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/masa-finance/masa-oracle/pkg/workers"
	datatypes "github.com/masa-finance/masa-oracle/pkg/workers/types"
)

func TestWorkers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Workers Suite")
}

var _ = Describe("Worker Selection", func() {
	var (
		oracleNode1 *node.OracleNode
		oracleNode2 *node.OracleNode
		category    pubsub.WorkerCategory
	)

	BeforeEach(func() {
		ctx := context.Background()

		// Start the first node with a random identity
		n1, err := node.NewOracleNode(ctx, node.EnableStaked, node.EnableRandomIdentity, node.IsTwitterScraper)
		Expect(err).ToNot(HaveOccurred())
		err = n1.Start()
		Expect(err).ToNot(HaveOccurred())

		// Get the address of the first node to use as a bootstrap node
		addrs, err := n1.GetP2PMultiAddrs()
		Expect(err).ToNot(HaveOccurred())

		var bootNodes []string
		for _, addr := range addrs {
			bootNodes = append(bootNodes, addr.String())
		}

		// Start the second node with a random identity and bootstrap to the first node
		n2, err := node.NewOracleNode(ctx, node.EnableStaked, node.EnableRandomIdentity, node.IsTelegramScraper, node.WithBootNodes(bootNodes...))
		Expect(err).ToNot(HaveOccurred())
		err = n2.Start()
		Expect(err).ToNot(HaveOccurred())

		n2.Host = &MockHost{id: "mockHostID1"}
		oracleNode1 = n1
		oracleNode2 = n2
		category = pubsub.CategoryTwitter
	})

	AfterEach(func() {
		//oracleNode1.Stop()
		//oracleNode2.Stop()
	})

	Describe("GetEligibleWorkers", func() {
		It("should return empty remote workers and a local worker", func() {
			// Wait for the nodes to see each other
			Eventually(func() bool {
				datas := oracleNode1.NodeTracker.GetAllNodeData()
				return len(datas) == 2
			}, "30s").Should(BeTrue())

			Eventually(func() bool {
				datas := oracleNode2.NodeTracker.GetAllNodeData()
				return len(datas) == 2
			}, "30s").Should(BeTrue())

			remoteWorkers, localWorker := workers.GetEligibleWorkers(oracleNode1, category, 1)

			Expect(remoteWorkers).To(BeEmpty())
			Expect(localWorker).ToNot(BeNil())
		})
	})
})

var _ = Describe("WorkHandlerManager", func() {
	var (
		oracleNode *node.OracleNode
		manager    *workers.WorkHandlerManager
	)

	BeforeEach(func() {
		manager = workers.NewWorkHandlerManager(workers.EnableTwitterWorker)
		ctx := context.Background()
		var err error
		// Start the first node with a random identity
		oracleNode, err = node.NewOracleNode(ctx, node.EnableStaked, node.EnableRandomIdentity, node.IsTwitterScraper)
		Expect(err).ToNot(HaveOccurred())
		err = oracleNode.Start()
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("Add and Get WorkHandler", func() {
		It("should add and retrieve a work handler", func() {
			handler, exists := manager.GetWorkHandler(datatypes.Twitter)
			Expect(exists).To(BeTrue())
			Expect(handler).ToNot(BeNil())
		})

		It("should return false for non-existent work handler", func() {
			_, exists := manager.GetWorkHandler(datatypes.WorkerType("NonExistent"))
			Expect(exists).To(BeFalse())
		})
	})

	Describe("DistributeWork", func() {
		It("should distribute work to eligible workers", func() {
			workRequest := datatypes.WorkRequest{
				WorkType: datatypes.Twitter,
				Data:     []byte(`{"query": "test", "count": 10}`),
			}
			response := manager.DistributeWork(oracleNode, workRequest)
			Expect(response.Error).To(BeEmpty())
		})

		It("should handle errors in work distribution", func() {
			workRequest := datatypes.WorkRequest{
				WorkType: datatypes.WorkerType("InvalidType"),
				Data:     []byte(`{"query": "test", "count": 10}`),
			}
			response := manager.DistributeWork(nil, workRequest)
			Expect(response.Error).ToNot(BeEmpty())
		})
	})
})
