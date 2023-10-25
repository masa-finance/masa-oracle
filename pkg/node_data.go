package masa

import (
	"time"

	"github.com/multiformats/go-multiaddr"
)

type NodeData struct {
	Multiaddr         multiaddr.Multiaddr
	LastJoined        time.Time
	LastLeft          time.Time
	LastUpdated       time.Time
	CurrentUptime     time.Duration
	AccumulatedUptime time.Duration
}

func NewNodeData(multiaddr multiaddr.Multiaddr) *NodeData {
	return &NodeData{
		Multiaddr:         multiaddr,
		LastJoined:        time.Now(),
		CurrentUptime:     0,
		AccumulatedUptime: 0,
	}
}

func (n *NodeData) Joined() {
	now := time.Now()
	n.LastJoined = now
	n.LastUpdated = now
}

func (n *NodeData) Left() {
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
