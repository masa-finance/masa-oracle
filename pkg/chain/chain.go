package chain

import (
	"context"
	"time"

	pubsub2 "github.com/libp2p/go-libp2p-pubsub"
	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/sirupsen/logrus"
)

var (
	blocksCh = make(chan *pubsub2.Message)
)

// func SubscribeToBlocks(node *masa.OracleNode) {
// 	node.BlockTracker = &pubsub.BlockEventTracker{BlocksCh: blocksCh}
// 	err := node.PubSubManager.AddSubscription(config.TopicWithVersion(config.BlockTopic), node.BlockTracker, true)
// 	if err != nil {
// 		logrus.Errorf("Subscribe error %v", err)
// 	}
// }

func SubscribeToBlocks(ctx context.Context, node *masa.OracleNode) {
	node.BlockTracker = &pubsub.BlockEventTracker{BlocksCh: blocksCh}
	err := node.PubSubManager.AddSubscription(config.TopicWithVersion(config.BlockTopic), node.BlockTracker, true)
	if err != nil {
		logrus.Errorf("Subscribe error %v", err)
	}

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logrus.Debug("tick")
		case block := <-node.BlockTracker.BlocksCh:
			logrus.Info("[+] Sending work to network", block)
		case <-ctx.Done():
			return
		}
	}
}
