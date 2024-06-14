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
	FirstJoined          time.Time       `json:"-"`
	FirstJoinedUnix      int64           `json:"firstJoined,omitempty"`
	LastJoined           time.Time       `json:"-"`
	LastJoinedUnix       int64           `json:"lastJoined,omitempty"`
	LastLeft             time.Time       `json:"-"`
	LastLeftUnix         int64           `json:"-"`
	LastUpdated          time.Time       `json:"-"`
	LastUpdatedUnix      int64           `json:"lastUpdated,omitempty"`
	CurrentUptime        time.Duration   `json:"uptime,omitempty"`
	CurrentUptimeStr     string          `json:"uptimeStr,omitempty"`
	AccumulatedUptime    time.Duration   `json:"accumulatedUptime,omitempty"`
	AccumulatedUptimeStr string          `json:"accumulatedUptimeStr,omitempty"`
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
	Version              string          `json:"version"`
}

// NewNodeData creates a new NodeData struct initialized with the given
// parameters. It is used to represent data about a node in the network.
func NewNodeData(addr multiaddr.Multiaddr, peerId peer.ID, publicKey string, activity int) *NodeData {
	multiaddrs := make([]JSONMultiaddr, 0)
	multiaddrs = append(multiaddrs, JSONMultiaddr{addr})
	cfg := config.GetInstance()
	wn, _ := strconv.ParseBool(cfg.Validator)
	ts := cfg.TwitterScraper
	ds := cfg.DiscordScraper
	ws := cfg.WebScraper
	ver := cfg.Version

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
		Version:           ver,
	}
}

// UpdateUnixTimestamps updates the Unix timestamp fields based on the time fields.
func (n *NodeData) UpdateUnixTimestamps() {
	n.FirstJoinedUnix = n.FirstJoined.Unix()
	n.LastJoinedUnix = n.LastJoined.Unix()
	n.LastLeftUnix = n.LastLeft.Unix()
	n.LastUpdatedUnix = n.LastUpdated.Unix()
}

// CalculateCurrentUptime calculates the current uptime based on Unix timestamps.
func (n *NodeData) CalculateCurrentUptime() {
	if n.Activity == ActivityJoined {
		n.CurrentUptime = time.Duration(n.LastUpdatedUnix-n.LastJoinedUnix) * time.Second
	} else {
		n.CurrentUptime = 0
	}
	n.CurrentUptimeStr = n.CurrentUptime.String()
}

// CalculateAccumulatedUptime calculates the accumulated uptime based on Unix timestamps.
func (n *NodeData) CalculateAccumulatedUptime() {
	if n.FirstJoinedUnix > 0 && n.LastLeftUnix > 0 {
		n.AccumulatedUptime = time.Duration(n.LastLeftUnix-n.FirstJoinedUnix) * time.Second
	} else {
		n.AccumulatedUptime = 0
	}
	n.AccumulatedUptimeStr = n.AccumulatedUptime.String()
}

// Address returns a string representation of the NodeData's multiaddress
// and peer ID in the format "/ip4/127.0.0.1/tcp/4001/p2p/QmcgpsyWgH8Y8ajJz1Cu72KnS5uo2Aa2LpzU7kinSupNKC".
// This can be used by other nodes to connect to this node.
func (n *NodeData) Address() string {
	return fmt.Sprintf("%s/p2p/%s", n.Multiaddrs[0].String(), n.PeerId.String())
}

// TwitterScraper checks if the current node is configured as a Twitter scraper.
// It retrieves the configuration instance and returns the value of the TwitterScraper field.
func (n *NodeData) TwitterScraper() bool {
	cfg := config.GetInstance()
	return cfg.TwitterScraper
}

// DiscordScraper checks if the current node is configured as a Discord scraper.
// It retrieves the configuration instance and returns the value of the DiscordScraper field.
func (n *NodeData) DiscordScraper() bool {
	cfg := config.GetInstance()
	return cfg.DiscordScraper
}

// WebScraper checks if the current node is configured as a Web scraper.
// It retrieves the configuration instance and returns the value of the WebScraper field.
func (n *NodeData) WebScraper() bool {
	cfg := config.GetInstance()
	return cfg.WebScraper
}

// Joined updates the NodeData when the node joins the network.
// It sets the join times, activity, active status, and logs based on stake status.
func (n *NodeData) Joined() {
	now := time.Now()
	if n.FirstJoined.IsZero() {
		n.FirstJoined = now
	}
	n.LastJoined = now
	n.LastUpdated = now
	n.Activity = ActivityJoined
	n.IsActive = true

	if n.FirstJoinedUnix == 0 {
		n.FirstJoinedUnix = time.Now().Unix()
	}
	n.UpdateUnixTimestamps()
	n.CalculateCurrentUptime()
	n.CalculateAccumulatedUptime()

	n.Version = "0.0.7-beta"

	logMessage := fmt.Sprintf("[+] %s node joined: %s", map[bool]string{true: "Staked", false: "Unstaked"}[n.IsStaked], n.Address())
	if n.IsStaked {
		logrus.Info(logMessage)
	} else {
		logrus.Debug(logMessage)
	}
}

// Left updates the NodeData when the node leaves the network.
// It sets the leave time, stops uptime, sets activity to left,
// sets node as inactive, and logs based on stake status.
func (n *NodeData) Left() {
	if n.Activity == ActivityLeft {
		return
	}
	n.LastLeft = time.Now()
	n.LastUpdated = n.LastLeft
	n.AccumulatedUptime += n.GetCurrentUptime()
	n.CurrentUptime = 0
	n.Activity = ActivityLeft
	n.IsActive = false
	n.UpdateUnixTimestamps()
	n.CalculateCurrentUptime()
	n.CalculateAccumulatedUptime()

	logMessage := fmt.Sprintf("Node left: %s", n.Address())
	if n.IsStaked {
		logrus.Info(logMessage)
	} else {
		logrus.Debug(logMessage)
	}
}

// GetCurrentUptime returns the node's current uptime duration.
// If the node is active, it calculates the time elapsed since the last joined time.
// If the node is marked as left, it returns 0.
func (n *NodeData) GetCurrentUptime() time.Duration {
	if n.Activity == ActivityJoined {
		return time.Since(n.LastJoined)
	}
	return 0
}

// GetAccumulatedUptime returns the total accumulated uptime for the node.
// It calculates this by adding the current uptime to the stored accumulated uptime.
func (n *NodeData) GetAccumulatedUptime() time.Duration {
	currentUptime := n.GetCurrentUptime()
	return n.AccumulatedUptime + currentUptime
}

// UpdateAccumulatedUptime updates the accumulated uptime of the node to account for any
// discrepancy between the last left and last joined times from gossipsub events.
// It calculates the duration between last left and joined if the node is marked as left.
// Otherwise, it uses the time since the last joined event.
func (n *NodeData) UpdateAccumulatedUptime() {
	if n.Activity == ActivityLeft {
		n.AccumulatedUptime += n.LastLeft.Sub(n.LastJoined)
		return
	}
	n.AccumulatedUptime += time.Since(n.LastJoined)
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
