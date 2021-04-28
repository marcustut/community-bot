package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
	"github.com/marcustut/community-bot/pkg/utils"
)

type client struct {
	Address string
}

type Message struct {
	AvatarURL string `json:"avatarURL"`
	Text      string `json:"text"`
}

var (
	DISCORD_BOT_TOKEN string
	clients           = make(map[string]client)    // Store connected clients
	register          = make(chan string)          // To register client
	unregister        = make(chan string)          // To unregister client
	forward           = make(chan *Message)        // Channel for forwarding messages
	storeWebConn      = make(chan *websocket.Conn) // To store current web client's connection
	webConn           *websocket.Conn              // Current webClient's connection
)

func init() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	// Get the bot token from environment
	DISCORD_BOT_TOKEN = os.Getenv("DISCORD_BOT_TOKEN")
}

func runHub() {
	for {
		select {
		case address := <-register:
			clients[address] = client{Address: address}
			log.Printf("Client %s registered\n", clients[address].Address)

		case connection := <-storeWebConn:
			webConn = connection

		case address := <-unregister:
			// Remove client from the hub
			log.Printf("Client %s unregistered\n", clients[address].Address)
			delete(clients, address)

		case message := <-forward:
			log.Println("Received message: ", message)

			if webConn == nil {
				log.Println("Web client is not connected")
				return
			}

			// Send the message to websocket
			// err := webConn.WriteMessage(websocket.TextMessage, []byte(message))
			err := webConn.WriteJSON(message)
			if err != nil {
				log.Println("Error forwarding message: ", err)

				webConn.WriteMessage(websocket.CloseMessage, []byte{})
				webConn.Close()
				delete(clients, webConn.RemoteAddr().String())
			}

			log.Printf("Forwarded '%s' to %s\n", message, webConn.RemoteAddr().String())
		}
	}
}

func fiberWS() {
	app := fiber.New()

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/messages", websocket.New(func(c *websocket.Conn) {
		// When the function returns, close the connection
		defer func() {
			unregister <- c.RemoteAddr().String()
			c.Close()
		}()

		// Register the client
		register <- c.RemoteAddr().String()

		// Store the web client's connection
		storeWebConn <- c

		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error: ", err)
				}
				return // Calls the deferred function, ie closes the connection on error
			}

			log.Printf("Recevied message of %v: %s", messageType, message)
		}
	}))

	log.Fatal(app.Listen(":8000"))
}

func main() {
	// Runs the hub
	go runHub()

	// Runs the websocket server
	go fiberWS()

	// Create a new discord session
	discord, err := discordgo.New("Bot " + DISCORD_BOT_TOKEN)
	if err != nil {
		log.Fatal("error connecting with Discord", err)
	}

	// Connect to websocket server
	register <- "Community Discord Bot"

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

		if strings.Contains(m.Content, ";clear") {
			parsed := strings.Split(m.Content, " ")

			numLine, err := strconv.ParseInt(parsed[1], 10, 64)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "You did not specify a number")
				return
			}

			if numLine > 100 {
				s.ChannelMessageSend(m.ChannelID, "You cannot clear more than 100 messages at once")
				return
			}

			if numLine < 1 {
				s.ChannelMessageSend(m.ChannelID, "You have to clear at least 1 message")
				return
			}

			channel, err := s.Channel(m.ChannelID)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error clearing message")
				log.Fatalf("Channel %s is not fetched", m.ChannelID)
				return
			}

			messageIDs := make([]string, 100)

			// TODO: Fix the stat-cached channels bug
			for i := 0; int64(i) < numLine && int(64) < 100; i++ {
				fmt.Println(channel)
				// messageIDs = append(messageIDs, channel.Messages[i].ChannelID)
			}

			s.ChannelMessagesBulkDelete(m.ChannelID, messageIDs)
			return
		}

		log.Print(fmt.Sprintf("%s: %s", m.Author.Username, m.Content))
		forward <- &Message{AvatarURL: m.Author.AvatarURL("64"), Text: m.Content}
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
