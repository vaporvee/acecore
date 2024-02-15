package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var tag_command Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:        "tag",
		Description: "A command to show and edit saved presaved messages.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "get",
				Description: "A command to get messages saved to the bot.",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:         discordgo.ApplicationCommandOptionString,
						Name:         "tag",
						Description:  "Your predefined tag for the saved message",
						Required:     true,
						Autocomplete: true,
					},
				},
			},
		},
	}}

func (tag_command Command) Interaction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Options[0].Name {
	case "get":
		if i.Type == discordgo.InteractionApplicationCommandAutocomplete {
			commandUseCount++
			choices := generateDynamicChoices(commandUseCount)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{
					Choices: choices,
				},
			})
		}
		if i.Type == discordgo.InteractionApplicationCommand {
			option := i.ApplicationCommandData().Options[0].Options[0]
			if option.Name == "tag" {
				value := option.Value.(string)
				response := fmt.Sprintf("You provided the tag: %s", value)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: response,
					},
				})
			}
		}
	}
}
