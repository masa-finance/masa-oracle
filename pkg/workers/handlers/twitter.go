package handlers

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/scrapers/twitter"
	data_types "github.com/masa-finance/masa-oracle/pkg/workers/types"
)

type TwitterQueryHandler struct{}
type TwitterFollowersHandler struct{}
type TwitterProfileHandler struct{}
type TwitterTweetHandler struct{}

func (h *TwitterTweetHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] TwitterTweetHandler input: %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		logrus.Errorf("[+] TwitterTweetHandler error parsing data: %v", err)
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to parse tweet data: %v", err)}
	}
	tweetID := dataMap["id"].(string)

	logrus.Infof("[+] Fetching tweet with ID: %s", tweetID)

	resp, loginEvent, err := twitter.ScrapeTweetByID(tweetID)
	if err != nil {
		logrus.Errorf("[+] TwitterTweetHandler error fetching tweet: %v", err)
		return data_types.WorkResponse{Error: err.Error(), LoginEvent: loginEvent}
	}

	logrus.Infof("[+] TwitterTweetHandler Work response for %s: tweet returned", data_types.TwitterTweet)
	return data_types.WorkResponse{Data: resp, RecordCount: 1}
}

func (h *TwitterQueryHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] TwitterQueryHandler input: %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		logrus.Errorf("[+] TwitterQueryHandler error parsing data: %v", err)
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to parse twitter query data: %v", err)}
	}
	count := int(dataMap["count"].(float64))
	query := dataMap["query"].(string)

	logrus.Infof("[+] Scraping tweets for query: %s, count: %d", query, count)

	resp, loginEvent, err := twitter.ScrapeTweetsByQuery(query, count)
	if err != nil {
		logrus.Errorf("[+] TwitterQueryHandler error scraping tweets: %v", err)
		return data_types.WorkResponse{Error: err.Error(), LoginEvent: loginEvent}
	}

	logrus.Infof("[+] TwitterQueryHandler Work response for %s: %d tweets returned", data_types.Twitter, len(resp))
	if len(resp) > 0 && resp[0].Tweet != nil {
		tweet := resp[0].Tweet
		logrus.Infof("[+] First tweet: ID: %s, Text: %s, Author: %s, CreatedAt: %s",
			tweet.ID, tweet.Text, tweet.Username, tweet.TimeParsed)
	}
	return data_types.WorkResponse{Data: resp, RecordCount: len(resp)}
}

func (h *TwitterFollowersHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] TwitterFollowersHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to parse twitter followers data: %v", err)}
	}
	username := dataMap["username"].(string)
	count := int(dataMap["count"].(float64))
	resp, loginEvent, err := twitter.ScrapeFollowersForProfile(username, count)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to get twitter followers: %v", err), LoginEvent: loginEvent}
	}

	logrus.Infof("[+] TwitterFollowersHandler Work response for %s: %d records returned", data_types.TwitterFollowers, len(resp))
	return data_types.WorkResponse{Data: resp, RecordCount: len(resp)}
}

func (h *TwitterProfileHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] TwitterProfileHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to parse twitter profile data: %v", err)}
	}
	username := dataMap["username"].(string)
	resp, loginEvent, err := twitter.ScrapeTweetsProfile(username)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to get twitter profile: %v", err), LoginEvent: loginEvent}
	}
	logrus.Infof("[+] TwitterProfileHandler Work response for %s: %d records returned", data_types.TwitterProfile, 1)
	return data_types.WorkResponse{Data: resp, RecordCount: 1}
}
