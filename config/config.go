package config

import (
	"os"
	"fmt"
)

var (
	Token string
	Prefix string
	EmojiDir string
)

func LoadConfig() {
	Token = os.Getenv("TOKEN")
	if Token == "" {
		fmt.Println("$TOKEN not set")
	}
	fmt.Println("Token: " + Token)

	Prefix = os.Getenv("PREFIX")
	if Prefix == "" {
		fmt.Println("$PREFIX not set")
	}
	fmt.Println("Prefix: " + Prefix)

	EmojiDir = os.Getenv("EMOJIDIR")
	if EmojiDir == "" {
		fmt.Println("$EMOJIDIR not set")
	}
	fmt.Println("EmojiDir: " + EmojiDir)

	return
}