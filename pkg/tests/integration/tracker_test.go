// Nodetracker integration test
package masa_test

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"

	. "github.com/masa-finance/masa-oracle/node"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/network"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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

	Context("CONNECT proxy", func() {
		It("tunnels the connection", func() {
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			data := []byte("Time is an illusion. Lunchtime, doubly so.\n")

			// Simple echo server
			lis, err := net.Listen("tcp", "127.0.0.1:14242")
			Expect(err).ToNot(HaveOccurred())

			go func(l net.Listener) {
				conn, err := l.Accept()
				Expect(err).ToNot(HaveOccurred())

				buf := make([]byte, len(data))
				_, err = conn.Read(buf)
				Expect(err).ToNot(HaveOccurred())
				Expect(buf).To(Equal(data))

				_, err = conn.Write(buf)
				Expect(err).ToNot(HaveOccurred())

				err = conn.Close()
				Expect(err).ToNot(HaveOccurred())
			}(lis)

			n, n2 := setupNodes(ctx, IsProxy)

			// Create and start the proxies
			p, err := network.NewProxy(n.Host, "127.0.0.1", 24242, 14242)
			Expect(err).ToNot(HaveOccurred())
			go func(ctx context.Context) {
				p.Start(ctx)
			}(ctx)

			p2, err := network.NewProxy(n2.Host, "127.0.0.1", 34242, 14242)
			Expect(err).ToNot(HaveOccurred())
			go func(ctx context.Context) {
				p2.Start(ctx)
			}(ctx)

			// Establish the proxy connection
			target := fmt.Sprintf("%s:0", n2.Host.ID())
			// This is ridiculous but it seems that Go's http.Request makes assumptions that don't work with CONNECT
			// BUT we need the req to properly read the response ¯\_(ツ)_/¯
			rawReq := fmt.Sprintf("CONNECT %s HTTP/1.1\n\n", target)
			req := &http.Request{
				Method: "CONNECT",
				URL:    &url.URL{Host: target},
				Header: make(http.Header),
			}

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

			// Send the CONNECT request, wait for the 200 to indicate that the tunnel is established
			_, err = conn.Write([]byte(rawReq))
			Expect(err).ToNot(HaveOccurred())

			resp, err := http.ReadResponse(bufio.NewReader(conn), req)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			// Now send the data and wait for it to come back
			_, err = conn.Write(data)
			Expect(err).ToNot(HaveOccurred())

			buf := make([]byte, len(data))
			_, err = resp.Body.Read(buf)
			Expect(err).ToNot(HaveOccurred())
			Expect(buf).To(Equal(data))

			resp.Body.Close()
		})
	})
})
