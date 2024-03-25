package main

import (
	"sync"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/rivo/tview"
)

type AppConfig struct {
	Address         string
	Model           string
	TwitterUser     string
	TwitterPassword string
}

var appConfig = AppConfig{}

var mainFlex *tview.Flex

type Gossip struct {
	Content  string
	Metadata map[string]string
}

type SpeakRequest struct {
	Text          string `json:"text"`
	VoiceSettings struct {
		Stability       float64 `json:"stability"`
		SimilarityBoost float64 `json:"similarity_boost"`
	} `json:"voice_settings"`
}

type SubscriptionHandler struct {
	Gossips     []Gossip
	GossipTopic *pubsub.Topic
	mu          sync.Mutex
}
