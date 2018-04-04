package main

import (
	"fmt"
	"os"

	"github.com/evanstan/gbf-discordgo"
	"github.com/evanstan/gbf-discordgo/config"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("$PORT must be set")
	}

	s := &gbfbot.GBFBot{
		Config: config.LoadConfig(),
	}

	s.StartSession()

	<-make(chan struct{})

	s.CloseSession()
}