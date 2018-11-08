package commands

import "../state"


// Return a string representing the scoreboard
func Scores(guildID string) string {
	return "**Scores**:\n" + state.GetScoreboard(guildID)
}
