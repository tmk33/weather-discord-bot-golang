package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"strings"

	"github.com/bwmarrin/discordgo"
	"weather-discord-bot-golang/weather"
)

// Variables used for command line parameters
var (
	Token string
)


func main() {
	
	Token = os.Getenv("DISCORD_TOKEN")
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	
	if !strings.HasPrefix(m.Content, "!bot") {
		return
	}
	
	args := strings.Split(m.Content, " ")[1:]
	
	
	switch args[0] {
	case "ping":
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	case "weather":
		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "COMMAND: !bot weather [location]")
		} else {
			resName, resTempC, resCondition  :=  weather.GetWeather(args[1])
			res := fmt.Sprintf("Location: %+v\nTemperature: %+v°C\nCondition: %+v", resName, resTempC, resCondition)
			s.ChannelMessageSend(m.ChannelID, res)
		}
	default:
		s.ChannelMessageSend(m.ChannelID, "Invalid command!")
	}
}
