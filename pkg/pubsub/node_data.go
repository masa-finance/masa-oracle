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
	MultiaddrsString     string          `json:"multiaddrsString,omitempty"`
	PeerId               peer.ID         `json:"peerId"`
	FirstJoinedUnix      int64           `json:"firstJoined,omitempty"`
	LastJoinedUnix       int64           `json:"lastJoined,omitempty"`
	LastLeftUnix         int64           `json:"-"`
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
	IsTelegramScraper    bool            `json:"isTelegramScraper"`
	IsWebScraper         bool            `json:"isWebScraper"`
	Records              any             `json:"records,omitempty"`
	Version              string          `json:"version"`
	WorkerTimeout        time.Time       `json:"workerTimeout,omitempty"`
	ReturnedTweets       int             `json:"returnedTweets"` // a running count of the number of tweets returned
	LastReturnedTweet    time.Time       `json:"lastReturnedTweet"`
	TweetTimeout         bool            `json:"tweetTimeout"`
	TweetTimeouts        int             `json:"tweetTimeouts"` // a running countthe number of times a tweet request times out
	LastTweetTimeout     time.Time       `json:"lastTweetTimeout"`
	LastNotFoundTime     time.Time       `json:"lastNotFoundTime"`
	NotFoundCount        int             `json:"notFoundCount"` // a running count of the number of times a node is not found
}

// NewNodeData creates a new NodeData struct initialized with the given
// parameters. It is used to represent data about a node in the network.
func NewNodeData(addr multiaddr.Multiaddr, peerId peer.ID, publicKey string, activity int) *NodeData {
	multiaddrs := make([]JSONMultiaddr, 0)
	multiaddrs = append(multiaddrs, JSONMultiaddr{addr})

	return &NodeData{
		PeerId:            peerId,
		Multiaddrs:        multiaddrs,
		LastUpdatedUnix:   time.Now().Unix(),
		CurrentUptime:     0,
		AccumulatedUptime: 0,
		EthAddress:        publicKey,
		Activity:          activity,
		SelfIdentified:    false,
	}
}

// Address returns a string representation of the NodeData's multiaddress
// and peer ID in the format "/ip4/127.0.0.1/tcp/4001/p2p/QmcgpsyWgH8Y8ajJz1Cu72KnS5uo2Aa2LpzU7kinSupNKC".
// This can be used by other nodes to connect to this node.
func (n *NodeData) Address() string {
	// Add a check for empty addresses
	if len(n.Multiaddrs) == 0 {
		return ""
	}
	return fmt.Sprintf("%s/p2p/%s", n.Multiaddrs[0].String(), n.PeerId.String())
}

// WorkerCategory represents the main categories of workers
type WorkerCategory int

const (
	CategoryDiscord WorkerCategory = iota
	CategoryTelegram
	CategoryTwitter
	CategoryWeb
)

// String returns the string representation of the WorkerCategory
func (wc WorkerCategory) String() string {
	return [...]string{"Discord", "Telegram", "Twitter", "Web"}[wc]
}

// CanDoWork checks if the node can perform work of the specified WorkerType.
// It returns true if the node is configured for the given worker type, false otherwise.
func (n *NodeData) CanDoWork(workerType WorkerCategory) bool {
	if !n.IsStaked {
		return false
	}
	switch workerType {
	case CategoryTwitter:
		return n.IsTwitterScraper
	case CategoryDiscord:
		return n.IsDiscordScraper
	case CategoryTelegram:
		return n.IsTelegramScraper
	case CategoryWeb:
		return n.IsWebScraper
	default:
		return false
	}
}

// TwitterScraper checks if the current node is configured as a Twitter scraper.
// It retrieves the configuration instance and returns the value of the TwitterScraper field.
func (n *NodeData) TwitterScraper() bool {
	return n.IsTwitterScraper
}

// DiscordScraper checks if the current node is configured as a Discord scraper.
// It retrieves the configuration instance and returns the value of the DiscordScraper field.
func (n *NodeData) DiscordScraper() bool {
	return n.IsDiscordScraper
}

// TelegramScraper checks if the current node is configured as a Discord scraper.
// It retrieves the configuration instance and returns the value of the TelegramScraper field.
func (n *NodeData) TelegramScraper() bool {
	return n.IsTelegramScraper
}

// WebScraper checks if the current node is configured as a Web scraper.
// It retrieves the configuration instance and returns the value of the WebScraper field.
func (n *NodeData) WebScraper() bool {
	return n.IsWebScraper
}

// Joined updates the NodeData when the node joins the network.
// It sets the join times, activity, active status, and logs based on stake status.
func (n *NodeData) Joined(nodeVersion string) {
	now := time.Now()
	if n.FirstJoinedUnix == 0 {
		n.FirstJoinedUnix = now.Unix()
	}
	n.LastJoinedUnix = now.Unix()
	n.LastUpdatedUnix = now.Unix()
	n.Activity = ActivityJoined
	n.IsActive = true

	n.Version = nodeVersion

	logMessage := fmt.Sprintf("[+] %s node joined: %s", map[bool]string{true: "Staked", false: "Unstaked"}[n.IsStaked], n.MultiaddrsString)
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
	n.LastLeftUnix = time.Now().Unix()
	n.LastUpdatedUnix = n.LastLeftUnix
	n.CurrentUptime = 0
	n.Activity = ActivityLeft
	n.IsActive = false

	n.UpdateAccumulatedUptime()
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
		return time.Since(time.Unix(n.LastJoinedUnix, 0))
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
		// Calculate the uptime for the most recent active period
		recentUptime := n.LastLeftUnix - n.LastJoinedUnix
		// Add this to the accumulated uptime
		n.AccumulatedUptime += time.Duration(recentUptime) * time.Second
	} else if n.Activity == ActivityJoined {
		// If the node is currently active, calculate the uptime since it first joined
		// This should only be done if the node is active and hasn't been updated yet
		currentUptime := time.Now().Unix() - n.FirstJoinedUnix
		// Update the accumulated uptime only if it's less than the current uptime
		if currentUptime > int64(n.AccumulatedUptime.Seconds()) {
			n.AccumulatedUptime = time.Duration(currentUptime) * time.Second
		}
	}
	// Ensure the accumulated uptime does not exceed the maximum possible uptime
	if n.FirstJoinedUnix > 0 && n.LastLeftUnix > 0 {
		maxAccumulatedUptime := time.Duration(n.LastLeftUnix-n.FirstJoinedUnix) * time.Second
		if n.AccumulatedUptime > maxAccumulatedUptime {
			n.AccumulatedUptime = maxAccumulatedUptime
		}
	}
	n.AccumulatedUptimeStr = n.AccumulatedUptime.String()
}

func (n *NodeData) MergeMultiaddresses(addr multiaddr.Multiaddr) {
	addrExists := false
	for _, existingAddr := range n.Multiaddrs {
		if addr.Equal(existingAddr.Multiaddr) {
			addrExists = true
			break
		}
	}
	if !addrExists {
		n.Multiaddrs = append(n.Multiaddrs, JSONMultiaddr{Multiaddr: addr})
	}
}

func (nd *NodeData) UpdateTwitterFields(fields NodeData) {
	if fields.ReturnedTweets != 0 {
		nd.ReturnedTweets += fields.ReturnedTweets
	}
	if !fields.LastReturnedTweet.IsZero() {
		nd.LastReturnedTweet = fields.LastReturnedTweet
	}
	if fields.TweetTimeout {
		nd.TweetTimeout = fields.TweetTimeout
		nd.TweetTimeouts += fields.TweetTimeouts
		nd.LastTweetTimeout = fields.LastTweetTimeout
	}
	if !fields.LastNotFoundTime.IsZero() {
		nd.LastNotFoundTime = fields.LastNotFoundTime
		nd.NotFoundCount += fields.NotFoundCount
	}
	nd.LastUpdatedUnix = time.Now().Unix()
}
