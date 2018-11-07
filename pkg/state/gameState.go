package state

import (
	"../db"
	"github.com/bwmarrin/discordgo"
)


// Maps from username to list of Row object
var activeQuestions map[string][]*db.Row = make(map[string][]*db.Row)

// Maps from username to list of Row objects
var pendingQuestions map[string][]*db.Row = make(map[string][]*db.Row)

// Send prompts to the user
func GetPrompts(user *discordgo.User, num int) []*db.Row {
	prompts := []*db.Row{}

	// Check for open questions
	for u, rows := range activeQuestions {
		// Don't make the same user answer the same question twice
		if u != user.ID {
			// Use whole row or just some of it
			if len(rows) < num {
				prompts = append(prompts, rows...)
				activeQuestions[u] = []*db.Row{}
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
		row := db.Sample()
		prompts = append(prompts, row)
		activeQuestions[user.ID] = append(activeQuestions[user.ID], row)
	}

	// Mark questions as pending
	pendingQuestions[user.ID] = prompts[:]

	return prompts
}
