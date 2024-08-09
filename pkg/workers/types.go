package workers

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	pubsub2 "github.com/libp2p/go-libp2p-pubsub"
)

type WorkerType string

const (
	Discord                 WorkerType = "discord"
	DiscordProfile          WorkerType = "discord-profile"
	DiscordChannelMessages  WorkerType = "discord-channel-messages"
	DiscordSentiment        WorkerType = "discord-sentiment"
	TelegramSentiment       WorkerType = "telegram-sentiment"
	TelegramChannelMessages WorkerType = "telegram-channel-messages"
	DiscordGuildChannels    WorkerType = "discord-guild-channels"
	DiscordUserGuilds       WorkerType = "discord-user-guilds"
	LLMChat                 WorkerType = "llm-chat"
	Twitter                 WorkerType = "twitter"
	TwitterFollowers        WorkerType = "twitter-followers"
	TwitterProfile          WorkerType = "twitter-profile"
	TwitterSentiment        WorkerType = "twitter-sentiment"
	TwitterTrends           WorkerType = "twitter-trends"
	Web                     WorkerType = "web"
	WebSentiment            WorkerType = "web-sentiment"
	Test                    WorkerType = "test"
)

var WORKER = struct {
	Discord, DiscordProfile, DiscordChannelMessages, DiscordSentiment, TelegramSentiment, TelegramChannelMessages, DiscordGuildChannels, DiscordUserGuilds, LLMChat, Twitter, TwitterFollowers, TwitterProfile, TwitterSentiment, TwitterTrends, Web, WebSentiment, Test WorkerType
}{
	Discord:                 Discord,
	DiscordProfile:          DiscordProfile,
	DiscordChannelMessages:  DiscordChannelMessages,
	DiscordSentiment:        DiscordSentiment,
	TelegramSentiment:       TelegramSentiment,
	TelegramChannelMessages: TelegramChannelMessages,
	DiscordGuildChannels:    DiscordGuildChannels,
	DiscordUserGuilds:       DiscordUserGuilds,
	LLMChat:                 LLMChat,
	Twitter:                 Twitter,
	TwitterFollowers:        TwitterFollowers,
	TwitterProfile:          TwitterProfile,
	TwitterSentiment:        TwitterSentiment,
	TwitterTrends:           TwitterTrends,
	Web:                     Web,
	WebSentiment:            WebSentiment,
	Test:                    Test,
}

var workerTypeMap = map[string]WorkerType{
	"discord":                   Discord,
	"discord-profile":           DiscordProfile,
	"discord-channel-messages":  DiscordChannelMessages,
	"discord-sentiment":         DiscordSentiment,
	"telegram-sentiment":        TelegramSentiment,
	"telegram-channel-messages": TelegramChannelMessages,
	"discord-guild-channels":    DiscordGuildChannels,
	"discord-user-guilds":       DiscordUserGuilds,
	"llm-chat":                  LLMChat,
	"twitter":                   Twitter,
	"twitter-followers":         TwitterFollowers,
	"twitter-profile":           TwitterProfile,
	"twitter-sentiment":         TwitterSentiment,
	"twitter-trends":            TwitterTrends,
	"web":                       Web,
	"web-sentiment":             WebSentiment,
	"test":                      Test,
}

func StringToWorkerType(s string) (WorkerType, error) {
	if workerType, ok := workerTypeMap[s]; ok {
		return workerType, nil
	}
	return "", fmt.Errorf("invalid WorkerType: %s", s)
}

var (
	clients        = actor.NewPIDSet()
	workerStatusCh = make(chan *pubsub2.Message)
	workerDoneCh   chan *pubsub2.Message
)

func init() {
	config, _ := LoadConfig()
	workerDoneCh = make(chan *pubsub2.Message, config.WorkerBufferSize)
}

type ChanResponse struct {
	Response  map[string]interface{}
	ChannelId string
}

type roundRobinIterator struct {
	workers []Worker
	index   int
	tried   map[int]bool
}
