package discord

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// BotConfig holds the configuration details for the Discord bot
type BotConfig struct {
	Token string
}

// LoadConfig reads configuration from environment variables or a .env file
func LoadConfig() *BotConfig {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading configuration from environment")
	}

	return &BotConfig{
		Token: os.Getenv("DISCORD_BOT_TOKEN"),
	}
}
