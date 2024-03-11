package api

import (
	"encoding/json"
	"net/http"

	"github.com/masa-finance/masa-oracle/pkg/llmbridge"
	"github.com/masa-finance/masa-oracle/pkg/twitter"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

type SearchTweetsRequest struct {
	Query string `json:"query"`
	Count int    `json:"count"`
}

// SearchTweetsAndAnalyzeSentiment handles searching tweets and analyzing their sentiment
func SearchTweetsAndAnalyzeSentiment(w http.ResponseWriter, r *http.Request) {
	var reqBody SearchTweetsRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the request body
	if reqBody.Query == "" {
		http.Error(w, "Query parameter is missing", http.StatusBadRequest)
		return
	}
	if reqBody.Count <= 0 {
		reqBody.Count = 50
	}

	// Initialize a new Twitter scraper
	scraper := twitterscraper.New()

	// Fetch tweets using the Twitter API
	tweets, err := twitter.ScrapeTweetsByQuery(scraper, reqBody.Query, reqBody.Count, twitterscraper.SearchLatest)
	if err != nil {
		http.Error(w, "Failed to fetch tweets: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Analyze sentiment of tweets
	sentimentSummary, err := llmbridge.AnalyzeSentiment(tweets)
	if err != nil {
		http.Error(w, "Failed to analyze tweets: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with analysis results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"sentiment": sentimentSummary})
}
