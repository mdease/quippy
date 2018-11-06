package callbacks

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)


const botPrefix string = "q."

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot messages
	if m.Author.Bot {
		return
	}

	// Ignore messages without bot prefix
	if !strings.HasPrefix(m.Content, botPrefix) {
		return
	}

	// Construct a reply message
	s.ChannelMessageSend(m.ChannelID, "hi :)")
}
