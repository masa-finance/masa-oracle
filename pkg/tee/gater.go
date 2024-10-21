package tee

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/control"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	ma "github.com/multiformats/go-multiaddr"
)

type RemoteAttestationConnectionGater struct {
	sync.Mutex
	ctx            context.Context
	node           host.Host
	signer         []byte
	production     bool
	verifyingPeers map[peer.ID]struct{}
	badPeers       map[peer.ID]struct{}
	goodPeers      map[peer.ID]struct{}
	cleanupTime    time.Duration
}

func NewRemoteAttestationConnectionGater(ctx context.Context, signer []byte, prod bool) *RemoteAttestationConnectionGater {
	return &RemoteAttestationConnectionGater{
		ctx:            ctx,
		signer:         signer,
		production:     prod,
		verifyingPeers: make(map[peer.ID]struct{}),
		cleanupTime:    5 * time.Minute,
	}
}

func (g *RemoteAttestationConnectionGater) SetNode(node host.Host) {
	g.node = node
}

func (g *RemoteAttestationConnectionGater) checkNodeIsStarted() bool {
	if g.node == nil {
		fmt.Println("RemoteAttestationConnectionGater: node not yet initialized")
		return false
	}

	return true
}

func (g *RemoteAttestationConnectionGater) checkIsGoodPeer(p peer.ID) bool {
	g.Lock()
	defer g.Unlock()

	_, isGoodPeer := g.goodPeers[p]
	return isGoodPeer
}

func (g *RemoteAttestationConnectionGater) checkIsBadPeer(p peer.ID) bool {
	g.Lock()
	defer g.Unlock()

	_, isBadPeer := g.badPeers[p]
	return isBadPeer
}

func (g *RemoteAttestationConnectionGater) isPeerBeingVerified(p peer.ID) bool {
	g.Lock()
	defer g.Unlock()

	// if the peer is already dialing, return true
	_, dialing := g.verifyingPeers[p]
	if dialing {
		return true
	}

	g.verifyingPeers[p] = struct{}{}

	return false
}

func (g *RemoteAttestationConnectionGater) cleanupVerification(p peer.ID) {
	g.Lock()
	defer g.Unlock()
	delete(g.verifyingPeers, p)
}

func (g *RemoteAttestationConnectionGater) verifyNode(p peer.ID) bool {
	// verify if the node is allowed to connect by challenging it,
	// if doesn't pass the challenge, add it to the badPeers list
	// so we don't run the challenging again on the same node for a certain time

	allowed := VerifyNode(g.ctx, p, g.node, g.signer, g.production) == nil
	if !allowed {
		// flag the node as bad peer
		g.Lock()
		g.badPeers[p] = struct{}{}
		g.Unlock()

		// Remove the peer from the badPeers list after a certain time
		// This allows e.g. a node which was misconfigured to not be "banned"
		// by the network forever
		go func() {
			time.Sleep(g.cleanupTime)
			g.Lock()
			delete(g.badPeers, p)
			g.Unlock()
		}()
	}

	g.Lock()
	defer g.Unlock()
	g.goodPeers[p] = struct{}{}

	return allowed
}

func (g *RemoteAttestationConnectionGater) waitPeerVerification(p peer.ID) error {
	maxAttempts := 10

	for {
		if maxAttempts == 0 {
			return fmt.Errorf("peer %s is still dialing", p)
		}
		if !g.isPeerBeingVerified(p) {
			return nil
		}
		time.Sleep(1 * time.Second)
		maxAttempts--
	}
}

func (g *RemoteAttestationConnectionGater) checkNodeValidity(p peer.ID) bool {
	if !g.checkNodeIsStarted() || g.checkIsBadPeer(p) {
		return false
	}

	if g.checkIsGoodPeer(p) {
		return true
	}

	// Until we verify a peer, any other request is on "hold"
	if err := g.waitPeerVerification(p); err != nil {
		fmt.Println("Failed to wait for peer to finish dialing:", err)
		return false
	}

	// Once done, remove the peer from the dialing map
	defer g.cleanupVerification(p)

	return g.verifyNode(p)
}

func (g *RemoteAttestationConnectionGater) InterceptPeerDial(p peer.ID) (allow bool) {
	return g.checkNodeValidity(p)
}

func (g *RemoteAttestationConnectionGater) InterceptAddrDial(p peer.ID, ma ma.Multiaddr) (allow bool) {
	return g.checkNodeValidity(p)
}

func (g *RemoteAttestationConnectionGater) InterceptAccept(ma network.ConnMultiaddrs) (allow bool) {
	return true
}

func (g *RemoteAttestationConnectionGater) InterceptSecured(dir network.Direction, p peer.ID, ma network.ConnMultiaddrs) (allow bool) {
	// Let's not verify outbound connections.
	// Peers should verify us in this case.
	if dir == network.DirOutbound {
		return true
	}

	return g.checkNodeValidity(p)
}

func (g *RemoteAttestationConnectionGater) InterceptUpgraded(c network.Conn) (allow bool, reason control.DisconnectReason) {
	return true, 0
}
