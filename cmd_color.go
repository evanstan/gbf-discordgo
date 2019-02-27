package gbfbot

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"os"
	"time"
	"math/rand"
	"image"
	"image/draw"
	"log"
	"image/png"
	"image/color"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ColorConfig struct {
	Admins map[string]string `yaml:"Admins"`
	Colors map[string]int `yaml:"Colors"`

}

type Roles struct {
	ID   string
	Name string
}

var (
	colorconfig ColorConfig
	CreatedRoles = map[string]map[string]Roles{}
	SpamChannel = "429307503537422336"
	HelpText = `Help for Color-Bot
<<PrintColors   https://nayu.moe/colors
<<NewColor   "Assign a random color to the current user"
<<NewColor ColorName   "Assign the specified color to the current user"
<<PreviewColor ColorName   "Post a preview image of the color"`
)

func init() {
	createConfig(&colorconfig)
}

func createConfig(conf *ColorConfig) {
	yamlFile, err := ioutil.ReadFile("colorconfig.yaml")
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func CheckAdmin(UserID string) (bool) {
	if _, ok := colorconfig.Admins[UserID]; ok {
		return true
	}
	return false
}

func LoadRoles(session *discordgo.Session, GuildID string) {
	GuildRoles, err := session.GuildRoles(GuildID)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	// Initialise nested map with GuildID as key
	CreatedRoles[GuildID] = map[string]Roles{}
	for _, Role := range GuildRoles {
		if _, ok := colorconfig.Colors[Role.Name]; ok {
			CreatedRoles[GuildID][Role.Name] = Roles{Role.ID, Role.Name}
			CreatedRoles[GuildID][Role.ID]   = Roles{Role.ID, Role.Name}
		}
	}
}

func JoinedNewGuild(session *discordgo.Session, GuildID string) {
	// Initialise nested map with GuildID as key
	CreatedRoles[GuildID] = map[string]Roles{}
	fmt.Printf("Joined a new server: %s\n", GuildID)
	CreateAllRoles(session, GuildID)
}

func AddColorToAllMember(session *discordgo.Session, GuildID string) {
	fmt.Printf("Updating all member with new a color.\n")
	session.RequestGuildMembers(GuildID, "", 0)
}

func RemoveAllColors(session *discordgo.Session, GuildID string) {
	GuildRoles, err := session.GuildRoles(GuildID)
	if err != nil {
		fmt.Println("Permission error")
		return
	}

	for _, Role := range GuildRoles {
		if _, ok := colorconfig.Colors[Role.Name]; ok {
			fmt.Println("Remove role " + Role.Name + " " + Role.ID)
			err = session.GuildRoleDelete(GuildID, Role.ID)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func UpdateMemberColor(session *discordgo.Session, GuildID, MemberID, RoleName string) {
	err := session.GuildMemberRoleAdd(GuildID, MemberID, CreatedRoles[GuildID][RoleName].ID)
	if err != nil {
		LoadRoles(session, GuildID)
		UpdateMemberColorRandom(session, GuildID, MemberID)
	}
}

func UpdateMemberColorRandom(session *discordgo.Session, GuildID, MemberID string) {
	key := rand.Intn(len(colorconfig.Colors))
	randcolor := ""
	for Name, _ := range colorconfig.Colors {
		key -= 1
		if key == 0 {
			randcolor = Name
			break;
		}
	}
	session.GuildMemberRoleAdd(GuildID, MemberID, CreatedRoles[GuildID][randcolor].ID)
}

func CreateColorRole(session *discordgo.Session, GuildID, Name string, Color int) {
	role, err := session.GuildRoleCreate(GuildID)
	if err != nil {
		fmt.Println("Permission error")
		return
	}

	fmt.Printf("Name: %s     Int: %d \n", Name, Color)
	Role, _ := session.GuildRoleEdit(GuildID, role.ID, Name, Color, false, 0, false)
	CreatedRoles[GuildID][Role.Name] = Roles{Role.ID, Role.Name}
	CreatedRoles[GuildID][Role.ID]   = Roles{Role.ID, Role.Name}
}

func CreateNewRoles(session *discordgo.Session, GuildID string) {
	for Name, Color := range colorconfig.Colors {
		if _, ok := CreatedRoles[GuildID][Name]; !ok {
			CreateColorRole(session, GuildID, Name, Color)
		}
	}
}

func CreateAllRoles(session *discordgo.Session, GuildID string) {
	for Name, Color := range colorconfig.Colors {
		CreateColorRole(session, GuildID, Name, Color)
	}
}

func RemoveColorFromMember(session *discordgo.Session, GuildID, MemberID string) (bool) {
	Member, err := session.GuildMember(GuildID, MemberID)
	if err != nil {
		fmt.Printf("Can't get the guild.\n")
		fmt.Printf("Error:\n%s", err)
		return true
	}

	for _, RoleID := range Member.Roles {
		if _, ok := CreatedRoles[GuildID][RoleID]; ok {
			session.GuildMemberRoleRemove(GuildID, MemberID, RoleID)
		}
	}
	return false
}

func PreviewRole(session *discordgo.Session, RoleName string) discordgo.MessageEmbed {
	CreateImageWithColor(colorconfig.Colors[RoleName], RoleName)
	Embed := CreateImageEmbed(session, RoleName)

	return Embed
}

// Create Preview Image
func CreateImageWithColor(ColorInt int, ColorName string) {
	size := image.Rect(0, 0, 200, 100)
	rgbaImage := image.NewRGBA(size)
	red := uint8((ColorInt >> 16) & 0xff)
	green := uint8((ColorInt >> 8) & 0xff)
	blue := uint8(ColorInt & 0xff)
	c := color.RGBA{R:red, G:green, B:blue, A:255}
	draw.Draw(rgbaImage, rgbaImage.Bounds(), &image.Uniform{c}, image.ZP, draw.Src)

	f, err := os.Create(ColorName + ".png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, rgbaImage); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func CreateImageEmbed(session *discordgo.Session, ColorName string) discordgo.MessageEmbed {
	Embed := discordgo.MessageEmbed{Title: ColorName, Color: colorconfig.Colors[ColorName]}
	FileReader, _ := os.Open(ColorName + ".png")
	Msg, err := session.ChannelFileSend(SpamChannel, ColorName + ".png", FileReader)
	if err != nil {
		log.Fatal(err)
		return Embed
	}
	Image := discordgo.MessageEmbedImage{URL: Msg.Attachments[0].URL, Height: 100, Width: 200}
	Embed.Image = &Image

	return Embed
}

func KickMemberAfterTime(session *discordgo.Session, GuildID, MemberID string) {
	time.Sleep(30 * time.Minute)

	Member, err := session.GuildMember(GuildID, MemberID)
	if err != nil {
		fmt.Printf("Member already leaved.\n")
		return
	}

	for _, RoleID := range Member.Roles {
		if _, ok := colorconfig.Colors[CreatedRoles[GuildID][RoleID].Name]; ok {
			continue
		} else {
			return
		}
	}

	err = session.GuildMemberDeleteWithReason(GuildID, MemberID, "Not enough roles after 30min.")
	if err != nil {
		return
	}

	PrivateChannel, err := session.UserChannelCreate(MemberID)
	if err != nil {
		fmt.Printf("Can't send the message.\n")
		fmt.Printf("Error:\n%s", err)
		return
	}

	session.ChannelMessageSend(PrivateChannel.ID, "You got kicked from the server. Please read the welcome channel.")
}

func DeleteMessageAfterTime(session *discordgo.Session, Message *discordgo.Message, Time time.Duration) {
	time.Sleep(Time * time.Minute)
	session.ChannelMessageDelete(Message.ChannelID, Message.ID)
}

func SendMessageAndDeleteAfterTime(session *discordgo.Session, ChannelID, Content string) {
	Message, err := session.ChannelMessageSend(ChannelID, Content)
	if err != nil {
		fmt.Printf("Can't send the message.\n")
		fmt.Printf("Error:\n%s", err)
		return
	}

	go DeleteMessageAfterTime(session, Message, 5)
}

func SendEmbedAndDeleteAfterTime(session *discordgo.Session, ChannelID string, Embed discordgo.MessageEmbed) {
	Message, err := session.ChannelMessageSendEmbed(ChannelID, &Embed)
	if err != nil {
		fmt.Printf("Can't send embed.\n")
		fmt.Printf("Error:\n%s", err)
		return
	}

	go DeleteMessageAfterTime(session, Message, 5)
}

func NewColor(session *discordgo.Session, m *discordgo.MessageCreate, color string) {
	Channel, _ := session.State.Channel(m.ChannelID)
	resp := RemoveColorFromMember(session, Channel.GuildID, m.Author.ID)
	if resp {
		return
	}

	if color == "" {
		UpdateMemberColorRandom(session, Channel.GuildID, m.Author.ID)
	} else {
		if _, ok := CreatedRoles[Channel.GuildID][color]; ok {
			UpdateMemberColor(session, Channel.GuildID, m.Author.ID, color)
		}else if _, ok := colorconfig.Colors[color]; ok {
			UpdateMemberColorRandom(session, Channel.GuildID, m.Author.ID)
			SendMessageAndDeleteAfterTime(session, m.ChannelID, "Color not found available on this server.")
		} else {
			UpdateMemberColorRandom(session, Channel.GuildID, m.Author.ID)
			SendMessageAndDeleteAfterTime(session, m.ChannelID, "Color not found, pls use <<PrintColors.")
		}
	}
	session.ChannelMessageDelete(m.ChannelID, m.ID)
}