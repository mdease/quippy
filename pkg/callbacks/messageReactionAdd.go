package callbacks

import (
	"regexp"
	"strconv"
	"strings"

	"../db"
	"../state"
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
		voteRegex, err := regexp.Compile(`^\*\*It's time to vote\*\*`)

		if err == nil && voteRegex.FindString(message.Content) != "" {
			// Count the votes
			votes1, err1 := s.MessageReactions(m.ChannelID, m.MessageID, "\x31\xE2\x83\xA3", 3)
			votes2, err2 := s.MessageReactions(m.ChannelID, m.MessageID, "\x32\xE2\x83\xA3", 3)

			if err1 != nil || err2 != nil {
				return
			}

			// If one response has 3 votes, find the winner
			if len(votes1) < 3 && len(votes2) < 3 {
				return
			}

			var winner string

			if winner = "\x31\xE2\x83\xA3"; len(votes2) > len(votes1) {
				winner = "\x32\xE2\x83\xA3"
			}

			// Announce the winner
			// Get the prompt string
			promptIndex := strings.Index(message.Content, ")")
			promptIndex += 2
			promptLength := 0

			for message.Content[promptIndex + promptLength:promptIndex + promptLength + 1] != "\n" {
				promptLength++
			}

			prompt := message.Content[promptIndex:promptIndex + promptLength]

			// Get the response string
			responseIndex := strings.Index(message.Content, winner)
			responseIndex += 6
			responseLength := 0

			for message.Content[responseIndex + responseLength:responseIndex + responseLength + 1] != "\n" {
				responseLength++
			}

			response := message.Content[responseIndex:responseIndex + responseLength]

			// Get the user ID
			userIndex := strings.Index(response, "<@")
			userIndex += 2
			userLength := 0

			for response[userIndex + userLength:userIndex + userLength + 1] != ">" {
				userLength++
			}

			userID := response[userIndex:userIndex + userLength]

			// Calculate and record the score
			var score int
			if score = (len(votes1) - len(votes2)) * 200; score < 0 {
				score = -score
			}

			state.RecordScore(m.GuildID, userID, score)

			// Send the win message
			reply := winner + "wins!\n\n" +
				prompt + "\n" +
				response + "\n\n" +
				"Scores have been updated."

			s.ChannelMessageSend(m.ChannelID, reply)
		}
	}
}
