package pubsub

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/masa-finance/masa-oracle/pkg/config"

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

// UnmarshalJSON implements json.Unmarshaler. It parses the JSON
// representation of a multiaddress and stores the result in m.
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
	FirstJoined          time.Time       `json:"firstJoined,omitempty"`
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
	IsValidator          bool            `json:"isValidator"`
	IsTwitterScraper     bool            `json:"isTwitterScraper"`
	IsDiscordScraper     bool            `json:"isDiscordScraper"`
	IsWebScraper         bool            `json:"isWebScraper"`
	BytesScraped         int             `json:"bytesScraped"`
	Records              any             `json:"records,omitempty"`
}

// NewNodeData creates a new NodeData struct initialized with the given
// parameters. It is used to represent data about a node in the network.
func NewNodeData(addr multiaddr.Multiaddr, peerId peer.ID, publicKey string, activity int) *NodeData {
	multiaddrs := make([]JSONMultiaddr, 0)
	multiaddrs = append(multiaddrs, JSONMultiaddr{addr})
	cfg := config.GetInstance()
	wn, _ := strconv.ParseBool(cfg.WriterNode)
	ts := cfg.TwitterScraper
	ds := cfg.DiscordScraper
	ws := cfg.WebScraper

	return &NodeData{
		PeerId:            peerId,
		Multiaddrs:        multiaddrs,
		LastUpdated:       time.Now(),
		CurrentUptime:     0,
		AccumulatedUptime: 0,
		EthAddress:        publicKey,
		Activity:          activity,
		SelfIdentified:    false,
		IsValidator:       wn,
		IsTwitterScraper:  ts,
		IsDiscordScraper:  ds,
		IsWebScraper:      ws,
		BytesScraped:      0,
	}
}

// Address returns a string representation of the NodeData's multiaddress
// and peer ID in the format "/ip4/127.0.0.1/tcp/4001/p2p/QmcgpsyWgH8Y8ajJz1Cu72KnS5uo2Aa2LpzU7kinSupNKC".
// This can be used by other nodes to connect to this node.
func (n *NodeData) Address() string {
	return fmt.Sprintf("%s/p2p/%s", n.Multiaddrs[0].String(), n.PeerId.String())
}

func (n *NodeData) TwitterScraper() bool {
	cfg := config.GetInstance()
	return cfg.TwitterScraper
}

func (n *NodeData) DiscordScraper() bool {
	cfg := config.GetInstance()
	return cfg.DiscordScraper
}

func (n *NodeData) WebScraper() bool {
	cfg := config.GetInstance()
	return cfg.WebScraper
}

// Joined updates the NodeData when the node joins the network.
// It sets the join times, activity, active status, and logs based on stake status.
func (n *NodeData) Joined() {
	now := time.Now()
	n.FirstJoined = now.Add(-n.AccumulatedUptime)
	n.LastJoined = now
	n.LastUpdated = now
	n.Activity = ActivityJoined
	n.IsActive = true
	if n.IsStaked {
		logrus.Info("Staked node joined: ", n.Address())
	} else {
		logrus.Debug("Unstaked node joined: ", n.Address())
	}
}

// Left updates the NodeData when the node leaves the network.
// It sets the leave time, stops uptime, sets activity to left,
// sets node as inactive, and logs based on stake status.
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

// GetCurrentUptime returns the node's current uptime duration.
// If the node is active, it calculates the time elapsed since the last joined time.
// If the node is marked as left, it returns 0.
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

// GetAccumulatedUptime returns the total accumulated uptime for the node.
// It calculates this by adding the current uptime to the stored accumulated uptime.
func (n *NodeData) GetAccumulatedUptime() time.Duration {
	return n.AccumulatedUptime + n.GetCurrentUptime()
}

// UpdateAccumulatedUptime updates the accumulated uptime of the node to account for any
// discrepancy between the last left and last joined times from gossipsub events.
// It calculates the duration between last left and joined if the node is marked as left.
// Otherwise, it uses the time since the last joined event.
func (n *NodeData) UpdateAccumulatedUptime() {
	if n.Activity == ActivityLeft {
		n.AccumulatedUptime += n.LastLeft.Sub(n.LastJoined)
	} else {
		n.AccumulatedUptime += time.Since(n.LastJoined)
	}
}

// GetSelfNodeDataJson converts the local node's data into a JSON byte array.
// It populates a NodeData struct with the node's ID, staking status, and Ethereum address.
// The NodeData struct is then marshalled into a JSON byte array.
// Returns nil if there is an error marshalling to JSON.
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
