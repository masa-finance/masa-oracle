# Masa Node Twitter Sentiment Analysis Feature

The Masa Node introduces a powerful feature for analyzing the sentiment of tweets. This functionality leverages advanced language models to interpret the sentiment behind a collection of tweets, providing valuable insights into public perception and trends.

## Overview

The Twitter sentiment analysis feature is part of the broader capabilities of the Masa Node, designed to interact with social media data in a meaningful way. It uses state-of-the-art language models to evaluate the sentiment of tweets, categorizing them into positive, negative, or neutral sentiments.

## How It Works

The sentiment analysis process involves fetching tweets based on specific queries, and then analyzing these tweets using selected language models. The system supports various models, including Claude and GPT variants, allowing for flexible and powerful sentiment analysis.

### Fetching Tweets

Tweets are fetched using the Twitter Scraper library, as seen in the [llmbridge](file:///Users/john/Projects/masa/masa-oracle/pkg/llmbridge/sentiment_twitter.go#1%2C9-1%2C9) package. This process does not require Twitter API keys, making it accessible and straightforward.

```go
func AnalyzeSentiment(tweets []*twitterscraper.Tweet, model string) (string, string, error) { ... }
```

### Analyzing Sentiment

Once tweets are fetched, they are sent to the chosen language model for sentiment analysis. The system currently supports models prefixed with "claude-" and "gpt-", catering to a range of analysis needs.

### Integration with Masa Node CLI

The sentiment analysis feature is integrated into the Masa Node CLI, allowing users to interact with it directly from the command line. Users can specify the query, the number of tweets to analyze, and the model to use for analysis.

```go
var countMessage string
var userMessage string

inputCountField := tview.NewInputField().
  SetLabel("# of Tweets to analyze ").
  SetFieldWidth(10)
```

### Example Usage

o analyze the sentiment of tweets, users can follow these steps:

1. Launch the Masa Node CLI.
2. Navigate to the sentiment analysis section.
3. Enter the query for fetching tweets.
4. Specify the number of tweets to analyze.
5. Choose the language model for analysis.

The system will then display the sentiment analysis results, providing insights into the overall sentiment of the tweets related to the query.

### Conclusion

The Twitter sentiment analysis feature of the Masa Node offers a powerful tool for understanding public sentiment on various topics. By leveraging advanced language models, it provides deep insights into the emotional tone behind tweets, making it a valuable asset for data analysis and decision-making.
