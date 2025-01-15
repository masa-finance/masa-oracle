package pubsub

import (
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/assert"
)

func TestNodeData(t *testing.T) {
	// Setup test data
	testAddr, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/4001")
	assert.NoError(t, err)

	testPeerID, err := peer.Decode("QmcgpsyWgH8Y8ajJz1Cu72KnS5uo2Aa2LpzU7kinSupNKC")
	assert.NoError(t, err)

	t.Run("NewNodeData with multiple addresses", func(t *testing.T) {
		addr2, err := multiaddr.NewMultiaddr("/ip4/192.168.1.1/tcp/4001")
		assert.NoError(t, err)

		nodeData := NewNodeData(
			[]multiaddr.Multiaddr{testAddr, addr2},
			testPeerID,
			"0x123",
			ActivityJoined,
		)

		assert.Equal(t, 2, len(nodeData.Multiaddrs))
		assert.Contains(t, nodeData.MultiaddrsString, testAddr.String())
		assert.Contains(t, nodeData.MultiaddrsString, addr2.String())
	})

	t.Run("UpdateAccumulatedUptime with different scenarios", func(t *testing.T) {
		nodeData := NewNodeData(
			[]multiaddr.Multiaddr{testAddr},
			testPeerID,
			"0x123",
			ActivityJoined,
		)

		// Test case 1: Node just joined
		now := time.Now()
		nodeData.FirstJoinedUnix = now.Add(-2 * time.Hour).Unix()
		nodeData.LastJoinedUnix = now.Add(-1 * time.Hour).Unix()
		nodeData.Activity = ActivityJoined
		nodeData.IsActive = true

		// Let's simulate some time passage
		time.Sleep(10 * time.Millisecond)

		currentUptime := nodeData.GetCurrentUptime()
		assert.True(t, currentUptime > 0, "Current uptime should be greater than 0")

		// Test case 2: Node left
		nodeData.LastLeftUnix = now.Unix()
		nodeData.Activity = ActivityLeft
		nodeData.IsActive = false
		nodeData.UpdateAccumulatedUptime()

		// The accumulated uptime should be around 1 hour (difference between LastLeft and LastJoined)
		expectedUptime := time.Duration(nodeData.LastLeftUnix-nodeData.LastJoinedUnix) * time.Second
		assert.True(t, nodeData.AccumulatedUptime > 0,
			"Expected accumulated uptime > 0, got %v", nodeData.AccumulatedUptime)
		assert.True(t, nodeData.AccumulatedUptime <= expectedUptime,
			"Expected accumulated uptime <= %v, got %v", expectedUptime, nodeData.AccumulatedUptime)
		assert.NotEmpty(t, nodeData.AccumulatedUptimeStr)

		// Test case 3: Multiple join/leave cycles
		firstCycleUptime := nodeData.AccumulatedUptime

		// Second cycle
		laterTime := now.Add(1 * time.Hour)
		nodeData.LastJoinedUnix = laterTime.Unix()
		nodeData.LastLeftUnix = laterTime.Add(30 * time.Minute).Unix()
		nodeData.Activity = ActivityLeft
		nodeData.UpdateAccumulatedUptime()

		// Accumulated uptime should increase
		assert.True(t, nodeData.AccumulatedUptime > firstCycleUptime,
			"Expected accumulated uptime to increase after second cycle")
	})

	t.Run("WorkerCategory String representation", func(t *testing.T) {
		assert.Equal(t, "Discord", CategoryDiscord.String())
		assert.Equal(t, "Telegram", CategoryTelegram.String())
		assert.Equal(t, "Twitter", CategoryTwitter.String())
		assert.Equal(t, "Web", CategoryWeb.String())
	})

	t.Run("UpdateTwitterFields with zero values", func(t *testing.T) {
		nodeData := NewNodeData(
			[]multiaddr.Multiaddr{testAddr},
			testPeerID,
			"0x123",
			ActivityJoined,
		)

		// Initial update with non-zero values
		initialUpdate := NodeData{
			ReturnedTweets:    5,
			LastReturnedTweet: time.Now(),
			TweetTimeout:      true,
			TweetTimeouts:     1,
		}
		nodeData.UpdateTwitterFields(initialUpdate)

		// Update with zero values
		zeroUpdate := NodeData{
			ReturnedTweets: 0,
			TweetTimeout:   false,
			TweetTimeouts:  0,
		}
		nodeData.UpdateTwitterFields(zeroUpdate)

		// Original values should be preserved
		assert.Equal(t, 5, nodeData.ReturnedTweets)
		assert.Equal(t, 1, nodeData.TweetTimeouts)
	})

	t.Run("CanDoWork with different worker types", func(t *testing.T) {
		nodeData := NewNodeData(
			[]multiaddr.Multiaddr{testAddr},
			testPeerID,
			"0x123",
			ActivityJoined,
		)

		nodeData.IsStaked = true

		// Test each worker category
		testCases := []struct {
			category WorkerCategory
			setup    func()
			expected bool
		}{
			{
				category: CategoryDiscord,
				setup:    func() {},
				expected: false,
			},
			{
				category: CategoryTelegram,
				setup:    func() {},
				expected: false,
			},
			{
				category: CategoryTwitter,
				setup: func() {
					nodeData.IsTwitterScraper = true
				},
				expected: true,
			},
			{
				category: CategoryWeb,
				setup: func() {
					nodeData.IsWebScraper = true
				},
				expected: true,
			},
		}

		for _, tc := range testCases {
			tc.setup()
			result := nodeData.CanDoWork(tc.category)
			assert.Equal(t, tc.expected, result, "Category %v test failed", tc.category)
		}
	})

	t.Run("GetCurrentUptime with various states", func(t *testing.T) {
		nodeData := NewNodeData(
			[]multiaddr.Multiaddr{testAddr},
			testPeerID,
			"0x123",
			ActivityJoined,
		)

		// Test active node
		now := time.Now()
		nodeData.LastJoinedUnix = now.Add(-1 * time.Hour).Unix()
		nodeData.Activity = ActivityJoined
		uptime := nodeData.GetCurrentUptime()
		assert.True(t, uptime >= time.Hour-time.Second)
		assert.True(t, uptime <= time.Hour+time.Second)

		// Test inactive node
		nodeData.Activity = ActivityLeft
		assert.Equal(t, time.Duration(0), nodeData.GetCurrentUptime())
	})

	t.Run("MergeMultiaddresses with duplicate addresses", func(t *testing.T) {
		nodeData := NewNodeData(
			[]multiaddr.Multiaddr{testAddr},
			testPeerID,
			"0x123",
			ActivityJoined,
		)

		// Add same address multiple times
		for i := 0; i < 3; i++ {
			nodeData.MergeMultiaddresses(testAddr)
		}
		assert.Equal(t, 1, len(nodeData.Multiaddrs))

		// Add multiple different addresses
		addrs := []string{
			"/ip4/192.168.1.1/tcp/4001",
			"/ip4/192.168.1.2/tcp/4001",
			"/ip4/192.168.1.1/tcp/4001", // Duplicate
		}

		for _, addrStr := range addrs {
			addr, _ := multiaddr.NewMultiaddr(addrStr)
			nodeData.MergeMultiaddresses(addr)
		}
		assert.Equal(t, 3, len(nodeData.Multiaddrs))
	})
}
