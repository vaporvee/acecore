package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var tag_command Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:        "get",
		Description: "A command to get messages saved to the bot.",
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
}

func (tag_command Command) Interaction(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		if len(i.ApplicationCommandData().Options) > 0 {
			// Loop through the options and handle them
			for _, option := range i.ApplicationCommandData().Options {
				switch option.Name {
				case "tag":
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
}
