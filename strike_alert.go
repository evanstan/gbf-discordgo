package gbfbot

import (
	"fmt"
	"strings"
	"math/rand"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/bwmarrin/discordgo"
)

func (g *GBFBot) StartStrikeAlert() {
	times := strings.SplitN(g.Config.StrikeTime, "_", 2)
	if len(times) == 0 {
		return
	}

	gocron.Every(1).Day().At(times[0]).Do(g.postStrikeAlert)
	gocron.Every(1).Day().At(times[1]).Do(g.postStrikeAlert)
	<- gocron.Start()
}

func (g *GBFBot) postStrikeAlert() {
	if g.session == nil {
		fmt.Println("Session does not exists")
		return 
	}

	var (
		err error
	)

	messages := make([]string, 0)
	messages = append(messages,
		"Stop lazing around and join some raids!",
		"Start grinding for those SSRs!",
		"Strike time is here!",
		"Time to go all out!",
		)

	rand.Seed(time.Now().Unix())

	em := &discordgo.MessageEmbed{
		Title: "Strike Time",
		Description: messages[rand.Intn(len(messages))],
		Image: &discordgo.MessageEmbedImage{
			URL:    GBFStikeImageURL,
		},
	}

	_, err = g.session.ChannelMessageSendEmbed("431069931476484108", em)
	if err != nil {
		fmt.Println(err.Error())
	}
}