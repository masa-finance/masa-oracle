package network

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/sirupsen/logrus"
)

// This is an HTTP CONNECT proxy, used for proxying Cosmos blockchain communications over
// Libp2p. It expects the client to send an HTTP CONNECT to the target PeerID (the port
// number is ignored), it then establishes a Libp2p connection to the target host, which
// then connects to its localhost:TargetPort.
//
// To use it:
// * You must run the Cosmos node and the Masa node side by side
// * In the Cosmos node:
//   * Set the environment variable `HTTP_PROXY=127.0.0.1:<proxyListenPort>
//   * Add to the config: `external_address = <Masa node PeerID>:0` (the port number is ignored)
// * In the Masa node:
//   * Set the following configuration parameters:
//       proxy_enabled = true
//       proxyListenAddr = <IP address for the proxy to listen on>
//       proxyListenPort = <Port number for the proxy to listen on>
//       proxyTargetPort = <Port number that the Cosmos node is listening on>

type Proxy struct {
	host       host.Host
	listenAddr net.IP
	listenPort uint16
	targetPort uint16
}

const proxyProtocol = "/connect-proxy/1.0.0"

func NewProxy(host host.Host, listenAddr string, listenPort uint16, targetPort uint16) (*Proxy, error) {
	addr := net.ParseIP(listenAddr)
	if addr == nil {
		return nil, fmt.Errorf("invalid Listen Address: %s", listenAddr)
	}

	return &Proxy{host, addr, listenPort, targetPort}, nil
}

// streamHandler handles the incoming libp2p streams. We know that the stream will contain an
// HTTP request, but strictly speaking we don't care (since CONNECT should act as a transparent
// tunnel).
func (p *Proxy) streamHandler(stream network.Stream) {
	target := fmt.Sprintf("localhost:%d", p.targetPort)
	conn, err := net.Dial("tcp", target)
	if err != nil {
		_ = stream.Reset()
		logrus.Errorf("Error connecting to target host %s: %v", target, err)
	}

	go transfer(conn, stream)
	go transfer(stream, conn)
}

// handleTunnel handles the HTTP CONNECT requests
func (p *Proxy) handleTunnel(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodConnect {
		http.Error(w, "Proxy only supports CONNECT requests", http.StatusBadRequest)
		logrus.Errorf("Received invalid request: %v", *req)
		return
	}

	parts := strings.Split(req.URL.Host, ":")
	peerID, err := peer.IDFromBytes([]byte(parts[0]))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid peerID '%s'", parts[0]), http.StatusBadRequest)
		logrus.Errorf("Invalid PeerID '%s' in host '%s'", parts[0], req.Host)
		return
	}

	destStream, err := p.host.NewStream(ctx, peerID, proxyProtocol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		logrus.Errorf("Error while creating stream to peer %s: %v", peerID, err)
		return
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		logrus.Error("Hijacking not supported")
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		logrus.Errorf("Error while hijacking: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)

	go transfer(destStream, clientConn)
	go transfer(clientConn, destStream)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer closeStream(source)
	defer closeStream(destination)

	if _, err := io.Copy(destination, source); err != nil {
		logrus.Errorf("Error during transfer: %v", err)
	}
}

func closeStream(s io.Closer) {
	err := s.Close()
	if err != nil {
		logrus.Errorf("Error closing stream: %v", err)
	}
}

func (p *Proxy) Start(ctx context.Context) {
	p.host.SetStreamHandler(proxyProtocol, p.streamHandler)

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", p.listenAddr, p.listenPort),
		Handler: http.HandlerFunc(
			func(w http.ResponseWriter, req *http.Request) {
				go p.handleTunnel(ctx, w, req)
			}),
		// Disable HTTP/2.
		// TODO Is this necessary, since we're not doing HTTPS?
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	if err := server.ListenAndServe(); err != nil {
		logrus.Error(err)
	}
}
