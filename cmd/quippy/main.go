package main

import (
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"../../pkg/db"
	"../../pkg/callbacks"
)

// Start the bot!
func main() {
	// Fetch the bot token
	token, found := os.LookupEnv("BOT_TOKEN")

	if !found || token == "" {
		panic("Bot token not found")
	}

	// Create a Discord session
	session, err := discordgo.New("Bot " + token)

	if err != nil {
		panic("Could not create session")
	}

	// Open the session
	err = session.Open()

	if err != nil {
		panic("Could not establish session")
	}

	// Load the database of questions
	db.Load()

	// Set Discord event handlers
	session.AddHandler(callbacks.MessageCreate)
	session.AddHandler(callbacks.MessageReactionAdd)

	// Keep the bot running
	for {
		time.Sleep(1 * time.Hour)
	}

	// TODO kill signal?

	// Remember to close the session
	defer session.Close()

	// Save any changes to the database
	defer db.Save()
}
