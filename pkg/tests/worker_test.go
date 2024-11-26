package tests

import (
	"context"
	"os"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/node"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/masa-finance/masa-oracle/pkg/workers"
	datatypes "github.com/masa-finance/masa-oracle/pkg/workers/types"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debug("Log level set to Debug")
}

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
		err := os.Setenv("TWITTER_ACCOUNTS", "test:test")
		if err != nil {
			logrus.Error(err)
			return
		}

		// Start the first node with a random identity
		n1, err := node.NewOracleNode(ctx,
			node.EnableStaked,
			node.EnableRandomIdentity,
			node.IsTwitterScraper,
			node.IsValidator,
			node.UseLocalWorkerAsRemote,
			node.WithPageSize(config.PageSize),
			node.WithOracleProtocol(config.OracleProtocol),
		)
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
		n2, err := node.NewOracleNode(ctx,
			node.EnableStaked,
			node.EnableRandomIdentity,
			node.IsTelegramScraper,
			node.IsValidator,
			node.WithBootNodes(bootNodes...),
			node.WithPageSize(config.PageSize),
			node.WithOracleProtocol(config.OracleProtocol),
		)
		Expect(err).ToNot(HaveOccurred())
		err = n2.Start()
		Expect(err).ToNot(HaveOccurred())

		n2.Host = &MockHost{id: "mockHostID1"}
		oracleNode1 = n1
		oracleNode2 = n2
		category = pubsub.CategoryTwitter
	})

	Describe("GetEligibleWorkers", func() {
		It("should return remote workers and a local worker", func() {
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

			Expect(remoteWorkers).ToNot(BeNil())
			Expect(localWorker).ToNot(BeNil())
		})
	})
})

var _ = Describe("WorkHandlerManager - Local", func() {
	var (
		oracleNode *node.OracleNode
		manager    *workers.WorkHandlerManager
	)

	BeforeEach(func() {
		err := os.Setenv("TWITTER_ACCOUNTS", "test:test")
		if err != nil {
			logrus.Error(err)
			return
		}

		manager = workers.NewWorkHandlerManager(workers.EnableTwitterWorker)
		ctx := context.Background()
		// Start the first node with a random identity
		oracleNode, err = node.NewOracleNode(ctx, node.EnableStaked, node.EnableRandomIdentity, node.IsTwitterScraper, node.WithOracleProtocol(config.OracleProtocol))
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
			if strings.Contains(response.Error, "Twitter authentication failed") {
				logrus.Warn("Passing test as twitter authentication failed")
				return
			} else {
				Expect(response.Error).To(BeEmpty())
			}
		})

		It("should handle errors in work distribution", func() {
			workRequest := datatypes.WorkRequest{
				WorkType: datatypes.WorkerType("InvalidType"),
				Data:     []byte(`{"query": "test", "count": 10}`),
			}
			response := manager.DistributeWork(oracleNode, workRequest)
			Expect(response.Error).ToNot(BeEmpty())
		})
	})
})

var _ = Describe("WorkHandlerManager - Remote", func() {
	var (
		localNode  *node.OracleNode
		remoteNode *node.OracleNode
		manager    *workers.WorkHandlerManager
	)

	BeforeEach(func() {
		err := os.Setenv("TWITTER_ACCOUNTS", "test:test")
		if err != nil {
			logrus.Error(err)
			return
		}

		manager = workers.NewWorkHandlerManager(workers.EnableTwitterWorker)
		ctx := context.Background()

		// Start the first node with a random identity
		workerManagerOptions := []workers.WorkerOptionFunc{workers.EnableTwitterWorker}
		workHandlerManager := workers.NewWorkHandlerManager(workerManagerOptions...)
		protocolOptions := node.WithMasaProtocolHandler(
			config.WorkerProtocol,
			workHandlerManager.HandleWorkerStream,
		)

		remoteNode, err = node.NewOracleNode(ctx,
			protocolOptions,
			node.EnableStaked,
			node.EnableRandomIdentity,
			node.IsTwitterScraper,
			node.IsValidator,
			node.UseLocalWorkerAsRemote,
			node.WithPageSize(config.PageSize),
			node.WithOracleProtocol(config.OracleProtocol),
		)
		Expect(err).ToNot(HaveOccurred())
		err = remoteNode.Start()
		Expect(err).ToNot(HaveOccurred())

		// Get the address of the first node to use as a bootstrap node
		addrs, err := remoteNode.GetP2PMultiAddrs()
		Expect(err).ToNot(HaveOccurred())

		var bootNodes []string
		for _, addr := range addrs {
			bootNodes = append(bootNodes, addr.String())
		}

		// Start the second node with a random identity and bootstrap to the first node
		localNode, err = node.NewOracleNode(ctx,
			protocolOptions,
			node.EnableStaked,
			node.EnableRandomIdentity,
			node.IsTwitterScraper,
			node.IsValidator,
			node.WithBootNodes(bootNodes...),
			node.WithPageSize(config.PageSize),
			node.WithOracleProtocol(config.OracleProtocol),
		)
		Expect(err).ToNot(HaveOccurred())
		err = localNode.Start()
		Expect(err).ToNot(HaveOccurred())

		localNode.Host = &MockHost{id: "mockHostID1"}
	})

	Describe("DistributeWork - remote", func() {
		It("should distribute work to remote nodes", func() {
			// Wait for the nodes to see each other
			Eventually(func() bool {
				nodeDataList := remoteNode.NodeTracker.GetAllNodeData()
				return len(nodeDataList) == 2
			}, "30s").Should(BeTrue())

			workRequest := datatypes.WorkRequest{
				WorkType: datatypes.Twitter,
				Data:     []byte(`{"query": "test", "count": 10}`),
			}
			response := manager.DistributeWork(remoteNode, workRequest)
			if strings.Contains(response.Error, "Twitter authentication failed") {
				logrus.Warn("Passing test as twitter authentication failed")
				return
			} else {
				Expect(response.Error).To(BeEmpty())
			}
		})
	})
})
