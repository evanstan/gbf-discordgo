package config

import (
	"os"
	"fmt"
)

type Config struct {
	Token string
	Prefix string
	EmojiDir string
	StrikeTime string
}

func LoadConfig() (config *Config) {
	config = new(Config)

	config.Token = os.Getenv("TOKEN")
	if config.Token == "" {
		fmt.Println("$TOKEN not set")
	}

	config.Prefix = os.Getenv("PREFIX")
	if config.Prefix == "" {
		fmt.Println("$PREFIX not set")
	}

	config.EmojiDir = os.Getenv("EMOJIDIR")
	if config.EmojiDir == "" {
		fmt.Println("$EMOJIDIR not set")
	}

	config.StrikeTime = os.Getenv("STRIKETIME")
	if config.StrikeTime == "" {
		fmt.Println("$STRIKETIME not set")
	}

	return config
}