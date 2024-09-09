package twitterscraper

import (
	"log"
	"net/url"
)

type Response struct {
	Data struct {
		User struct {
			Result struct {
				Timeline struct {
					Timeline struct {
						Instructions []struct {
							Entries []struct {
								Content struct {
									ItemContent struct {
										UserResults struct {
											Result struct {
												Legacy Legacy `json:"legacy"`
											} `json:"result"`
										} `json:"user_results"`
									} `json:"itemContent"`
								} `json:"content"`
							} `json:"entries"`
						} `json:"instructions"`
					} `json:"timeline"`
				} `json:"timeline"`
			} `json:"result"`
		} `json:"user"`
	} `json:"data"`
}

type Legacy struct {
	CanDM                bool     `json:"can_dm"`
	CanMediaTag          bool     `json:"can_media_tag"`
	CreatedAt            string   `json:"created_at"`
	DefaultProfile       bool     `json:"default_profile"`
	DefaultProfileImage  bool     `json:"default_profile_image"`
	Description          string   `json:"description"`
	Entities             Entities `json:"entities"`
	FastFollowersCount   int      `json:"fast_followers_count"`
	FavouritesCount      int      `json:"favourites_count"`
	FollowersCount       int      `json:"followers_count"`
	FriendsCount         int      `json:"friends_count"`
	HasCustomTimelines   bool     `json:"has_custom_timelines"`
	IsTranslator         bool     `json:"is_translator"`
	ListedCount          int      `json:"listed_count"`
	Location             string   `json:"location"`
	MediaCount           int      `json:"media_count"`
	Name                 string   `json:"name"`
	NormalFollowersCount int      `json:"normal_followers_count"`
	PinnedTweetIdsStr    []string `json:"pinned_tweet_ids_str"`
	PossiblySensitive    bool     `json:"possibly_sensitive"`
	ProfileBannerUrl     string   `json:"profile_banner_url"`
	ProfileImageUrlHttps string   `json:"profile_image_url_https"`
	ScreenName           string   `json:"screen_name"`
	StatusesCount        int      `json:"statuses_count"`
	TranslatorType       string   `json:"translator_type"`
	Url                  string   `json:"url"`
	Verified             bool     `json:"verified"`
	WantRetweets         bool     `json:"want_retweets"`
	WithheldInCountries  []string `json:"withheld_in_countries"`
}

type Entities struct {
	Description Description `json:"description"`
	Url         Url         `json:"url"`
}

type Description struct {
	Urls []UrlInfo `json:"urls"`
}

type Url struct {
	Urls []UrlInfo `json:"urls"`
}

type UrlInfo struct {
	DisplayUrl  string `json:"display_url"`
	ExpandedUrl string `json:"expanded_url"`
	Url         string `json:"url"`
	Indices     []int  `json:"indices"`
}

// FetchFollowers gets the list of followers for a given user, via the Twitter frontend GraphQL API.
func (s *Scraper) FetchFollowers(username string, maxUsersNbr int, cursor string) ([]Legacy, string, error) {
	if maxUsersNbr > 200 {
		maxUsersNbr = 200
	}

	// Use GetProfile to get the userID from the username
	profile, err := s.GetProfile(username)
	if err != nil {
		return nil, "", err
	}
	userID := profile.UserID

	req, err := s.newRequest("GET", "https://twitter.com/i/api/graphql/o1YfmoGa-hb8Z6yQhoIBhg/Followers")
	if err != nil {
		return nil, "", err
	}

	variables := map[string]interface{}{
		"userId":                 userID,
		"count":                  maxUsersNbr,
		"includePromotedContent": false,
	}
	features := map[string]interface{}{"rweb_tipjar_consumption_enabled": true, "responsive_web_graphql_exclude_directive_enabled": true, "verified_phone_label_enabled": false, "creator_subscriptions_tweet_preview_api_enabled": true, "responsive_web_graphql_timeline_navigation_enabled": true, "responsive_web_graphql_skip_user_profile_image_extensions_enabled": false, "communities_web_enable_tweet_community_results_fetch": true, "c9s_tweet_anatomy_moderator_badge_enabled": true, "articles_preview_enabled": true, "tweetypie_unmention_optimization_enabled": true, "responsive_web_edit_tweet_api_enabled": true, "graphql_is_translatable_rweb_tweet_is_translatable_enabled": true, "view_counts_everywhere_api_enabled": true, "longform_notetweets_consumption_enabled": true, "responsive_web_twitter_article_tweet_consumption_enabled": true, "tweet_awards_web_tipping_enabled": false, "creator_subscriptions_quote_tweet_preview_enabled": false, "freedom_of_speech_not_reach_fetch_enabled": true, "standardized_nudges_misinfo": true, "tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true, "tweet_with_visibility_results_prefer_gql_media_interstitial_enabled": true, "rweb_video_timestamps_enabled": true, "longform_notetweets_rich_text_read_enabled": true, "longform_notetweets_inline_media_enabled": true, "responsive_web_enhance_cards_enabled": false}

	if cursor != "" {
		variables["cursor"] = cursor
	}

	query := url.Values{}
	query.Set("variables", mapToJSONString(variables))
	query.Set("features", mapToJSONString(features))
	req.URL.RawQuery = query.Encode()

	var response Response
	err = s.RequestAPI(req, &response)
	if err != nil {
		// Handle the error, for example, log it or return it to the caller
		log.Printf("Error making API request: %v", err)
		return nil, "", err
	}

	legacies, nextCursor, err := response.parseFollowing()
	if err != nil {
		// Handle the parsing error
		log.Printf("Error parsing following response: %v", err)
		return nil, "", err
	}

	// If err is nil here, it means both the API request and parsing were successful
	return legacies, nextCursor, nil
}

func (fr Response) parseFollowing() ([]Legacy, string, error) {
	var legacies []Legacy
	log.Println("Starting to parse following...") // Log the start of the parsing process

	for _, instruction := range fr.Data.User.Result.Timeline.Timeline.Instructions {
		for _, entry := range instruction.Entries {
			// Append the address of Legacy struct to the slice
			legacies = append(legacies, entry.Content.ItemContent.UserResults.Result.Legacy)
		}
	}

	// Assuming the next cursor is part of your response, you need to extract it here.
	// This is a placeholder for where you would extract the cursor from your response.
	// Adjust this according to your actual JSON structure.
	nextCursor := "" // Placeholder: Extract the actual cursor from the response

	return legacies, nextCursor, nil
}
