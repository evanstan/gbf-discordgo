package gbfbot

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/arbovm/levenshtein"
	"github.com/bwmarrin/discordgo"
)

const (
	emojiSuggestions            = 8
	emojiSuggestionsMaxDistance = 3
)

// cmdEmoji handles the emoji command.
func (g *GBFBot) cmdEmoji(ds *discordgo.Session, m *discordgo.MessageCreate, emoji string) error {
	if emoji == "" {
		return errors.New("Please specify an emoji.")
	}

	var (
		err error
		f   *os.File
	)

	name := fmt.Sprintf("%s.png", emoji)
	file := path.Join(g.EmojiDir, name)

	f, err = os.Open(file)
	if err == nil {
		_, err = ds.ChannelFileSend(m.ChannelID, name, f)

		f.Close()

		return err
	}

	var (
		entries      []string
		names        []string
		similarNames []string
	)

	f, err = os.Open(g.EmojiDir)
	if err != nil {
		return err
	}
	defer f.Close()

	entries, err = f.Readdirnames(-1)
	if err != nil {
		return err
	}

	names = make([]string, 0, len(entries))
	for _, entry := range entries {
		if !strings.HasSuffix(entry, ".png") {
			continue
		}

		name := strings.TrimSuffix(entry, ".png")
		name = fmt.Sprintf("%s", name)

		names = append(names, name)
	}

	if len(names) > 0 {
		sort.StringSlice(names).Sort()

		namesByDistance := [emojiSuggestionsMaxDistance + 1][]string{}
		for i := 0; i <= emojiSuggestionsMaxDistance; i++ {
			namesByDistance[i] = make([]string, 0)
		}

		for i := 0; i < len(names); i++ {
			d := levenshtein.Distance(emoji, names[i])
			if d <= emojiSuggestionsMaxDistance {
				nameCode := fmt.Sprintf("`%s`", names[i])
				namesByDistance[d] = append(namesByDistance[d], nameCode)
				similarNames = append(similarNames, nameCode)
			}
		}

		similarNames = make([]string, 0, emojiSuggestions)
		for i := 0; i <= emojiSuggestionsMaxDistance; i++ {
			if len(similarNames) > 0 {
				break
			}

			for j := 0; j < len(namesByDistance[i]) && len(similarNames) < emojiSuggestions; j++ {
				similarNames = append(similarNames, namesByDistance[i][j])
			}
		}
	}

	if len(similarNames) == 0 {
		_, err = ds.ChannelMessageSend(m.ChannelID, "Unknown emoji.")

		return err
	}

	response := fmt.Sprintf("Did you mean: %s", strings.Join(similarNames, ", "))

	_, err = ds.ChannelMessageSend(m.ChannelID, response)

	return err
}