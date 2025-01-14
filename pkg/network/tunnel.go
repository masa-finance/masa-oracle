package network

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/sirupsen/logrus"
)

// Tunnel is a TCP tunnel over Libp2p. It contains a map[string]string, where the key is a port number and the value
// a PeerID.
//
// The tunnel listens on all the given ports (the keys of the map) of the given address. When it receives a connection
// on a given port, it will open a stream to the target peerID, using the protocol configured by the `TunnelProtocol`
// constant.
//
// On the other node, the Tunnel will receive the stream and open a connection to the given targetPort on localhost,
// transparently forwarding all traffic back and forth.
//
// In the Masa node this is used to forward traffic between the local Cosmos node and the Cosmos nodes in all the
// validators.
//
// TODO: Move this to docs or someplace else
//
// To use it:
// * You must run the Cosmos node and the Masa node side by side (in the same host)
// * In the Cosmos nodes:
//   * Set the following in the `[p2p]` section of config.toml with Ignite, it's at `$HOME/.bobtestchain/config`:
//     * For the validators: `external_address = localhost:<port number>` (note that the port number must be unique between all validators)
//     * For all other nodes: `seeds = localhost:<port number>` (for all the port numbers of all the validators)
//     * You might also want to set some `persistent_peers`
// * In the Masa node:
//   * Set the following configuration parameters:
//       tunnel_enabled = true
//       tunnelListenAddr = <IP address for the tunnels to listen on>
//       tunnelPorts = <map of port number -> peerID>
//       tunnelTargetPort = <Port number that the Cosmos node is listening on>
// * You will need to start up the Masa nodes BEFORE the Cosmos nodes, and the bootstrappers before the other nodes

type Tunnel struct {
	host        host.Host
	listenAddr  net.IP
	tunnelPorts map[string]string
	targetPort  uint16
}

const TunnelProtocol = "/masa/tunnel/0.0.1"

func NewTunnel(host host.Host, listenAddr string, tunnelPorts map[string]string, targetPort uint16) (*Tunnel, error) {
	addr := net.ParseIP(listenAddr)
	if addr == nil {
		return nil, fmt.Errorf("invalid Listen Address: %s", listenAddr)
	}

	return &Tunnel{host, addr, tunnelPorts, targetPort}, nil
}

func (t *Tunnel) Start(ctx context.Context) error {
	t.host.SetStreamHandler(TunnelProtocol, t.streamHandler)

	for portStr, peerIdStr := range t.tunnelPorts {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("Invalid port number %s, must be numeric", portStr)
		}
		if port < 1 || port > 65535 {
			return fmt.Errorf("Invalid port number %d, must be in the range 1-65535", port)
		}

		targetPeer, err := peer.Decode(peerIdStr)
		if err != nil {
			return fmt.Errorf("Invalid peerId %s", peerIdStr)
		}

		addr := fmt.Sprintf("%s:%d", t.listenAddr, port)

		go t.listenOn(ctx, addr, targetPeer)

	}

	return nil
}

// streamHandler handles the incoming libp2p streams, forwarding all the data between the stream and a TCP connection
// to the targetPort
func (t *Tunnel) streamHandler(stream network.Stream) {
	logrus.Infof("Received new stream %s from peer %s", stream.ID(), stream.Conn().RemotePeer())
	target := fmt.Sprintf("localhost:%d", t.targetPort)
	conn, err := net.Dial("tcp", target)
	if err != nil {
		_ = stream.Reset()
		logrus.Errorf("Error connecting to target host %s: %v", target, err)
		return
	}

	go transfer(conn, stream)
	go transfer(stream, conn)
}

func (t *Tunnel) handleConnection(ctx context.Context, conn net.Conn, targetPeer peer.ID) {
	logrus.Infof("Received connection from %s", conn.RemoteAddr())

	targetStream, err := t.host.NewStream(ctx, targetPeer, TunnelProtocol)
	if err != nil {
		logrus.Errorf("Error while creating stream to peer %s: %v", targetPeer, err)
		return
	}

	go transfer(conn, targetStream)
	go transfer(targetStream, conn)
}

func (t *Tunnel) listenOn(ctx context.Context, address string, targetPeer peer.ID) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		logrus.Errorf("Error while listening on %s: %#v", address, err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			logrus.Errorf("Error while accepting connection on %s: %#v", address, err)
			return
		}

		go t.handleConnection(ctx, conn, targetPeer)
	}
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
