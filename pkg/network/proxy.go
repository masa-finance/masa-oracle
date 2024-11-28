package network

import (
	"context"
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

const ProxyProtocol = "/masa/connect-proxy/0.0.1"

func NewProxy(host host.Host, listenAddr string, listenPort uint16, targetPort uint16) (*Proxy, error) {
	addr := net.ParseIP(listenAddr)
	if addr == nil {
		return nil, fmt.Errorf("invalid Listen Address: %s", listenAddr)
	}

	return &Proxy{host, addr, listenPort, targetPort}, nil
}

// streamHandler handles the incoming libp2p streams. We know that the stream will contain an
// HTTP request, but strictly speaking we don't care (since CONNECT should act as a transparent
// tunnel), so we just forward the data.
func (p *Proxy) streamHandler(stream network.Stream) {
	logrus.Infof("Received new stream %s from peer %s", stream.ID(), stream.Conn().RemotePeer())
	target := fmt.Sprintf("localhost:%d", p.targetPort)
	conn, err := net.Dial("tcp", target)
	if err != nil {
		_ = stream.Reset()
		logrus.Errorf("Error connecting to target host %s: %v", target, err)
		return
	}

	go transfer(conn, stream)
	go transfer(stream, conn)
}

// handleTunnel handles the HTTP CONNECT requests
func (px *Proxy) handleTunnel(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodConnect {
		http.Error(w, "Proxy only supports CONNECT requests", http.StatusMethodNotAllowed)
		logrus.Errorf("[-] Received invalid request: %#v", *req)
		return
	}

	logrus.Debugf("Received CONNECT request %#v", req)
	parts := strings.Split(req.RequestURI, ":")
	peerID, err := peer.Decode(parts[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid peerID '%s'", parts[0]), http.StatusBadRequest)
		logrus.Errorf("[-] Invalid PeerID '%s' in host '%s'", parts[0], req.Host)
		return
	}

	if peerID == px.host.ID() {
		http.Error(w, fmt.Sprintf("Cannot establish tunnel to myself: %s", peerID), http.StatusBadRequest)
		logrus.Errorf("[-] Tried to establish tunnel to myself")
		return
	}

	logrus.Infof("Creating CONNECT tunnel from %s to peer %s", req.RemoteAddr, peerID)

	destStream, err := px.host.NewStream(ctx, peerID, ProxyProtocol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		logrus.Errorf("Error while creating stream to peer %s: %v", peerID, err)
		return
	}

	logrus.Debug("Stream established, hijacking")
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

	logrus.Debug("Sending response header")
	hdr := fmt.Sprintf("%s 200 OK\n\n", req.Proto)
	_, err = clientConn.Write([]byte(hdr))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Errorf("Error while sending response header: %v", err)
		return
	}

	logrus.Debug("Starting transfer")
	go transfer(clientConn, destStream)
	go transfer(destStream, clientConn)
}

func transfer(dst io.WriteCloser, src io.ReadCloser) {
	defer closeStream(src)
	defer closeStream(dst)

	if _, err := io.Copy(dst, src); err != nil {
		logrus.Errorf("Error during transfer: %v", err)
	}
}

func closeStream(s io.Closer) {
	if err := s.Close(); err != nil {
		logrus.Errorf("Error closing stream: %v", err)
	}
}

func (p *Proxy) Start(ctx context.Context) {
	p.host.SetStreamHandler(ProxyProtocol, p.streamHandler)

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", p.listenAddr, p.listenPort),
		Handler: http.HandlerFunc(
			func(w http.ResponseWriter, req *http.Request) {
				p.handleTunnel(ctx, w, req)
			}),
	}

	if err := server.ListenAndServe(); err != nil {
		logrus.Error(err)
	}
}
