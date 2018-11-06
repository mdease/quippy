package main

import (
	"os"
	"time"

	"github.com/mdease/quippy/pkg/callbacks"
	"github.com/bwmarrin/discordgo"
)

func main() {
	token, found := os.LookupEnv("BOT_TOKEN")

	if !found || token == "" {
		panic("Bot token not found!")
	}

	session, err := discordgo.New("Bot " + token)

	if err != nil {
		panic("Could not create session!")
	}

	err = session.Open()

	if err != nil {
		panic("Could not establish session!")
	}

	session.AddHandler(callbacks.MessageCreate)

	time.Sleep(30 * time.Second)

	defer session.Close()
}
