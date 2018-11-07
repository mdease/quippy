package commands


func Help(cmd string, args []string) string {
	// If no argument provided, return the default help message
	if len(args) == 0 {
		return "**Available commands**:\n" +
		"`help`: Display this message\n" + 
		"`help [command]`: Display more information about a certain command"
	}

	// Only check the first argument
	switch arg := args[0]; arg {
	case "help":
		return "`help` displays information on commands and how to use them"
	default:
		return "Command not found: `" + arg + "`"
	}
}
