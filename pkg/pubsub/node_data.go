package pubsub

import (
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

const (
	ActivityJoined = iota
	ActivityLeft
)

type NodeData struct {
	Multiaddr         multiaddr.Multiaddr
	PeerId            peer.ID
	LastJoined        time.Time
	LastLeft          time.Time
	LastUpdated       time.Time
	CurrentUptime     time.Duration
	AccumulatedUptime time.Duration
	Activity          int
}

func NewNodeData(multiaddr multiaddr.Multiaddr, peerId peer.ID, activity int) *NodeData {
	return &NodeData{
		PeerId:            peerId,
		Multiaddr:         multiaddr,
		LastJoined:        time.Now(),
		CurrentUptime:     0,
		AccumulatedUptime: 0,
		Activity:          activity,
	}
}

func (n *NodeData) Address() string {
	return fmt.Sprintf("%s/p2p/%s", n.Multiaddr.String(), n.PeerId.String())
}

func (n *NodeData) Joined() {
	now := time.Now()
	n.LastJoined = now
	n.LastUpdated = now
	logrus.Info("Node joined: ", n.Address())
}

func (n *NodeData) Left() {
	logrus.Info("Node left: ", n.Multiaddr.String())
	now := time.Now()
	n.LastLeft = now
	n.LastUpdated = now
	n.AccumulatedUptime += n.GetCurrentUptime()
	n.CurrentUptime = 0
}

func (n *NodeData) GetCurrentUptime() time.Duration {
	return time.Since(n.LastJoined)
}

func (n *NodeData) GetAccumulatedUptime() time.Duration {
	return n.AccumulatedUptime + n.GetCurrentUptime()
}

// UpdateAccumulatedUptime updates the accumulated uptime of the node in the cases where there is a discrepancy between
// the last left and last joined times that came in from the gossip sub events
func (n *NodeData) UpdateAccumulatedUptime() {
	if n.Activity == ActivityLeft {
		n.AccumulatedUptime += n.LastLeft.Sub(n.LastJoined)
	} else {
		n.AccumulatedUptime += time.Since(n.LastJoined)
	}
}
