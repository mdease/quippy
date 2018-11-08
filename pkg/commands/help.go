package commands


func Help(args []string) string {
	// If no argument provided, return the default help message
	if len(args) == 0 {
		return "**Available commands**:\n" +
		"`help`: Display this message\n" + 
		"`help [command]`: Display more information about a certain command\n" +
		"`prompt [n]`: Request `n` prompts to be sent to you\n" +
		"`scores`: Print the scoreboard for this guild"
	}

	// Only check the first argument
	switch arg := args[0]; arg {
	case "help":
		return "`help` displays information on commands and how to use them."
	case "prompt":
		return "`prompt` requests prompts to be sent to you via a DM. " +
			"Specify the number of prompts to send by including a number after the command (e.g. `prompt 3`). " +
			"Omitting the number causes one prompt to be sent. " +
			"Answer the prompts by replying to the DM."
	case "scores":
		return "`scores` displays the scoreboard for this guild."
	default:
		return "Command not found: `" + arg + "`"
	}
}
