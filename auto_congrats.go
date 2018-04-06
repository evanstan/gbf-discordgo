package gbfbot

import (
	"math/rand"
	"time"
	
	"github.com/bwmarrin/discordgo"
)

func (g *GBFBot) sendCongrats(m *discordgo.MessageCreate) {

	messages := make([]string, 0)
	messages = append(messages,
		"I-I suppose you deserve that.",
		"Looks like you have gotten something nice, congratulations!",
		"Will you share some of that luck with Cagliostro?",
		)

	rand.Seed(time.Now().Unix())

	_, _ = g.session.ChannelMessageSend(m.ChannelID, messages[rand.Intn(len(messages))])
}
