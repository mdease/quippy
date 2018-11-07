package commands

import (
	"strconv"

	"../db"
	"../state"
	"github.com/bwmarrin/discordgo"
)


// Send prompts to the user via direct message
func Prompt(s *discordgo.Session, user *discordgo.User, args []string) string {
	// First check if the user has any outstanding prompts
	outstanding := state.GetOutstandingPrompts(user)

	if len(outstanding) > 0 {
		sendPromptMessage(s, user, outstanding)

		return "Please answer all of your prompts first! They have been sent to you again."
	}

	// Send one prompt by default
	var numPrompts int64 = 1

	// Try to parse a number from the arguments
	if len(args) > 0 {
		n, err := strconv.ParseInt(args[0], 10, 64)

		if err == nil {
			if n < 1 {
				n = 1
			}

			numPrompts = n
		}
	}
	
	// Get the prompts
	prompts := state.GetPrompts(user, int(numPrompts))

	// Send the prompts
	sendPromptMessage(s, user, prompts)

	// Reply with the number of prompts sent
	if numPrompts == 1 {
		return "1 prompt sent!"
	}

	numStr := strconv.FormatInt(numPrompts, 10)

	return numStr + " prompts sent!"
}

// Send a messages to the user with the prompts
func sendPromptMessage(s *discordgo.Session, user *discordgo.User, prompts []*db.Row) {
	// Create a DM channel
	channel, err := s.UserChannelCreate(user.ID)

	if err != nil {
		panic("Could not create a direct message channel")
	}

	// Send the prompts as separate messages
	if len(prompts) == 1 {
		s.ChannelMessageSend(channel.ID, "**Here is your prompt**:")
	} else {
		s.ChannelMessageSend(channel.ID, "**Here are your prompts**:")
	}

	for i, p := range prompts {
		index := strconv.FormatInt(int64(i + 1), 10)
		qindex := strconv.FormatInt(int64(p.Index), 10)
		prompt := index + ". (#" + qindex + ") " + p.Question

		message, err := s.ChannelMessageSend(channel.ID, prompt)

		if err != nil {
			return
		}

		// Add thumbs up and down buttons
		s.MessageReactionAdd(channel.ID, message.ID, "\xF0\x9F\x91\x8D")
		s.MessageReactionAdd(channel.ID, message.ID, "\xF0\x9F\x91\x8E")
	}

	s.ChannelMessageSend(channel.ID,
		"Respond to each prompt with its list position followed by your response (e.g. 1 Stella and Oliver).\n" +
		"React to each prompt with thumbs up or thumbs down to make it appear more or less often, respectively.")
}
