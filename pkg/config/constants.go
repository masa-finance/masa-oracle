package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/spf13/viper"
)

// ModelType defines a type for model strings.
type ModelType string

// Define model constants.
const (
	ClaudeOpus                                 ModelType = "claude-3-opus"
	ClaudeOpus20240229                         ModelType = "claude-3-opus-20240229"
	ClaudeSonnet20240229                       ModelType = "claude-3-sonnet-20240229"
	ClaudeHaiku20240307                        ModelType = "claude-3-haiku-20240307"
	GPT4                                       ModelType = "gpt-4"
	GPT4o                                      ModelType = "gpt-4o"
	GPT4TurboPreview                           ModelType = "gpt-4-turbo-preview"
	GPT35Turbo                                 ModelType = "gpt-3.5-turbo"
	LLama2                                     ModelType = "ollama/llama2"
	LLama3                                     ModelType = "ollama/llama3"
	Mistral                                    ModelType = "ollama/mistral"
	Gemma                                      ModelType = "ollama/gemma"
	Mixtral                                    ModelType = "ollama/mixtral"
	OpenChat                                   ModelType = "ollama/openchat"
	NeuralChat                                 ModelType = "ollama/neural-chat"
	CloudflareQwen15Chat                       ModelType = "@cf/qwen/qwen1.5-0.5b-chat"
	CloudflareLlama27bChatFp16                 ModelType = "@cf/meta/llama-2-7b-chat-fp16"
	CloudflareLlama38bInstruct                 ModelType = "@cf/meta/llama-3-8b-instruct"
	CloudflareMistral7bInstruct                ModelType = "@cf/mistral/mistral-7b-instruct"
	CloudflareMistral7bInstructV01             ModelType = "@cf/mistral/mistral-7b-instruct-v0.1"
	HuggingFaceGoogleGemma7bIt                 ModelType = "@hf/google/gemma-7b-it"
	HuggingFaceNousresearchHermes2ProMistral7b ModelType = "@hf/nousresearch/hermes-2-pro-mistral-7b"
	HuggingFaceTheblokeLlama213bChatAwq        ModelType = "@hf/thebloke/llama-2-13b-chat-awq"
	HuggingFaceTheblokeNeuralChat7bV31Awq      ModelType = "@hf/thebloke/neural-chat-7b-v3-1-awq"
	CloudflareOpenchat35_0106                  ModelType = "@cf/openchat/openchat-3.5-0106"
	CloudflareMicrosoftPhi2                    ModelType = "@cf/microsoft/phi-2"
)

// Models holds the available models for easy access and iteration.
var Models = struct {
	ClaudeOpus, ClaudeOpus20240229, ClaudeSonnet, ClaudeHaiku, GPT4, GPT4o, GPT4Turbo, GPT35Turbo, LLama2, LLama3, Mistral, Gemma, Mixtral, OpenChat, NeuralChat, CloudflareQwen15Chat, CloudflareLlama27bChatFp16, CloudflareLlama38bInstruct, CloudflareMistral7bInstruct, CloudflareMistral7bInstructV01, HuggingFaceGoogleGemma7bIt, HuggingFaceNousresearchHermes2ProMistral7b, HuggingFaceTheblokeLlama213bChatAwq, HuggingFaceTheblokeNeuralChat7bV31Awq, CloudflareOpenchat35_0106, CloudflareMicrosoftPhi2 ModelType
}{
	ClaudeOpus:                     ClaudeOpus,
	ClaudeOpus20240229:             ClaudeOpus20240229,
	ClaudeSonnet:                   ClaudeSonnet20240229,
	ClaudeHaiku:                    ClaudeHaiku20240307,
	GPT4:                           GPT4,
	GPT4o:                          GPT4o,
	GPT4Turbo:                      GPT4TurboPreview,
	GPT35Turbo:                     GPT35Turbo,
	LLama2:                         LLama2,
	LLama3:                         LLama3,
	Mistral:                        Mistral,
	Gemma:                          Gemma,
	Mixtral:                        Mixtral,
	OpenChat:                       OpenChat,
	NeuralChat:                     NeuralChat,
	CloudflareQwen15Chat:           CloudflareQwen15Chat,
	CloudflareLlama27bChatFp16:     CloudflareLlama27bChatFp16,
	CloudflareLlama38bInstruct:     CloudflareLlama38bInstruct,
	CloudflareMistral7bInstruct:    CloudflareMistral7bInstruct,
	CloudflareMistral7bInstructV01: CloudflareMistral7bInstructV01,
	HuggingFaceGoogleGemma7bIt:     HuggingFaceGoogleGemma7bIt,
	HuggingFaceNousresearchHermes2ProMistral7b: HuggingFaceNousresearchHermes2ProMistral7b,
	HuggingFaceTheblokeLlama213bChatAwq:        HuggingFaceTheblokeLlama213bChatAwq,
	HuggingFaceTheblokeNeuralChat7bV31Awq:      HuggingFaceTheblokeNeuralChat7bV31Awq,
	CloudflareOpenchat35_0106:                  CloudflareOpenchat35_0106,
	CloudflareMicrosoftPhi2:                    CloudflareMicrosoftPhi2,
}

const (
	PrivKeyFile = "MASA_PRIV_KEY_FILE"
	BootNodes   = "BOOTNODES"
	MasaDir     = "MASA_DIR"
	RpcUrl      = "RPC_URL"
	PortNbr     = "PORT_NBR"
	UDP         = "UDP"
	TCP         = "TCP"
	PrivateKey  = "PRIVATE_KEY"
	StakeAmount = "STAKE_AMOUNT"
	LogLevel    = "LOG_LEVEL"
	LogFilePath = "LOG_FILEPATH"
	Environment = "ENV"
	AllowedPeer = "allowedPeer"
	Signature   = "signature"
	Debug       = "debug"
	Version     = "v0.0.4-beta"
	FilePath    = "FILE_PATH"
	WriterNode  = "WRITER_NODE"
	CachePath   = "CACHE_PATH"

	MasaPrefix           = "/masa"
	OracleProtocol       = "oracle_protocol"
	NodeDataSyncProtocol = "nodeDataSync"
	NodeGossipTopic      = "gossip"
	AdTopic              = "ad"
	PublicKeyTopic       = "bootNodePublicKey"
	WorkerTopic          = "workerTopic"
	Rendezvous           = "masa-mdns"
	PageSize             = 25

	TwitterUsername  = "TWITTER_USERNAME"
	TwitterPassword  = "TWITTER_PASSWORD"
	Twitter2FaCode   = "TWITTER_2FA_CODE"
	DiscordBotToken  = "DISCORD_BOT_TOKEN"
	ClaudeApiKey     = "CLAUDE_API_KEY"
	ClaudeApiURL     = "CLAUDE_API_URL"
	ClaudeApiVersion = "CLAUDE_API_VERSION"
	GPTApiKey        = "OPENAI_API_KEY"
	TwitterScraper   = "TWITTER_SCRAPER"
	DiscordScraper   = "DISCORD_SCRAPER"
	WebScraper       = "WEB_SCRAPER"
	LlmServer        = "LLM_SERVER"
	LlmChatUrl       = "LLM_CHAT_URL"
	LlmCfUrl         = "LLM_CF_URL"
)

// ProtocolWithVersion returns a libp2p protocol ID string
// with the configured version and environment suffix.
func ProtocolWithVersion(protocolName string) protocol.ID {
	if GetInstance().Environment == "" {
		return protocol.ID(fmt.Sprintf("%s/%s/%s", MasaPrefix, protocolName, Version))
	}
	return protocol.ID(fmt.Sprintf("%s/%s/%s-%s", MasaPrefix, protocolName, Version, viper.GetString(Environment)))
}

// TopicWithVersion returns a topic string with the configured version
// and environment suffix.
func TopicWithVersion(protocolName string) string {
	if GetInstance().Environment == "" {
		return fmt.Sprintf("%s/%s/%s", MasaPrefix, protocolName, Version)
	}
	return fmt.Sprintf("%s/%s/%s-%s", MasaPrefix, protocolName, Version, viper.GetString(Environment))
}

// Function to call the Cloudflare API and parse the response
func GetCloudflareModels() ([]string, error) {
	url := "https://api.cloudflare.com/client/v4/accounts/a72433aa3bb83aecaca1bc8acecdb166/ai/models/search"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	bearer := fmt.Sprintf("Bearer %s", os.Getenv("LLM_CF_TOKEN"))
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Result []struct {
			ID string `json:"id"`
		} `json:"result"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	var models []string
	for _, model := range result.Result {
		models = append(models, model.ID)
	}

	return models, nil
}
