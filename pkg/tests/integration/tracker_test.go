// Nodetracker integration test
package masa_test

import (
	"context"
	"fmt"
	"net"

	. "github.com/masa-finance/masa-oracle/node"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/network"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

// setupNodes sets up two Masa nodes and connects them together
// addOpts are additional options to add to the node for the specific test
func setupNodes(ctx context.Context, addOpts ...Option) (*OracleNode, *OracleNode) {
	opts := []Option{EnableStaked, EnableRandomIdentity}
	opts = append(opts, addOpts...)

	n, err := NewOracleNode(
		ctx,
		config.WithConstantOptions(opts...)...,
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
	opts = append(opts, WithBootNodes(bootNodes...))

	n2, err := NewOracleNode(
		ctx,
		config.WithConstantOptions(opts...)...,
	)
	Expect(err).ToNot(HaveOccurred())
	Expect(n.Host.ID()).ToNot(Equal(n2.Host.ID()))

	err = n2.Start()
	Expect(err).ToNot(HaveOccurred())

	// Wait for the nodes to see each other in their respective nodeTracker
	Eventually(func() bool {
		datas := n2.NodeTracker.GetAllNodeData()
		return len(datas) == 2
	}, "30s").Should(BeTrue())

	Eventually(func() bool {
		datas := n.NodeTracker.GetAllNodeData()
		return len(datas) == 2
	}, "30s").Should(BeTrue())

	return n, n2
}

var _ = Describe("Oracle integration tests", func() {
	Context("NodeData distribution", func() {
		It("is distributed across two nodes", func() {
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			n, n2 := setupNodes(ctx)

			data := n.NodeTracker.GetAllNodeData()

			peerIds := []string{}
			for _, d := range data {
				peerIds = append(peerIds, d.PeerId.String())
			}

			Expect(peerIds).To(ContainElement(n.Host.ID().String()))
			Expect(peerIds).To(ContainElement(n2.Host.ID().String()))
		})
	})

	Context("tunnel", func() {
		It("tunnels the connection", func() {
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			inData := []byte("Time is an illusion.\n")
			outData := []byte("Lunchtime, doubly so.\n")

			// Simple server
			lis, err := net.Listen("tcp", "127.0.0.1:14242")
			Expect(err).ToNot(HaveOccurred())

			go func(l net.Listener) {
				conn, err := l.Accept()
				Expect(err).ToNot(HaveOccurred())

				buf := make([]byte, len(inData))
				_, err = conn.Read(buf)
				Expect(err).ToNot(HaveOccurred())
				Expect(buf).To(Equal(inData))

				_, err = conn.Write(outData)
				Expect(err).ToNot(HaveOccurred())

				err = conn.Close()
				Expect(err).ToNot(HaveOccurred())
			}(lis)

			n, n2 := setupNodes(ctx, EnableTunnel)

			tunnelMap := map[string]string{
				"24242": n.Host.ID().String(),
				"34242": n2.Host.ID().String(),
			}
			// Create and start the tunnels
			p, err := network.NewTunnel(n.Host, "127.0.0.1", tunnelMap, 14242)
			Expect(err).ToNot(HaveOccurred())
			go func(ctx context.Context) {
				err := p.Start(ctx)
				if err != nil {
					logrus.Errorf("Error when starting tunnel: %#v", err)
				}
			}(ctx)

			p2, err := network.NewTunnel(n2.Host, "127.0.0.1", tunnelMap, 14242)
			Expect(err).ToNot(HaveOccurred())
			go func(ctx context.Context) {
				err := p2.Start(ctx)
				if err != nil {
					logrus.Errorf("Error when starting tunnel: %#v", err)
				}
			}(ctx)

			// Wait until the proxy is listening
			var conn net.Conn
			for {
				conn, err = net.Dial("tcp", "127.0.0.1:24242")
				if err == nil {
					break
				} else {
					Expect(err.Error()).To(ContainSubstring("connection refused"))
				}
			}

			// Send the data and wait for it to come back
			_, err = conn.Write(inData)
			Expect(err).ToNot(HaveOccurred())

			buf := make([]byte, len(outData))
			_, err = conn.Read(buf)
			Expect(err).ToNot(HaveOccurred())
			Expect(buf).To(Equal(outData))

			conn.Close()
		})
	})
})
