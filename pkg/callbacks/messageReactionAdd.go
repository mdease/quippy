package callbacks

import (
	"regexp"
	"strconv"

	"../db"
	"github.com/bwmarrin/discordgo"
)


const botID string = "439164276058488843"

// Handles a MessageReactionAdd Discord event
func MessageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	// Ignore reactions on messages not from this bot
	message, err := s.ChannelMessage(m.ChannelID, m.MessageID)

	if err != nil || message.Author.ID != botID || m.UserID == botID {
		return
	}

	// If the reaction is in a DM, it should be a thumbs up or down
	if m.GuildID == "" {
		promptRegex, err := regexp.Compile(`^[0-9]+`)

		if err == nil && promptRegex.FindString(message.Content) != "" {
			// Read the question index
			var start, numLen int

			for message.Content[start:start + 1] != "#" {
				start++
			}

			start++

			for message.Content[start + numLen:start + numLen + 1] != ")" {

				numLen++
			}

			index, err := strconv.ParseInt(message.Content[start:start + numLen], 10, 64)

			if err != nil {
				return
			}

			// Upvote or downvote
			switch m.Emoji.Name {
			case "\xF0\x9F\x91\x8D":
				db.Upvote(int(index))
			case "\xF0\x9F\x91\x8E":
				db.Downvote(int(index))
			default:
			}
		}
	} else { // Otherwise it is a vote

	}
}
