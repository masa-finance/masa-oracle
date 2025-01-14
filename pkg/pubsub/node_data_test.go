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

	t.Run("NewNodeData", func(t *testing.T) {
		nodeData := NewNodeData(
			[]multiaddr.Multiaddr{testAddr},
			testPeerID,
			"0x123",
			ActivityJoined,
		)

		assert.Equal(t, testPeerID, nodeData.PeerId)
		assert.Equal(t, 1, len(nodeData.Multiaddrs))
		assert.Equal(t, testAddr.String(), nodeData.Multiaddrs[0].String())
		assert.Equal(t, "0x123", nodeData.EthAddress)
		assert.Equal(t, ActivityJoined, nodeData.Activity)
		assert.False(t, nodeData.SelfIdentified)
	})

	t.Run("Address", func(t *testing.T) {
		nodeData := NewNodeData(
			[]multiaddr.Multiaddr{testAddr},
			testPeerID,
			"0x123",
			ActivityJoined,
		)

		expected := "/ip4/127.0.0.1/tcp/4001/p2p/QmcgpsyWgH8Y8ajJz1Cu72KnS5uo2Aa2LpzU7kinSupNKC"
		assert.Equal(t, expected, nodeData.Address())

		// Test empty addresses
		emptyNode := NewNodeData(
			[]multiaddr.Multiaddr{},
			testPeerID,
			"0x123",
			ActivityJoined,
		)
		assert.Equal(t, "", emptyNode.Address())
	})

	t.Run("Join and Leave", func(t *testing.T) {
		nodeData := NewNodeData(
			[]multiaddr.Multiaddr{testAddr},
			testPeerID,
			"0x123",
			ActivityJoined,
		)

		nodeData.IsStaked = true
		nodeVersion := "1.0.0"

		// Call Joined() first to set initial state
		nodeData.Joined(nodeVersion)

		// Then manually adjust the join time to be in the past
		pastTime := time.Now().Add(-1 * time.Hour).Unix()
		nodeData.FirstJoinedUnix = pastTime
		nodeData.LastJoinedUnix = pastTime

		// Verify initial state
		assert.Equal(t, ActivityJoined, nodeData.Activity)
		assert.True(t, nodeData.IsActive)
		assert.Equal(t, nodeVersion, nodeData.Version)

		// Wait a tiny bit to ensure time difference
		time.Sleep(time.Millisecond)

		// Test Leave
		nodeData.Left()

		assert.Equal(t, ActivityLeft, nodeData.Activity)
		assert.False(t, nodeData.IsActive)
		assert.NotZero(t, nodeData.LastLeftUnix)
		assert.Greater(t, nodeData.LastLeftUnix, nodeData.LastJoinedUnix,
			"LastLeftUnix (%d) should be greater than LastJoinedUnix (%d)",
			nodeData.LastLeftUnix, nodeData.LastJoinedUnix)

		// Calculate expected minimum uptime
		expectedMinUptime := time.Duration(nodeData.LastLeftUnix-nodeData.LastJoinedUnix) * time.Second

		// Verify uptime was updated
		assert.GreaterOrEqual(t, nodeData.AccumulatedUptime, expectedMinUptime,
			"Expected accumulated uptime to be at least %v, got %v",
			expectedMinUptime, nodeData.AccumulatedUptime)
		assert.NotEmpty(t, nodeData.AccumulatedUptimeStr)
	})

	t.Run("Uptime Calculations", func(t *testing.T) {
		nodeData := NewNodeData(
			[]multiaddr.Multiaddr{testAddr},
			testPeerID,
			"0x123",
			ActivityJoined,
		)

		// Set a fixed join time in the past
		pastTime := time.Now().Add(-1 * time.Hour)
		nodeData.LastJoinedUnix = pastTime.Unix()
		nodeData.FirstJoinedUnix = pastTime.Unix()
		nodeData.Activity = ActivityJoined
		nodeData.IsActive = true

		// Check current uptime while active
		currentUptime := nodeData.GetCurrentUptime()
		assert.True(t, currentUptime >= time.Hour-time.Second, "Expected uptime to be approximately 1 hour")

		// Leave and verify uptime calculations
		nodeData.Left()
		assert.Equal(t, time.Duration(0), nodeData.GetCurrentUptime(), "Uptime should be 0 after leaving")
		assert.True(t, nodeData.AccumulatedUptime >= time.Hour-time.Second, "Expected accumulated uptime to be approximately 1 hour")

		// Check final accumulated uptime
		accumulatedUptime := nodeData.GetAccumulatedUptime()
		assert.True(t, accumulatedUptime >= time.Hour-time.Second, "Expected final accumulated uptime to be approximately 1 hour")
	})

	t.Run("Worker Categories", func(t *testing.T) {
		nodeData := NewNodeData(
			[]multiaddr.Multiaddr{testAddr},
			testPeerID,
			"0x123",
			ActivityJoined,
		)

		// Test unstaked node
		assert.False(t, nodeData.CanDoWork(CategoryTwitter))
		assert.False(t, nodeData.CanDoWork(CategoryWeb))

		// Test staked node
		nodeData.IsStaked = true
		nodeData.IsTwitterScraper = true
		assert.True(t, nodeData.CanDoWork(CategoryTwitter))
		assert.False(t, nodeData.CanDoWork(CategoryWeb))

		nodeData.IsWebScraper = true
		assert.True(t, nodeData.CanDoWork(CategoryWeb))
	})

	t.Run("TwitterFields", func(t *testing.T) {
		nodeData := NewNodeData(
			[]multiaddr.Multiaddr{testAddr},
			testPeerID,
			"0x123",
			ActivityJoined,
		)

		updateFields := NodeData{
			ReturnedTweets:    5,
			LastReturnedTweet: time.Now(),
			TweetTimeout:      true,
			TweetTimeouts:     1,
			LastTweetTimeout:  time.Now(),
			LastNotFoundTime:  time.Now(),
			NotFoundCount:     1,
		}

		nodeData.UpdateTwitterFields(updateFields)

		assert.Equal(t, 5, nodeData.ReturnedTweets)
		assert.True(t, nodeData.TweetTimeout)
		assert.Equal(t, 1, nodeData.TweetTimeouts)
		assert.Equal(t, 1, nodeData.NotFoundCount)
		assert.False(t, nodeData.LastReturnedTweet.IsZero())
		assert.False(t, nodeData.LastTweetTimeout.IsZero())
		assert.False(t, nodeData.LastNotFoundTime.IsZero())
	})

	t.Run("MergeMultiaddresses", func(t *testing.T) {
		nodeData := NewNodeData(
			[]multiaddr.Multiaddr{testAddr},
			testPeerID,
			"0x123",
			ActivityJoined,
		)

		// Add same address - should not duplicate
		nodeData.MergeMultiaddresses(testAddr)
		assert.Equal(t, 1, len(nodeData.Multiaddrs))

		// Add new address
		newAddr, err := multiaddr.NewMultiaddr("/ip4/192.168.1.1/tcp/4001")
		assert.NoError(t, err)
		nodeData.MergeMultiaddresses(newAddr)
		assert.Equal(t, 2, len(nodeData.Multiaddrs))
		assert.Equal(t, newAddr.String(), nodeData.Multiaddrs[1].String())
	})
}
