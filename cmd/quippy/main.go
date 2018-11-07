package main

import (
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"../../pkg/db"
	"../../pkg/callbacks"
)

// Start the bot!
// This bot is currently intended to run in only one guild.
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

	time.Sleep(60 * time.Second)

	// Remember to close the session
	defer session.Close()

	// Save any changes to the database
	defer db.Save()
}
