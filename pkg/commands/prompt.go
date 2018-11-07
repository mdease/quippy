package commands

import (
	"strconv"

	"../db"
	"../state"
	"github.com/bwmarrin/discordgo"
)


// Send prompts to the user via direct message
func Prompt(s *discordgo.Session, user *discordgo.User, args []string) string {
	// Send one prompt by default
	var numPrompts int64 = 1

	// Try to parse a number from the arguments
	if len(args) > 0 {
		n, err := strconv.ParseInt(args[0], 10, 64)

		if err == nil {
			numPrompts = n
		}
	}
	
	// Get the prompts
	prompts := state.GetPrompts(user, int(numPrompts))

	// Create a DM channel
	channel, err := s.UserChannelCreate(user.ID)

	if err != nil {
		panic("Could not create a direct message channel")
	}

	// Send the prompts
	s.ChannelMessageSend(channel.ID, createPromptMessage(prompts))

	// Reply with the number of prompts sent
	if numPrompts == 1 {
		return "1 prompt sent!"
	}

	numStr := strconv.FormatInt(numPrompts, 10)

	return numStr + " prompts sent!"
}

// Create the message to send with the prompts to the user
func createPromptMessage(prompts []*db.Row) string {
	reply := "**Here are your prompts**:\n"

	for i, p := range prompts {
		index := strconv.FormatInt(int64(i), 10)
		reply += index + ". " + p.Question + "\n"
	}

	reply += "\nRespond to these prompts with the prompt number followed by your response (e.g. 1 Stella and Oliver)"

	return reply
}
