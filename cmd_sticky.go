package main

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
	"github.com/sirupsen/logrus"
	"github.com/vaporvee/acecore/custom"
)

var cmd_sticky Command = Command{
	Definition: discord.SlashCommandCreate{
		Name:                     "sticky",
		Description:              "Stick or unstick messages to the bottom of the current channel",
		DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionManageMessages),
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeGuild,
			discord.InteractionContextTypePrivateChannel},
		IntegrationTypes: []discord.ApplicationIntegrationType{
			discord.ApplicationIntegrationTypeGuildInstall},
		Options: []discord.ApplicationCommandOption{
			&discord.ApplicationCommandOptionString{
				Name:        "message",
				Description: "The message you want to stick to the bottom of this channel",
				Required:    false,
			},
		},
	},
	Interact: func(e *events.ApplicationCommandInteractionCreate) {
		if len(e.SlashCommandInteractionData().Options) == 0 {
			if hasSticky(e.GuildID().String(), e.Channel().ID().String()) {
				err := e.Client().Rest().DeleteMessage(e.Channel().ID(), snowflake.MustParse(getStickyMessageID(e.GuildID().String(), e.Channel().ID().String())))
				if err != nil {
					logrus.Error(err)
				}
				removeSticky(e.GuildID().String(), e.Channel().ID().String())
				err = e.CreateMessage(discord.NewMessageCreateBuilder().
					SetContent("The sticky message was removed from this channel!").SetEphemeral(true).
					Build())
				if err != nil {
					logrus.Error(err)
				}
			} else {
				err := e.CreateMessage(discord.NewMessageCreateBuilder().
					SetContent("This channel has no sticky message!").SetEphemeral(true).
					Build())
				if err != nil {
					logrus.Error(err)
				}
			}
		} else {
			inputStickyMessage(e)
		}
	},
}

var context_sticky Command = Command{
	Definition: discord.MessageCommandCreate{
		Name:                     "Stick to channel",
		DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionManageMessages),
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeGuild,
			discord.InteractionContextTypePrivateChannel},
		IntegrationTypes: []discord.ApplicationIntegrationType{
			discord.ApplicationIntegrationTypeGuildInstall},
	},
	Interact: func(e *events.ApplicationCommandInteractionCreate) {
		inputStickyMessage(e)
	},
}

func inputStickyMessage(e *events.ApplicationCommandInteractionCreate) {
	var messageText string
	if e.ApplicationCommandInteraction.Data.Type() == discord.ApplicationCommandTypeMessage {
		messageText = e.MessageCommandInteractionData().TargetMessage().Content //TODO add more data then just content
	} else {
		messageText = e.SlashCommandInteractionData().String("message")
	}
	if messageText == "" {
		err := e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Can't add empty sticky messages!").SetEphemeral(true).
			Build())
		if err != nil {
			logrus.Error(err)
		}
	} else {
		message, err := e.Client().Rest().CreateMessage(e.Channel().ID(), discord.MessageCreate{Embeds: []discord.Embed{
			{Description: messageText, Footer: &discord.EmbedFooter{Text: "📌 Sticky message"}, Color: custom.GetColor("primary")}}})
		if err != nil {
			logrus.Error(err)
		}

		if hasSticky(e.GuildID().String(), e.Channel().ID().String()) {
			err = e.Client().Rest().DeleteMessage(e.Channel().ID(), snowflake.MustParse(getStickyMessageID(e.GuildID().String(), e.Channel().ID().String())))
			if err != nil {
				logrus.Error(err, getStickyMessageID(e.GuildID().String(), e.Channel().ID().String()))
			}
			removeSticky(e.GuildID().String(), e.Channel().ID().String())
			addSticky(e.GuildID().String(), e.Channel().ID().String(), messageText, message.ID.String())
			err = e.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("Sticky message in this channel was updated!").SetEphemeral(true).
				Build())
			if err != nil {
				logrus.Error(err)
			}
		} else {
			addSticky(e.GuildID().String(), e.Channel().ID().String(), messageText, message.ID.String())
			err := e.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("Message sticked to the channel!").SetEphemeral(true).
				Build())
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}
