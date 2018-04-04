package gbfbot

import (
	"fmt"

	"github.com/evanstan/go-gbf"
	"github.com/bwmarrin/discordgo"
)

// cmdEvents handles the events command.
func (g *GBFBot) cmdEvents(s *discordgo.Session, m *discordgo.MessageCreate) error {
	var (
		err error

		currentEvents  []*CachedEvent
		upcomingEvents []*CachedEvent
	)

	currentEvents, err = g.CurrentEvents()
	if err != nil {
		return err
	}

	upcomingEvents, err = g.UpcomingEvents()
	if err != nil {
		return err
	}

	totalEvents := len(currentEvents) + len(upcomingEvents)

	em := &discordgo.MessageEmbed{
		Title: "Events",
		URL:   gbf.WikiEventsURL.String(),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL:    GBFLogoImageURL,
			Width:  GBFLogoImageWidth,
			Height: GBFLogoImageHeight,
		},
		Provider: &discordgo.MessageEmbedProvider{
			Name: "Granblue Fantasy Wiki",
			URL:  gbf.WikiHomeURL.String(),
		},
		Fields: make([]*discordgo.MessageEmbedField, 0, totalEvents),
	}

	for _, event := range currentEvents {
		value := "-"

		if !event.Details.EndsAt.IsZero() {
			endsAt := event.Details.EndsAt.UTC().Format("2006-01-02 15:04:05 MST")
			value = fmt.Sprintf("Ends on %s", endsAt)
		}

		em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
			Name:   event.Event.Title,
			Value:  value,
			Inline: false,
		})
	}
	for _, event := range upcomingEvents {
		value := "-"

		if event.Details != nil && !event.Details.StartsAt.IsZero() {
			startsAt := event.Details.StartsAt.UTC().Format("2006-01-02 15:04:05 MST")
			value = fmt.Sprintf("Begins on %s", startsAt)
		}

		em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
			Name:   event.Event.Title,
			Value:  value,
			Inline: false,
		})
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, em)
	if err != nil {
		fmt.Println(err.Error())
	}

	return err
}
