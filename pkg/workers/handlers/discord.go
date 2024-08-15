package handlers

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/scrapers/discord"
	"github.com/masa-finance/masa-oracle/pkg/workers/types"
)

type DiscordProfileHandler struct{}
type DiscordChannelHandler struct{}
type DiscordSentimentHandler struct{}
type DiscordGuildHandler struct{}
type DiscoreUserGuildsHandler struct{}

// HandleWork implements the WorkHandler interface for DiscordProfileHandler.
func (h *DiscordProfileHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] DiscordProfileHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Errorf("unable to parse discord json data: %v", err)}
	}
	userID := dataMap["userID"].(string)
	resp, err := discord.GetUserProfile(userID)
	return data_types.WorkResponse{Data: resp, Error: err}
}

// HandleWork implements the WorkHandler interface for DiscordChannelHandler.
func (h *DiscordChannelHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] DiscordChannelHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Errorf("unable to parse discord json data: %v", err)}
	}
	channelID := dataMap["channelID"].(string)
	limit := dataMap["limit"].(string)
	prompt := dataMap["prompt"].(string)
	resp, err := discord.GetChannelMessages(channelID, limit, prompt)
	return data_types.WorkResponse{Data: resp, Error: err}
}

// HandleWork implements the WorkHandler interface for DiscordSentimentHandler.
func (h *DiscordSentimentHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] DiscordSentimentHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Errorf("unable to parse discord json data: %v", err)}
	}
	channelID := dataMap["channelID"].(string)
	model := dataMap["model"].(string)
	prompt := dataMap["prompt"].(string)
	_, resp, err := discord.ScrapeDiscordMessagesForSentiment(channelID, model, prompt)

	return data_types.WorkResponse{Data: resp, Error: err}
}

// HandleWork implements the WorkHandler interface for DiscordGuildHandler.
func (h *DiscordGuildHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] DiscordGuildHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Errorf("unable to parse discord json data: %v", err)}
	}
	guildID := dataMap["guildID"].(string)
	resp, err := discord.GetGuildChannels(guildID)
	return data_types.WorkResponse{Data: resp, Error: err}
}

func (h *DiscoreUserGuildsHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] DiscordUserGuildsHandler %s", data)
	resp, err := discord.GetUserGuilds()
	return data_types.WorkResponse{Data: resp, Error: err}
}
