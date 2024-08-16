package data_types

import (
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/pubsub"
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

// WorkerTypeToCategory maps WorkerType to WorkerCategory
func WorkerTypeToCategory(wt WorkerType) pubsub.WorkerCategory {
	logrus.Infof("Mapping WorkerType %s to WorkerCategory", wt)
	switch wt {
	case Discord, DiscordProfile, DiscordChannelMessages, DiscordSentiment, DiscordGuildChannels, DiscordUserGuilds:
		logrus.Info("WorkerType is related to Discord")
		return pubsub.CategoryDiscord
	case TelegramSentiment, TelegramChannelMessages:
		logrus.Info("WorkerType is related to Telegram")
		return pubsub.CategoryTelegram
	case Twitter, TwitterFollowers, TwitterProfile, TwitterSentiment, TwitterTrends:
		logrus.Info("WorkerType is related to Twitter")
		return pubsub.CategoryTwitter
	case Web, WebSentiment:
		logrus.Info("WorkerType is related to Web")
		return pubsub.CategoryWeb
	default:
		logrus.Warn("WorkerType is invalid or not recognized")
		return -1 // Invalid category
	}
}
