package main

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/json"
	"github.com/sirupsen/logrus"
)

// TODO: make user installable tag command using userIDs instead of guildIDs
var cmd_tag Command = Command{
	Definition: discord.SlashCommandCreate{
		Name:                     "tag",
		DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionManageGuild),
		Description:              "A command to show and edit saved presaved messages.",
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeGuild,
			discord.InteractionContextTypePrivateChannel},
		IntegrationTypes: []discord.ApplicationIntegrationType{
			discord.ApplicationIntegrationTypeGuildInstall},
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
			AddTagCommand(e)
		case "remove":
			removeTag(e.GuildID().String(), e.SlashCommandInteractionData().String("tag"))
			err := e.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("Tag removed!").SetEphemeral(true).
				Build())
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
	Definition: discord.SlashCommandCreate{
		Name:        "g",
		Description: "A short command to get presaved messages.",
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeGuild,
			discord.InteractionContextTypePrivateChannel},
		IntegrationTypes: []discord.ApplicationIntegrationType{
			discord.ApplicationIntegrationTypeGuildInstall},
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:         "tag",
				Description:  "Your predefined tag for the saved message",
				Required:     true,
				Autocomplete: true,
			},
		},
	},
	Interact: func(e *events.ApplicationCommandInteractionCreate) {
		GetTagCommand(e)
	},
	Autocomplete: func(e *events.AutocompleteInteractionCreate) {
		AutocompleteTag(e)
	},
}

var context_tag Command = Command{
	Definition: discord.MessageCommandCreate{
		Name:                     "Save as tag",
		DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionManageGuild),
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeGuild,
			discord.InteractionContextTypePrivateChannel},
		IntegrationTypes: []discord.ApplicationIntegrationType{
			discord.ApplicationIntegrationTypeGuildInstall},
	},
	Interact: func(e *events.ApplicationCommandInteractionCreate) {
		AddTagCommand(e)
	},
}

func GetTagCommand(e *events.ApplicationCommandInteractionCreate) {
	err := e.CreateMessage(discord.NewMessageCreateBuilder().
		SetContent(getTagContent(e.GuildID().String(), e.SlashCommandInteractionData().String("tag"))).
		Build())
	if err != nil {
		logrus.Error(err)
	}
}

func AddTagCommand(e *events.ApplicationCommandInteractionCreate) {
	var prevalue string
	if e.ApplicationCommandInteraction.Data.Type() == discord.ApplicationCommandTypeMessage {
		prevalue = e.MessageCommandInteractionData().TargetMessage().Content
	}
	err := e.Modal(discord.ModalCreate{
		CustomID: "tag_add_modal" + e.User().ID.String(),
		Title:    "Add a custom tag command",
		Components: []discord.ContainerComponent{
			discord.ActionRowComponent{
				discord.TextInputComponent{
					CustomID:  "tag_add_modal_name",
					Label:     "Name",
					Style:     discord.TextInputStyleShort,
					Required:  true,
					MaxLength: 20,
					Value:     "",
				},
			},
			discord.ActionRowComponent{
				discord.TextInputComponent{
					CustomID:  "tag_add_modal_content",
					Label:     "Content",
					Style:     discord.TextInputStyleParagraph,
					Required:  true,
					MaxLength: 2000,
					Value:     prevalue,
				},
			},
		},
	})
	if err != nil {
		logrus.Error(err)
	}
}

func AutocompleteTag(e *events.AutocompleteInteractionCreate) {
	err := e.AutocompleteResult(generateTagChoices(e.GuildID().String()))
	if err != nil {
		logrus.Error(err)
	}
}

func generateTagChoices(guildID string) []discord.AutocompleteChoice {
	choices := []discord.AutocompleteChoice{}
	IDs, err := getTagIDs(guildID)
	if err != nil {
		logrus.Error(err)
		return choices
	}
	for _, id := range IDs {
		id_name := getTagName(guildID, id)
		choices = append(choices, &discord.AutocompleteChoiceString{
			Name:  id_name,
			Value: id,
		})
	}
	return choices
}
