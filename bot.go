package gbfbot

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)


//Main struct that wraps all the data required
type GBFBot struct {
	Token 	 string
	Prefix 	 string
	EmojiDir string

	eventsMutex    sync.Mutex
	currentEvents  []*CachedEvent
	upcomingEvents []*CachedEvent

	session *discordgo.Session
}

var (
	BotID string
)

//Event handler
func (g *GBFBot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	if g.Prefix == "" {
		return
	}
	if !strings.HasPrefix(m.Content, g.Prefix) {
		return
	}

	var (
		err error

		head string
		tail string
	)

	parts := strings.SplitN(m.Content, " ", 2)
	if len(parts) == 0 {
		return
	}

	head = strings.TrimPrefix(parts[0], g.Prefix)
	if len(parts) > 1 {
		tail = parts[1]
	}

	fmt.Println("head=%#v tail=%#v", head, tail)

	switch head {
	case "events":
		err = g.cmdEvents(s, m)
	case "emo":
		err = g.cmdEmoji(s, m, tail)
	}

	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, err.Error())
	}
}

//Start connection with Discord
func (g *GBFBot) StartSession() {
	if g.session != nil {
		fmt.Println("Session already exists")
		return 
	}

	s, err := discordgo.New("Bot " + g.Token)
	if err != nil {
		fmt.Println(err.Error())
	}

	u, err := s.User("@me")
	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID

	s.AddHandler(g.messageHandler)
	fmt.Println("Opening session...")
	err = s.Open()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Session opened successfully")

	g.session = s

	return
}

//Close connection with Discord
func (g *GBFBot) CloseSession() {
	if g.session == nil {
		fmt.Println("Session does not exists")
		return 
	}

	fmt.Println("Closing session")

	err := g.session.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Session Closed")

	return
}