package twitter

import (
	"context"
	"time"

	twitterscraper "github.com/imperatrona/twitter-scraper"

	data_types "github.com/masa-finance/masa-oracle/pkg/workers/types"
)

type SimpleTweet struct {
	ConversationID    string
	GIFs              []twitterscraper.GIF
	Hashtags          []string
	HTML              string
	ID                string
	InReplyToStatusID string
	IsQuoted          bool
	IsPin             bool
	IsReply           bool
	IsRetweet         bool
	IsSelfThread      bool
	Likes             int
	Name              string
	Mentions          []twitterscraper.Mention
	PermanentURL      string
	Photos            []twitterscraper.Photo
	Place             *twitterscraper.Place
	QuotedStatusID    string
	Replies           int
	Retweets          int
	RetweetedStatusID string
	Text              string
	TimeParsed        time.Time
	Timestamp         int64
	URLs              []string
	UserID            string
	Username          string
	Videos            []twitterscraper.Video
	Views             int
	SensitiveContent  bool
}

type TweetResult struct {
	Tweet *SimpleTweet
	Error error
}

// Conversion function
func convertToSimpleTweet(tweet *twitterscraper.Tweet) *SimpleTweet {
	return &SimpleTweet{
		ConversationID:    tweet.ConversationID,
		GIFs:              tweet.GIFs,
		Hashtags:          tweet.Hashtags,
		HTML:              tweet.HTML,
		ID:                tweet.ID,
		InReplyToStatusID: tweet.InReplyToStatusID,
		IsQuoted:          tweet.IsQuoted,
		IsPin:             tweet.IsPin,
		IsReply:           tweet.IsReply,
		IsRetweet:         tweet.IsRetweet,
		IsSelfThread:      tweet.IsSelfThread,
		Likes:             tweet.Likes,
		Name:              tweet.Name,
		Mentions:          tweet.Mentions,
		PermanentURL:      tweet.PermanentURL,
		Photos:            tweet.Photos,
		Place:             tweet.Place,
		QuotedStatusID:    tweet.QuotedStatusID,
		Replies:           tweet.Replies,
		Retweets:          tweet.Retweets,
		RetweetedStatusID: tweet.RetweetedStatusID,
		Text:              tweet.Text,
		TimeParsed:        tweet.TimeParsed,
		Timestamp:         tweet.Timestamp,
		URLs:              tweet.URLs,
		UserID:            tweet.UserID,
		Username:          tweet.Username,
		Videos:            tweet.Videos,
		Views:             tweet.Views,
		SensitiveContent:  tweet.SensitiveContent,
	}
}

func ScrapeTweetByID(id string) (*SimpleTweet, *data_types.LoginEvent, error) {
	scraper, account, loginEvent, err := getAuthenticatedScraper()
	if err != nil {
		return nil, loginEvent, err
	}

	tweet, err := scraper.GetTweet(id)
	if err != nil {
		if handleRateLimit(err, account) {
			return nil, loginEvent, err
		}
		return nil, loginEvent, err
	}
	return convertToSimpleTweet(tweet), loginEvent, nil
}

func ScrapeTweetsByQuery(query string, count int) ([]*TweetResult, *data_types.LoginEvent, error) {
	scraper, account, loginEvent, err := getAuthenticatedScraper()
	if err != nil {
		return nil, loginEvent, err
	}

	var tweets []*TweetResult
	ctx := context.Background()
	scraper.SetSearchMode(twitterscraper.SearchLatest)
	for tweet := range scraper.SearchTweets(ctx, query, count) {
		if tweet.Error != nil {
			if handleRateLimit(tweet.Error, account) {
				return nil, loginEvent, tweet.Error
			}
			return nil, loginEvent, tweet.Error
		}
		tweets = append(tweets, &TweetResult{Tweet: convertToSimpleTweet(&tweet.Tweet)})
	}
	return tweets, loginEvent, nil
}
