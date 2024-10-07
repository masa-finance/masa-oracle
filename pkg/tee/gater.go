package tee

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/control"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	ma "github.com/multiformats/go-multiaddr"
)

type RemoteAttestationConnectionGater struct {
	ctx        context.Context
	node       host.Host
	signer     []byte
	production bool
}

func NewRemoteAttestationConnectionGater(ctx context.Context, signer []byte, prod bool) *RemoteAttestationConnectionGater {
	return &RemoteAttestationConnectionGater{
		ctx:        ctx,
		signer:     signer,
		production: prod,
	}
}

func (g *RemoteAttestationConnectionGater) SetNode(node host.Host) {
	g.node = node
}

func (g *RemoteAttestationConnectionGater) checkNode() bool {
	if g.node == nil {
		fmt.Println("RemoteAttestationConnectionGater: node not yet initialized")
		return false
	}

	return true
}

func (g *RemoteAttestationConnectionGater) InterceptPeerDial(p peer.ID) (allow bool) {
	if !g.checkNode() {
		return false
	}
	return VerifyNode(g.ctx, p, g.node, g.signer, g.production) == nil
}

func (g *RemoteAttestationConnectionGater) InterceptAddrDial(p peer.ID, ma ma.Multiaddr) (allow bool) {
	if !g.checkNode() {
		return false
	}
	return VerifyNode(g.ctx, p, g.node, g.signer, g.production) == nil
}

func (g *RemoteAttestationConnectionGater) InterceptAccept(ma network.ConnMultiaddrs) (allow bool) {
	return true
}

func (g *RemoteAttestationConnectionGater) InterceptSecured(dir network.Direction, p peer.ID, ma network.ConnMultiaddrs) (allow bool) {
	if !g.checkNode() {
		return false
	}
	return VerifyNode(g.ctx, p, g.node, g.signer, g.production) == nil
}

func (g *RemoteAttestationConnectionGater) InterceptUpgraded(c network.Conn) (allow bool, reason control.DisconnectReason) {
	return true, 0
}
