package main

import (
	"fmt"
	//"net/http"
	"os"

	"github.com/evanstan/gbf-discordgo"
	"github.com/evanstan/gbf-discordgo/config"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("$PORT must be set")
	}

	config.LoadConfig()

	s := &gbfbot.GBFBot{
		Token:    config.Token,
		Prefix:   config.Prefix,
		EmojiDir: config.EmojiDir,
	}

	/*err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}*/

	s.StartSession()

	<-make(chan struct{})

	s.CloseSession()
}