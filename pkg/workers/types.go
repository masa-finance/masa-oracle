package workers

import (
	"github.com/asynkron/protoactor-go/actor"
	pubsub2 "github.com/libp2p/go-libp2p-pubsub"
	masa "github.com/masa-finance/masa-oracle/pkg"
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

var (
	clients        = actor.NewPIDSet()
	workerStatusCh = make(chan *pubsub2.Message)
	workerDoneCh   = make(chan *pubsub2.Message)
)

type ChanResponse struct {
	Response  map[string]interface{}
	ChannelId string
}

type Worker struct {
	Node *masa.OracleNode
}
