package main

import (
	"github.com/bwmarrin/discordgo"
)

var cmd_dadjoke Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:        "dadjoke",
		Description: "Gives you a random joke that is as bad as your dad would tell them",
	},
	Interact: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: simpleGetFromAPI("joke", "https://icanhazdadjoke.com/").(string),
			},
		})
	},
}
