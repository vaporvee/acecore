package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/json"
	"github.com/sirupsen/logrus"
)

var cmd_tag Command = Command{
	Definition: discord.SlashCommandCreate{
		Name:                     "tag",
		DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionManageGuild),
		Description:              "A command to show and edit saved presaved messages.",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        "get",
				Description: "A command to get messages saved to the bot.",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:         "tag",
						Description:  "Your predefined tag for the saved message",
						Required:     true,
						Autocomplete: true,
					},
				},
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        "add",
				Description: "A command to add messages saved to the bot.",
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        "remove",
				Description: "A command to remove messages saved to the bot.",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:         "tag",
						Description:  "The tag you want to remove",
						Required:     true,
						Autocomplete: true,
					},
				},
			},
		},
	},
	Interact: func(e *events.ApplicationCommandInteractionCreate) {
		switch *e.SlashCommandInteractionData().SubCommandName {
		case "get":
			GetTagCommand(e)
		case "add":
			AddTagCommand(i, "")
		case "remove":
			removeTag(i.GuildID, i.ApplicationCommandData().Options[0].Options[0].StringValue())
			err := respond(i.Interaction, "Tag removed!", true)
			if err != nil {
				logrus.Error(err)
			}
		}
	},
	ModalIDs: []string{"tag_add_modal"},
	ModalSubmit: func(e *events.ModalSubmitInteractionCreate) {
		tagName := e.Data.Text("tag_add_modal_name")
		tagContent := e.Data.Text("tag_add_modal_content")
		addTag(e.GuildID().String(), tagName, tagContent)
		err := e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Tag \"" + tagName + "\" added!").SetEphemeral(true).
			Build())
		if err != nil {
			logrus.Error(err)
		}
	},
	Autocomplete: func(e *events.AutocompleteInteractionCreate) {
		AutocompleteTag(e)
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

var context_tag Command = Command{
	Definition: discordgo.ApplicationCommand{
		Type:                     discordgo.MessageApplicationCommand,
		Name:                     "Save as tag",
		DefaultMemberPermissions: int64Ptr(discordgo.PermissionManageServer),
	},
	Interact: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		AddTagCommand(i, i.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID].Content)
	},
}

func GetTagCommand(e *events.ApplicationCommandInteractionCreate) {
	err := e.CreateMessage(discord.NewMessageCreateBuilder().SetContent(getTagContent(guildID, option)).SetEphemeral(true).Build())
	if err != nil {
		logrus.Error(err)
	}
}

func AddTagCommand(i *discordgo.InteractionCreate, prevalue string) {
	err := bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
							Value:       prevalue,
						},
					},
				},
			},
		},
	})
	if err != nil {
		logrus.Error(err)
	}
}

func AutocompleteTag(i *discordgo.InteractionCreate) {
	err := bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: generateTagChoices(i.GuildID),
		},
	})
	if err != nil {
		logrus.Error(err)
	}
}

func generateTagChoices(guildID string) []*discordgo.ApplicationCommandOptionChoice {
	choices := []*discordgo.ApplicationCommandOptionChoice{}
	IDs, err := getTagIDs(guildID)
	if err != nil {
		logrus.Error(err)
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
