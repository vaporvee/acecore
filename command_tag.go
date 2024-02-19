package main

import (
	"fmt"
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
			{
				Name:        "remove",
				Description: "A command to remove messages saved to the bot.",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:         discordgo.ApplicationCommandOptionString,
						Name:         "tag",
						Description:  "The tag you want to remove",
						Required:     true,
						Autocomplete: true,
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
	AutocompleteTag(s, i)
	if i.Type == discordgo.InteractionApplicationCommand {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: getTagContent(i.GuildID, option.Value.(string)),
			},
		})
	}
}

func AutocompleteTag(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommandAutocomplete {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: generateDynamicChoices(i.GuildID),
			},
		})
	}
}

func generateDynamicChoices(guildID string) []*discordgo.ApplicationCommandOptionChoice {
	choices := []*discordgo.ApplicationCommandOptionChoice{}
	IDs, err := getTagIDs(guildID)
	if err != nil {
		log.Println("Error getting tag keys:", err)
		return choices
	}
	for _, id := range IDs {
		id_name := getTagName(guildID, id)
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  id_name,
			Value: id,
		})
	}
	return choices
}

// Yeeeahh the codebase sucks rn
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
	case "remove":
		AutocompleteTag(s, i)
		if i.Type == discordgo.InteractionApplicationCommand {
			fmt.Println("Trying to remove " + i.ApplicationCommandData().Options[0].Options[0].StringValue()) // so now it returns the content so wee reeeeaally need to start using UUIDs
			removeTag(i.GuildID, i.ApplicationCommandData().Options[0].Options[0].StringValue())
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Tag removed!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}
}

func (short_get_tag_command Command) tInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	GetTagCommand(s, i, i.ApplicationCommandData().Options[0])
}
