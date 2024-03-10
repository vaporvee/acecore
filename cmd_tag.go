package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var cmd_tag Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:                     "tag",
		DefaultMemberPermissions: int64Ptr(discordgo.PermissionManageMessages),
		Description:              "A command to show and edit saved presaved messages.",
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
	},
	Interact: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.ApplicationCommandData().Options[0].Name {
		case "get":
			GetTagCommand(i, i.ApplicationCommandData().Options[0].Options[0])
		case "add":
			jsonStringShowModal(i.Interaction, "tag_add_modal", "template_general")
		case "remove":
			removeTag(i.GuildID, i.ApplicationCommandData().Options[0].Options[0].StringValue())
			respond(i.Interaction, "Tag removed!", true)
		}
	},
	ModalIDs: []string{"tag_add_modal"},
	ModalSubmit: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		tagName := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		tagContent := i.ModalSubmitData().Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		addTag(i.GuildID, tagName, tagContent)
		respond(i.Interaction, "Tag \""+tagName+"\" added!", true)
	},
	Autocomplete: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		AutocompleteTag(i)
	},
}

var cmd_tag_short Command = Command{
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
	},
	Interact: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		GetTagCommand(i, i.ApplicationCommandData().Options[0])
	},
	Autocomplete: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		AutocompleteTag(i)
	},
}

func GetTagCommand(i *discordgo.InteractionCreate, option *discordgo.ApplicationCommandInteractionDataOption) {
	respond(i.Interaction, getTagContent(i.GuildID, option.Value.(string)), false)
}

func AutocompleteTag(i *discordgo.InteractionCreate) {
	bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: generateTagChoices(i.GuildID),
		},
	})
}

func generateTagChoices(guildID string) []*discordgo.ApplicationCommandOptionChoice {
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
