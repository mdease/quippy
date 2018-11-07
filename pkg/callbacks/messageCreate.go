package callbacks

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)


const botPrefix string = "q."

// Handles a MessageCreate Discord event
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot messages
	if m.Author.Bot {
		return
	}

	// Treat guild and direct messages differently
	// Guild messages must be valid commands
	// Direct messages may be in response to prompts
	if m.GuildID == "" {
		handleDirectMessage(s, m)
	} else if strings.HasPrefix(m.Content, botPrefix) {
		handleGuildMessage(s, m)
	}
}

// Handle a direct Discord message
func handleDirectMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Construct a reply message
	s.ChannelMessageSend(m.ChannelID, "hello :)")
}

// Handle a Discord guild message
func handleGuildMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Construct a reply message
	s.ChannelMessageSend(m.ChannelID, "hi :)")
}
