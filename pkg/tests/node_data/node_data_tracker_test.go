package node_data

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/masa-finance/masa-oracle/node"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
)

var _ = Describe("NodeDataTracker", func() {
	Context("UpdateNodeDataTwitter", func() {
		It("should correctly update NodeData Twitter fields", func() {
			testNode, err := NewOracleNode(
				context.Background(),
				EnableStaked,
				EnableRandomIdentity,
			)
			Expect(err).NotTo(HaveOccurred())

			err = testNode.Start()
			Expect(err).NotTo(HaveOccurred())

			initialData := pubsub.NodeData{
				PeerId:            testNode.Host.ID(),
				LastReturnedTweet: time.Now().Add(-1 * time.Hour),
				ReturnedTweets:    10,
				TweetTimeout:      true,
				TweetTimeouts:     2,
				LastTweetTimeout:  time.Now().Add(-1 * time.Hour),
				LastNotFoundTime:  time.Now().Add(-1 * time.Hour),
				NotFoundCount:     1,
			}

			err = testNode.NodeTracker.UpdateNodeDataTwitter(testNode.Host.ID().String(), initialData)
			Expect(err).NotTo(HaveOccurred())

			updates := pubsub.NodeData{
				ReturnedTweets:    5,
				LastReturnedTweet: time.Now(),
				TweetTimeout:      true,
				TweetTimeouts:     1,
				LastTweetTimeout:  time.Now(),
				LastNotFoundTime:  time.Now(),
				NotFoundCount:     1,
			}

			err = testNode.NodeTracker.UpdateNodeDataTwitter(testNode.Host.ID().String(), updates)
			Expect(err).NotTo(HaveOccurred())

			updatedData := testNode.NodeTracker.GetNodeData(testNode.Host.ID().String())
			Expect(updatedData.ReturnedTweets).To(Equal(15))
			Expect(updatedData.TweetTimeouts).To(Equal(3))
			Expect(updatedData.NotFoundCount).To(Equal(2))
		})
	})
})
