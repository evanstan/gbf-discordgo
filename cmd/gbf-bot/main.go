package main

import (
	"github.com/evanstan/gbf-discordgo"
	"github.com/evanstan/gbf-discordgo/config"
)

func main() {
	s := &gbfbot.GBFBot{
		Config: config.LoadConfig(),
	}

	s.StartSession()

	<-make(chan struct{})

	s.CloseSession()
}