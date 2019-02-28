package gbfbot

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/evanstan/gbf-discordgo/config"
)


//Main struct that wraps all the data required
type GBFBot struct {
	Config *config.Config

	eventsMutex    sync.Mutex
	currentEvents  []*CachedEvent
	upcomingEvents []*CachedEvent

	session *discordgo.Session
}

var (
	BotID string
	FirstTime bool = true
	HelpText = `!emo EmojiName "Post emoji"
!printcolors   "Post all available colors"
!newcolor   "Assign a random color to the current user"
!newcolor ColorName   "Assign the specified color to the current user"
!previewcolor ColorName   "Post a preview image of the color"`
	SpamChannel = "429307503537422336"
	TestChannel = "431069931476484108"
)

//Event handler
func (g *GBFBot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	//Congrats pock
	if(m.ChannelID == "431164963214852107" && m.Author.ID == "95625294064390144") {
		if(len(m.Attachments) > 0) {
			g.sendCongrats(m)
			return
		}
	}

	if g.Config.Prefix == "" {
		return
	}
	if !strings.HasPrefix(m.Content, g.Config.Prefix) {
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

	head = strings.TrimPrefix(parts[0], g.Config.Prefix)
	if len(parts) > 1 {
		tail = parts[1]
	}

	fmt.Println("head=%#v tail=%#v", head, tail)

	switch head {
	case "events":
		//err = g.cmdEvents(s, m)
	case "emo":
		err = g.cmdEmoji(s, m, tail)
	case "newcolor":
		NewColor(s, m, tail)
	case "previewcolor":
		PreviewColor(s, m, tail)
	case "printcolors":
		PrintColors(s, m)
	case "help":
		PostHelp(s, m)
	case "newserver":
		if CheckAdmin(m.Author.ID) {
			Channel, _ := s.State.Channel(m.ChannelID)
			JoinedNewGuild(s, Channel.GuildID)
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		}
	case "removeallcolors":
		if CheckAdmin(m.Author.ID) {
			Channel, _ := s.State.Channel(m.ChannelID)
			RemoveAllColors(s, Channel.GuildID)
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		}
	case "reloadcolors":
		if CheckAdmin(m.Author.ID) {
			Channel, _ := s.State.Channel(m.ChannelID)
			colorconfig = ColorConfig{}
			createConfig(&colorconfig)
			CreateNewRoles(s, Channel.GuildID)
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		}
	}



	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, err.Error())
	}
}

func onReady(session *discordgo.Session, Ready *discordgo.Ready) {
	if FirstTime {
		session.UpdateStatus(0, "Number one cutest ☆")

		for _, Guild := range Ready.Guilds {
			LoadRoles(session, Guild.ID)
			CreateNewRoles(session, Guild.ID)
		}
		fmt.Printf("Done.\n")
		FirstTime = false
	}
}

func OnMemberJoin(session *discordgo.Session, Member *discordgo.GuildMemberAdd) {
	UpdateMemberColorRandom(session, Member.GuildID, Member.User.ID)
}

func MemberChunkRequest(session *discordgo.Session, event *discordgo.GuildMembersChunk) {
	for _, Member := range event.Members {
		UpdateMemberColorRandom(session, event.GuildID, Member.User.ID)
	}
	fmt.Printf("Updated all members.\n")
}

//Start connection with Discord
func (g *GBFBot) StartSession() {
	if g.session != nil {
		fmt.Println("Session already exists")
		return 
	}

	s, err := discordgo.New("Bot " + g.Config.Token)
	if err != nil {
		fmt.Println(err.Error())
	}

	u, err := s.User("@me")
	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID

	s.AddHandler(g.messageHandler)
	s.AddHandler(onReady)
	s.AddHandler(OnMemberJoin)
	s.AddHandler(MemberChunkRequest)
	fmt.Println("Opening session...")
	err = s.Open()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Session opened successfully")

	g.session = s

	g.StartStrikeAlert()
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
}

func PostHelp(session *discordgo.Session, m * discordgo.MessageCreate) {
	if(m.ChannelID == SpamChannel || m.ChannelID == TestChannel) {
		em := discordgo.MessageEmbed{
			Title: "Cagloli Help!",
			Description: HelpText,
		}
		SendEmbedAndDeleteAfterTime(session, m.ChannelID, em)
	}
	session.ChannelMessageDelete(m.ChannelID, m.ID)
}