package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

var ping_command Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Returns the ping of the bot",
	},
	Interact: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		start := time.Now()

		// Create a new HTTP client with a timeout
		client := http.Client{
			Timeout: 5 * time.Second,
		}

		// Send a GET request to the Discord API
		resp, err := client.Get("https://discord.com/api/v9/gateway/bot")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()

		// Calculate the ping
		ping := time.Since(start)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Ping:  %.2fms", ping.Seconds()*1000),
			},
		})
	},
}
