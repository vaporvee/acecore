package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var cmd_tag Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:                     "tag",
		DefaultMemberPermissions: int64Ptr(discordgo.PermissionManageServer),
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
			GetTagCommand(s, i, i.ApplicationCommandData().Options[0].Options[0])
		case "add":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseModal,
				Data: &discordgo.InteractionResponseData{
					CustomID: "tag_add_modal" + i.Interaction.Member.User.ID,
					Title:    "Add a custom tag command",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:  "tag_add_modal_name",
									Label:     "Name",
									Style:     discordgo.TextInputShort,
									Required:  true,
									MaxLength: 20,
									Value:     "",
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:    "tag_add_modal_content",
									Label:       "Content",
									Placeholder: "Content that gets returned when the tag will be run",
									Style:       discordgo.TextInputParagraph,
									Required:    true,
									MaxLength:   2000,
									Value:       "",
								},
							},
						},
					},
				},
			})
		case "remove":
			removeTag(i.GuildID, i.ApplicationCommandData().Options[0].Options[0].StringValue())
			respondEphemeral(s, i.Interaction, "Tag removed!")
		}
	},
	ModalIDs: []string{"tag_add_modal"},
	ModalSubmit: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		tagName := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		tagContent := i.ModalSubmitData().Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		addTag(i.GuildID, tagName, tagContent)
		respondEphemeral(s, i.Interaction, "Tag \""+tagName+"\" added!")
	},
	Autocomplete: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		AutocompleteTag(s, i)
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
		GetTagCommand(s, i, i.ApplicationCommandData().Options[0])
	},
	Autocomplete: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		AutocompleteTag(s, i)
	},
}

func GetTagCommand(s *discordgo.Session, i *discordgo.InteractionCreate, option *discordgo.ApplicationCommandInteractionDataOption) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: getTagContent(i.GuildID, option.Value.(string)),
		},
	})
}

func AutocompleteTag(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
