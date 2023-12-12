package pubsub

import (
	"encoding/json"
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

type JSONMultiaddr struct {
	multiaddr.Multiaddr
}

func (m *JSONMultiaddr) UnmarshalJSON(b []byte) error {
	// Unmarshal the JSON as a string
	var multiaddrStr string
	if err := json.Unmarshal(b, &multiaddrStr); err != nil {
		return err
	}

	// Parse the string as a multiaddr
	multiaddr, err := multiaddr.NewMultiaddr(multiaddrStr)
	if err != nil {
		return err
	}

	m.Multiaddr = multiaddr
	return nil
}

type NodeData struct {
	Multiaddrs        []JSONMultiaddr
	PeerId            peer.ID
	LastJoined        time.Time
	LastLeft          time.Time
	LastUpdated       time.Time
	CurrentUptime     time.Duration
	AccumulatedUptime time.Duration
	PublicKey         string
	Activity          int
}

func NewNodeData(addr multiaddr.Multiaddr, peerId peer.ID, activity int) *NodeData {
	multiaddrs := make([]JSONMultiaddr, 0)
	multiaddrs = append(multiaddrs, JSONMultiaddr{addr})

	return &NodeData{
		PeerId:            peerId,
		Multiaddrs:        multiaddrs,
		LastJoined:        time.Now(),
		CurrentUptime:     0,
		AccumulatedUptime: 0,
		PublicKey:         "",
		Activity:          activity,
	}
}

func (n *NodeData) Address() string {
	return fmt.Sprintf("%s/p2p/%s", n.Multiaddrs[0].String(), n.PeerId.String())
}

func (n *NodeData) Joined() {
	now := time.Now()
	n.LastJoined = now
	n.LastUpdated = now
	logrus.Info("Node joined: ", n.Address())
}

func (n *NodeData) Left() {
	logrus.Info("Node left: ", n.Multiaddrs[0].String())
	now := time.Now()
	n.LastLeft = now
	n.LastUpdated = now
	n.AccumulatedUptime += n.GetCurrentUptime()
	n.CurrentUptime = 0
}

func (n *NodeData) GetCurrentUptime() time.Duration {
	var dur time.Duration
	// If the node is currently active, return the time since the last joined time
	if n.Activity == ActivityJoined {
		dur = time.Since(n.LastJoined)
	} else if n.Activity == ActivityLeft {
		dur = 0
	}
	return dur
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
