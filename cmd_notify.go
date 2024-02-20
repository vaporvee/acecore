package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// disabled
var notify_command Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:        "notify",
		Description: "Manage social media notifications.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add",
				Description: "Set channels where your social media notifications should appear.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "platform",
						Type:        discordgo.ApplicationCommandOptionString,
						Description: "The social media platform to receive notifications from.",
						Required:    true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{
								Name:  "Twitch",
								Value: "twitch",
							},
						},
					},
					{
						Name:        "username",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
						Description: "The social media platform to receive notifications from.",
					},
					{
						Name:        "channel",
						Type:        discordgo.ApplicationCommandOptionChannel,
						Required:    true,
						Description: "The social media platform to receive notifications from.",
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "remove",
				Description: "Remove a social media notification.",
			},
		},
	},
	Interact: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		options := i.ApplicationCommandData().Options[0]
		switch options.Name {
		case "add":
			switch options.Options[0].Value {
			case "twitch":
				fmt.Print("twitch")
			}
		}
	},
}
