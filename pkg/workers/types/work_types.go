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
	TelegramChannelMessages WorkerType = "telegram-channel-messages"
	DiscordGuildChannels    WorkerType = "discord-guild-channels"
	DiscordUserGuilds       WorkerType = "discord-user-guilds"
	Twitter                 WorkerType = "twitter"
	TwitterFollowers        WorkerType = "twitter-followers"
	TwitterProfile          WorkerType = "twitter-profile"
	TwitterTweetByID        WorkerType = "twitter-tweet-by-id"

	Web  WorkerType = "web"
	Test WorkerType = "test"

	DataSourceTwitter  = "twitter"
	DataSourceDiscord  = "discord"
	DataSourceWeb      = "web"
	DataSourceTelegram = "telegram"
)

// WorkerTypeToCategory maps WorkerType to WorkerCategory
func WorkerTypeToCategory(wt WorkerType) pubsub.WorkerCategory {
	logrus.Infof("Mapping WorkerType %s to WorkerCategory", wt)
	switch wt {
	case Discord, DiscordProfile, DiscordChannelMessages, DiscordGuildChannels, DiscordUserGuilds:
		logrus.Info("WorkerType is related to Discord")
		return pubsub.CategoryDiscord
	case TelegramChannelMessages:
		logrus.Info("WorkerType is related to Telegram")
		return pubsub.CategoryTelegram
	case Twitter, TwitterFollowers, TwitterProfile, TwitterTweetByID:
		logrus.Info("WorkerType is related to Twitter")
		return pubsub.CategoryTwitter
	case Web:
		logrus.Info("WorkerType is related to Web")
		return pubsub.CategoryWeb
	default:
		logrus.Warn("WorkerType is invalid or not recognized")
		return -1 // Invalid category
	}
}

// WorkerTypeToDataSource maps WorkerType to WorkerCategory
func WorkerTypeToDataSource(wt WorkerType) string {
	logrus.Infof("Mapping WorkerType %s to WorkerCategory", wt)
	switch wt {
	case Discord, DiscordProfile, DiscordChannelMessages, DiscordGuildChannels, DiscordUserGuilds:
		logrus.Info("WorkerType is related to Discord")
		return DataSourceDiscord
	case TelegramChannelMessages:
		logrus.Info("WorkerType is related to Telegram")
		return DataSourceTelegram
	case Twitter, TwitterFollowers, TwitterProfile, TwitterTweetByID:
		logrus.Info("WorkerType is related to Twitter")
		return DataSourceTwitter
	case Web:
		logrus.Info("WorkerType is related to Web")
		return DataSourceWeb
	default:
		logrus.Warn("WorkerType is invalid or not recognized")
		return "" // Invalid category
	}
}
