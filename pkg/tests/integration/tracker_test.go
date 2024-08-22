// Nodetracker integration test
package masa_test

import (
	"context"
	"os"

	. "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/masacrypto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Oracle integration tests", func() {
	Context("NodeData", func() {
		It("is distributed across two nodes", func() {
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			// Create a new temp file
			dir, err := os.MkdirTemp("", "config")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(dir)

			err = masacrypto.GenerateSelfSignedCert(dir+"/node1.cert", dir+"/node1.key")
			Expect(err).ToNot(HaveOccurred())

			// write now a config
			generateNodeKeys(dir+"/config.yaml", "node1")
			n2, err := NewOracleNode(ctx, config.EnableStaked, config.DisableCLIParse)
			Expect(err).ToNot(HaveOccurred())

			generateNodeKeys(dir+"/config.yaml", "node2")
			n, err := NewOracleNode(ctx, config.EnableStaked, config.DisableCLIParse)
			Expect(err).ToNot(HaveOccurred())

			Expect(n).ToNot(BeNil())

			err = n.Start()
			Expect(err).ToNot(HaveOccurred())

			err = n2.Start()
			Expect(err).ToNot(HaveOccurred())

			// Wait for the nodes to start
			Eventually(func() bool {
				datas := n.NodeTracker.GetAllNodeData()
				return len(datas) > 0
			}).Should(BeTrue())
		})
	})
})
