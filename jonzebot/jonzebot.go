package main

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

//VERSION is the semantic version
const VERSION = "0.0.1"

// MessageCreate This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	rand.Seed(time.Now().UnixNano())

	for _, v := range commands {
		v.Process(s, m)
	}
}
