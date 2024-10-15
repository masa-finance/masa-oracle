package tests

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/connmgr"
	"github.com/libp2p/go-libp2p/core/event"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
)

type MockHost struct {
	id peer.ID
}

func (m *MockHost) Peerstore() peerstore.Peerstore {
	fmt.Printf("Peerstore called\n")
	return nil
}

func (m *MockHost) Addrs() []multiaddr.Multiaddr {
	fmt.Printf("Addrs called\n")
	addr1, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/4001")
	return []multiaddr.Multiaddr{addr1}
}

func (m *MockHost) Network() network.Network {
	fmt.Printf("Network called\n")
	return &MockNetwork{}
}

func (m *MockHost) Mux() protocol.Switch {
	fmt.Printf("Mux called\n")
	return nil
}

func (m *MockHost) Connect(ctx context.Context, pi peer.AddrInfo) error {
	fmt.Printf("Connect called with peer info: %v\n", pi)
	if ctx == nil {
		fmt.Printf("nil context\n")
	}
	return nil
}

func (m *MockHost) SetStreamHandler(pid protocol.ID, handler network.StreamHandler) {
	fmt.Printf("SetStreamHandler called with protocol ID: %s\n", pid)
	if handler == nil {
		fmt.Printf("nil handler\n")
	}
}

func (m *MockHost) SetStreamHandlerMatch(id protocol.ID, f func(protocol.ID) bool, handler network.StreamHandler) {
	fmt.Printf("SetStreamHandlerMatch called with protocol ID: %s\n", id)
	if handler == nil {
		fmt.Printf("nil handler\n")
	}
}

func (m *MockHost) RemoveStreamHandler(pid protocol.ID) {
	fmt.Printf("RemoveStreamHandler called with protocol ID: %s\n", pid)
}

func (m *MockHost) NewStream(ctx context.Context, p peer.ID, pids ...protocol.ID) (network.Stream, error) {
	fmt.Printf("NewStream called with peer: %s, protocol IDs: %v\n", p, pids)
	if ctx == nil {
		fmt.Printf("nil context\n")
	}
	return nil, nil
}

func (m *MockHost) Close() error {
	fmt.Printf("Close called\n")
	return nil
}

func (m *MockHost) ConnManager() connmgr.ConnManager {
	fmt.Printf("ConnManager called\n")
	return nil
}

func (m *MockHost) EventBus() event.Bus {
	fmt.Printf("EventBus called\n")
	return nil
}

func (m *MockHost) ID() peer.ID {
	fmt.Printf("ID called\n")
	return m.id
}

type MockNetwork struct{}

func (m *MockNetwork) Close() error {
	fmt.Printf("Close called\n")
	return nil
}

func (m *MockNetwork) CanDial(p peer.ID, addr multiaddr.Multiaddr) bool {
	fmt.Printf("CanDial called with peer: %s, addr: %s\n", p, addr)
	return true
}

func (m *MockNetwork) DialPeer(ctx context.Context, id peer.ID) (network.Conn, error) {
	fmt.Printf("DialPeer called with peer: %s\n", id)
	if ctx == nil {
		fmt.Printf("nil context\n")
	}
	return nil, nil
}

func (m *MockNetwork) SetStreamHandler(handler network.StreamHandler) {
	fmt.Printf("SetStreamHandler called\n")
	if handler == nil {
		fmt.Printf("nil handler\n")
	}
}

func (m *MockNetwork) NewStream(ctx context.Context, id peer.ID) (network.Stream, error) {
	fmt.Printf("NewStream called with peer: %s\n", id)
	if ctx == nil {
		fmt.Printf("nil context\n")
	}
	return nil, nil
}

func (m *MockNetwork) Listen(m2 ...multiaddr.Multiaddr) error {
	fmt.Printf("Listen called with addresses: %v\n", m2)
	return nil
}

func (m *MockNetwork) ResourceManager() network.ResourceManager {
	fmt.Printf("ResourceManager called\n")
	return nil
}

func (m *MockNetwork) Peerstore() peerstore.Peerstore {
	fmt.Printf("Peerstore called\n")
	return nil
}

func (m *MockNetwork) LocalPeer() peer.ID {
	fmt.Printf("LocalPeer called\n")
	return "mockLocalPeerID"
}

func (m *MockNetwork) ListenAddresses() []multiaddr.Multiaddr {
	fmt.Printf("ListenAddresses called\n")
	addr1, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/4001")
	return []multiaddr.Multiaddr{addr1}
}

func (m *MockNetwork) InterfaceListenAddresses() ([]multiaddr.Multiaddr, error) {
	fmt.Printf("InterfaceListenAddresses called\n")
	addr1, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/4001")
	return []multiaddr.Multiaddr{addr1}, nil
}

func (m *MockNetwork) Connectedness(p peer.ID) network.Connectedness {
	fmt.Printf("Connectedness called with peer: %s\n", p)
	return network.NotConnected
}

func (m *MockNetwork) Peers() []peer.ID {
	fmt.Printf("Peers called\n")
	return []peer.ID{}
}

func (m *MockNetwork) Conns() []network.Conn {
	fmt.Printf("Conns called\n")
	return []network.Conn{}
}

func (m *MockNetwork) ConnsToPeer(p peer.ID) []network.Conn {
	fmt.Printf("ConnsToPeer called with peer: %s\n", p)
	return []network.Conn{}
}

func (m *MockNetwork) Notify(notifier network.Notifiee) {
	fmt.Printf("Notify called\n")
	if notifier == nil {
		fmt.Printf("nil notifier\n")
	}
}

func (m *MockNetwork) StopNotify(notifier network.Notifiee) {
	fmt.Printf("StopNotify called\n")
	if notifier == nil {
		fmt.Printf("nil notifier\n")
	}
}

func (m *MockNetwork) ClosePeer(p peer.ID) error {
	fmt.Printf("ClosePeer called with peer: %s\n", p)
	return nil
}

func (m *MockNetwork) RemovePeer(p peer.ID) {
	fmt.Printf("RemovePeer called with peer: %s\n", p)
}
