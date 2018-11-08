package callbacks

import (
	"strconv"
	"strings"

	"../commands"
	"../state"
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
	// Get the prompt index and response
	message := strings.TrimSpace(m.Content)
	split := strings.Index(message, " ")

	if split < 0 {
		return
	}

	indexStr := m.Content[:split]
	index, err := strconv.ParseInt(indexStr, 10, 64)

	if err != nil {
		return
	}

	response := message[split + 1:]

	// Record the response
	state.RecordResponse(s, m.Author, int(index), response)
}

// Handle a Discord guild message
func handleGuildMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Get the command
	strs := strings.Split(m.Content, " ")
	cmd := strs[0][2:]
	args := strs[1:]

	// Construct a reply message
	reply := ""

	// Make sure the command is valid
	switch cmd {
	case "help":
		reply = commands.Help(args)
	case "prompt":
		reply = commands.Prompt(s, m.Author, m.ChannelID, args)
	case "scores":
		reply = commands.Scores(m.GuildID)
	default:
	}

	// Send the reply
	if len(reply) > 0 {
		s.ChannelMessageSend(m.ChannelID, reply)
	}
}
