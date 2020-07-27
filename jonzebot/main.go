package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Create a new Discord session using the provided bot token.
	if fileExists(os.Getenv("PLAYBILL")) {
		unarchiveJSON(os.Getenv("PLAYBILL"), &playbill)
	}
	if fileExists(os.Getenv("VIDEOVAULT")) {
		unarchiveJSON(os.Getenv("VIDEOVAULT"), &vault)
	}
	if fileExists(os.Getenv("NOWPLAYING")) {
		unarchiveJSON(os.Getenv("NOWPLAYING"), &currentMovie)
	}
	if currentMovie.Start.Local().Add(6*time.Hour).Unix() < time.Now().Unix() {
		currentMovie = NowPlaying{}
		archiveJSON(os.Getenv("NOWPLAYING"), &currentMovie)
	}
	dg, err := discordgo.New("Bot " + os.Getenv("JONZEID"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(MessageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
