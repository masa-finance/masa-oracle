package pubsub

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"strconv"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/masacrypto"
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

	// Parse the string as a ma
	ma, err := multiaddr.NewMultiaddr(multiaddrStr)
	if err != nil {
		return err
	}

	m.Multiaddr = ma
	return nil
}

type NodeData struct {
	Multiaddrs           []JSONMultiaddr `json:"multiaddrs,omitempty"`
	PeerId               peer.ID         `json:"peerId"`
	LastJoined           time.Time       `json:"lastJoined,omitempty"`
	LastLeft             time.Time       `json:"lastLeft,omitempty"`
	LastUpdated          time.Time       `json:"lastUpdated,omitempty"`
	CurrentUptime        time.Duration   `json:"currentUptime,omitempty"`
	CurrentUptimeStr     string          `json:"readableCurrentUptime,omitempty"`
	AccumulatedUptime    time.Duration   `json:"accumulatedUptime,omitempty"`
	AccumulatedUptimeStr string          `json:"readableAccumulatedUptime,omitempty"`
	EthAddress           string          `json:"ethAddress,omitempty"`
	Activity             int             `json:"activity,omitempty"`
	IsActive             bool            `json:"isActive"`
	IsStaked             bool            `json:"isStaked"`
	SelfIdentified       bool            `json:"-"`
	IsWriterNode         bool            `json:"isWriterNode"`
}

func NewNodeData(addr multiaddr.Multiaddr, peerId peer.ID, publicKey string, activity int) *NodeData {
	multiaddrs := make([]JSONMultiaddr, 0)
	multiaddrs = append(multiaddrs, JSONMultiaddr{addr})
	wn, _ := strconv.ParseBool(viper.GetString("WRITER_NODE"))

	return &NodeData{
		PeerId:            peerId,
		Multiaddrs:        multiaddrs,
		LastUpdated:       time.Now(),
		CurrentUptime:     0,
		AccumulatedUptime: 0,
		EthAddress:        publicKey,
		Activity:          activity,
		SelfIdentified:    false,
		IsWriterNode:      wn,
	}
}

func (n *NodeData) Address() string {
	return fmt.Sprintf("%s/p2p/%s", n.Multiaddrs[0].String(), n.PeerId.String())
}

func (n *NodeData) Joined() {
	now := time.Now()
	n.LastJoined = now
	n.LastUpdated = now
	n.Activity = ActivityJoined
	n.IsActive = true
	if n.IsStaked {
		logrus.Info("Node joined: ", n.Address())
	} else {
		logrus.Debug("Node joined: ", n.Address())
	}
}

func (n *NodeData) Left() {
	if n.Activity == ActivityLeft {
		if n.IsStaked {
			logrus.Warnf("Node %s is already marked as left", n.Address())
		} else {
			logrus.Debugf("Node %s is already marked as left", n.Address())
		}
		return
	}
	now := time.Now()
	n.LastLeft = now
	n.LastUpdated = now
	n.AccumulatedUptime += n.GetCurrentUptime()
	n.CurrentUptime = 0
	n.Activity = ActivityLeft
	n.IsActive = false
	if n.IsStaked {
		logrus.Info("Node left: ", n.Address())
	} else {
		logrus.Debug("Node left: ", n.Address())
	}
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

func GetSelfNodeDataJson(host host.Host, isStaked bool) []byte {
	// Create and populate NodeData
	nodeData := NodeData{
		PeerId:     host.ID(),
		IsStaked:   isStaked,
		EthAddress: masacrypto.KeyManagerInstance().EthAddress,
	}

	// Convert NodeData to JSON
	jsonData, err := json.Marshal(nodeData)
	if err != nil {
		logrus.Error("Error marshalling NodeData:", err)
		return nil
	}
	return jsonData
}
