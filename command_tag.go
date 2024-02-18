package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/iancoleman/strcase"
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
			{
				Name:        "add",
				Description: "A command to add messages saved to the bot.",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "tag",
						Description: "Your tag for the saved message",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "content",
						Description: "Your content for the saved message",
						Required:    true,
					},
				},
			},
		},
	}}

var short_get_tag_command Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:        "g",
		Description: "A short command to get presaved messages.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "tag",
				Description:  "Your predefined tag for the saved message",
				Required:     true,
				Autocomplete: true,
			},
		},
	}}

func GetTagCommand(s *discordgo.Session, i *discordgo.InteractionCreate, option *discordgo.ApplicationCommandInteractionDataOption) {
	if i.Type == discordgo.InteractionApplicationCommandAutocomplete {
		choices := generateDynamicChoices(i.GuildID)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		})
	}
	if i.Type == discordgo.InteractionApplicationCommand {
		if option.Name == "tag" {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: option.Value.(string),
				},
			})
		}
	}
}

func generateDynamicChoices(guildID string) []*discordgo.ApplicationCommandOptionChoice {
	choices := []*discordgo.ApplicationCommandOptionChoice{}
	keys, err := getTagKeys(guildID)
	if err != nil {
		log.Println("Error getting tag keys:", err)
		return choices // Return empty choices if there's an error
	}

	for _, key := range keys {
		tagContent, err := getTag(guildID, key) // Assuming you have a getTag function
		if err != nil {
			log.Println("Error getting tag content:", err)
			continue // Skip this tag if there's an error
		}
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  key,
			Value: tagContent,
		})
	}

	return choices
}

func (tag_command Command) Interaction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Options[0].Name {
	case "get":
		GetTagCommand(s, i, i.ApplicationCommandData().Options[0].Options[0])
	case "add":
		option := i.ApplicationCommandData().Options[0]
		addTag(i.GuildID, strcase.ToSnake(option.Options[0].StringValue()) /*TODO: tag regex*/, option.Options[1].StringValue())
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Tag added!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}
}

func (short_get_tag_command Command) tInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	GetTagCommand(s, i, i.ApplicationCommandData().Options[0])
}
