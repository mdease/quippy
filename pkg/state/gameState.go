package state

import (
	"strconv"

	"../db"
	"github.com/bwmarrin/discordgo"
)


// Maps from user ID to list of unmatched questions
var activeQuestions map[string][]*db.Prompt = make(map[string][]*db.Prompt)

// Maps from user ID to list of pending questions
var pendingQuestions map[string][]*db.Prompt = make(map[string][]*db.Prompt)

// Maps from user ID to number of answered but pending questions
var numAnswered map[string]int = make(map[string]int)

// Maps from guild ID to user ID to score
var scoreboard map[string]map[string]int = make(map[string]map[string]int)

// Send prompts to the user
func GetPrompts(user *discordgo.User, channelID string, num int) []*db.Prompt {
	prompts := []*db.Prompt{}

	// Check for open questions
	for u, rows := range activeQuestions {
		// Don't make the same user answer the same question twice
		if u != user.ID {
			// Use whole row or just some of it
			if len(rows) < num {
				prompts = append(prompts, rows...)
				activeQuestions[u] = []*db.Prompt{}
				num -= len(rows)
			} else {
				prompts = append(prompts, rows[:num]...)
				activeQuestions[u] = rows[num:]
				num = 0
			}
		}
	}

	// Fetch new questions
	for i := 0; i < num; i++ {
		row := db.Sample(channelID)
		prompts = append(prompts, row)
		activeQuestions[user.ID] = append(activeQuestions[user.ID], row)
	}

	// Mark questions as pending
	pendingQuestions[user.ID] = prompts[:]

	return prompts
}

// Get outstanding prompts for the user
func GetOutstandingPrompts(user *discordgo.User) []*db.Prompt {
	return pendingQuestions[user.ID]
}

// Record a response to a prompt
func RecordResponse(s *discordgo.Session, user *discordgo.User, index int, response string) {
	prompts := pendingQuestions[user.ID]
	index -= 1

	// Invalid index
	if index < 0 || index >= len(prompts) {
		return
	}

	prompt := prompts[index]

	// Check to make sure user hasn't already answered
	switch len(prompt.Answers) {
	case 0:
		prompt.Answers = append(prompt.Answers, response)
		prompt.RespondentIDs = append(prompt.RespondentIDs, user.ID)
		numAnswered[user.ID]++
	case 1:
		if prompt.RespondentIDs[0] != user.ID {
			prompt.Answers = append(prompt.Answers, response)
			prompt.RespondentIDs = append(prompt.RespondentIDs, user.ID)
			numAnswered[user.ID]++

			startVote(s, prompt)
		}
	default:
	}

	// If user has answered all pending questions, reset
	if len(pendingQuestions[user.ID]) == numAnswered[user.ID] {
		pendingQuestions[user.ID] = []*db.Prompt{}
		numAnswered[user.ID] = 0
	}
}

// Start a vote on a complete prompt
func startVote(s *discordgo.Session, p *db.Prompt) {
	// Get the responses and the users who submitted them
	response1 := p.Answers[0]
	response2 := p.Answers[1]

	user1 := p.RespondentIDs[0]
	user2 := p.RespondentIDs[1]

	index := strconv.FormatInt(int64(p.Index), 10)

	// Send the vote message
	message, err := s.ChannelMessageSend(p.ChannelID,
		"**It's time to vote** \xE2\x80\xBC\n\n" +
		"**Prompt**: (#" + index + ") " + p.Question + "\n" +
		"\x31\xE2\x83\xA3: " + response1 + " (" + "<@" + user1 + ">" + ")\n" +
		"\x32\xE2\x83\xA3: " + response2 + " (" + "<@" + user2 + ">" + ")\n\n" +
		"Vote by reacting with \x31\xE2\x83\xA3 and \x32\xE2\x83\xA3")

	if err != nil {
		return
	}

	// Add the vote reactions
	s.MessageReactionAdd(p.ChannelID, message.ID, "\x31\xE2\x83\xA3")
	s.MessageReactionAdd(p.ChannelID, message.ID, "\x32\xE2\x83\xA3")
}

// Record a user's score
func RecordScore(guildID string, userID string, score int) {
	// Might need to init guild map
	_, ok := scoreboard[guildID]

	if !ok {
		scoreboard[guildID] = make(map[string]int)
	}

	scoreboard[guildID][userID] += score
}

// Return a string of users and their scores
func GetScoreboard(guildID string) string {
	// Might need to init guild map
	guild, ok := scoreboard[guildID]

	if !ok {
		scoreboard[guildID] = make(map[string]int)
		return ""
	}

	// Make the scoreboard string
	scores := ""

	for user, score := range guild {
		score := "<@" + user + ">: " + strconv.FormatInt(int64(score), 10) + "\n"
		scores += score
	}

	return scores
}
