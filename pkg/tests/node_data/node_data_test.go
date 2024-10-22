package node_data

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/masa-oracle/pkg/pubsub"
)

var _ = Describe("NodeData", func() {
	Describe("UpdateTwitterFields", func() {
		It("should correctly update Twitter fields", func() {
			initialData := pubsub.NodeData{
				ReturnedTweets: 10,
				TweetTimeouts:  2,
				NotFoundCount:  1,
			}

			updates := pubsub.NodeData{
				ReturnedTweets:    5,
				LastReturnedTweet: time.Now(),
				TweetTimeout:      true,
				TweetTimeouts:     1,
				LastTweetTimeout:  time.Now(),
				LastNotFoundTime:  time.Now(),
				NotFoundCount:     1,
			}

			initialData.UpdateTwitterFields(updates)

			Expect(initialData.ReturnedTweets).To(Equal(15))
			Expect(initialData.TweetTimeouts).To(Equal(3))
			Expect(initialData.NotFoundCount).To(Equal(2))
			Expect(initialData.LastReturnedTweet.IsZero()).To(BeFalse())
			Expect(initialData.LastTweetTimeout.IsZero()).To(BeFalse())
			Expect(initialData.LastNotFoundTime.IsZero()).To(BeFalse())
		})
	})
})
