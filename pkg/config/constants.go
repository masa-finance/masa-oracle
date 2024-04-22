package config

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/spf13/viper"
)

// ModelType defines a type for model strings.
type ModelType string

// Define model constants.
const (
	ClaudeOpus20240229   ModelType = "claude-3-opus-20240229"
	ClaudeSonnet20240229 ModelType = "claude-3-sonnet-20240229"
	ClaudeHaiku20240307  ModelType = "claude-3-haiku-20240307"
	GPT4                 ModelType = "gpt-4"
	GPT4TurboPreview     ModelType = "gpt-4-turbo-preview"
	GPT35Turbo           ModelType = "gpt-3.5-turbo"
	LLama2               ModelType = "llama2"
	LLama3               ModelType = "llama3"
	Mistral              ModelType = "mistral"
	Gemma                ModelType = "gemma"
	Mixtral              ModelType = "mixtral"
	OpenChat             ModelType = "openchat"
	NeuralChat           ModelType = "neural-chat"
)

// Models holds the available models for easy access and iteration.
var Models = struct {
	ClaudeOpus, ClaudeSonnet, ClaudeHaiku, GPT4, GPT4Turbo, GPT35Turbo, LLama2, LLama3, Mistral, Gemma, Mixtral, OpenChat, NeuralChat ModelType
}{
	ClaudeOpus:   ClaudeOpus20240229,
	ClaudeSonnet: ClaudeSonnet20240229,
	ClaudeHaiku:  ClaudeHaiku20240307,
	GPT4:         GPT4,
	GPT4Turbo:    GPT4TurboPreview,
	GPT35Turbo:   GPT35Turbo,
	LLama2:       LLama2,
	LLama3:       LLama3,
	Mistral:      Mistral,
	Gemma:        Gemma,
	Mixtral:      Mixtral,
	OpenChat:     OpenChat,
	NeuralChat:   NeuralChat,
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
	Version     = "v0.0.11-alpha"
	FilePath    = "FILE_PATH"
	WriterNode  = "WRITER_NODE"
	CachePath   = "CACHE_PATH"

	MasaPrefix           = "/masa"
	OracleProtocol       = "oracle_protocol"
	NodeDataSyncProtocol = "nodeDataSync"
	NodeGossipTopic      = "gossip"
	AdTopic              = "ad"
	PublicKeyTopic       = "bootNodePublicKey"
	CompletedWorkTopic   = "completedWork"
	Rendezvous           = "masa-mdns"
	PageSize             = 25

	TwitterUsername  = "TWITTER_USERNAME"
	TwitterPassword  = "TWITTER_PASSWORD"
	Twitter2FaCode   = "TWITTER_2FA_CODE"
	ClaudeApiKey     = "CLAUDE_API_KEY"
	ClaudeApiURL     = "CLAUDE_API_URL"
	ClaudeApiVersion = "CLAUDE_API_VERSION"
	GPTApiKey        = "OPENAI_API_KEY"
	TwitterScraper   = "TWITTER_SCRAPER"
	WebScraper       = "WEB_SCRAPER"
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
