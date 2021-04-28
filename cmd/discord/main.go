package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/marcustut/community-bot/pkg/utils"
)

var DISCORD_BOT_TOKEN string

func init() {
	// Load .env
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env")
	}

	// Get the bot token from environment
	DISCORD_BOT_TOKEN = os.Getenv("DISCORD_BOT_TOKEN")
}

func main() {
	// Create a new discord session
	discord, err := discordgo.New("Bot " + DISCORD_BOT_TOKEN)
	if err != nil {
		log.Fatal("error connecting with Discord", err)
	}

	// Add a handler to the bot
	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// If messages from bot
		if m.Author.ID == s.State.User.ID {
			return
		}

		// if message from outside of admin channel
		if m.ChannelID != utils.ChannelIDs["admin"] {
			return
		}

		log.Print(fmt.Sprintf("%s: %s", m.Author.Username, m.Content))
	})

	// Open a websocket connection to the bot
	err = discord.Open()
	if err != nil {
		log.Fatal("error opening connection, ", err)
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	// Open a channel to listen for system interrupt
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Close the connection if interrupted
	discord.Close()
}
