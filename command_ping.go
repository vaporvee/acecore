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

		client := http.Client{
			Timeout: 5 * time.Second,
		}

		resp, err := client.Get("https://discord.com/api/v9/gateway/bot")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()

		ping := time.Since(start)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Ping:  %.2fms", ping.Seconds()*1000),
			},
		})
	},
}
