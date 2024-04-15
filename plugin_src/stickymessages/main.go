package main

import (
	"database/sql"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
	"github.com/sirupsen/logrus"
	"github.com/vaporvee/acecore/custom"
	"github.com/vaporvee/acecore/shared"
)

var db *sql.DB

var dbCreateQuery string = `CREATE TABLE IF NOT EXISTS sticky (
	message_id TEXT NOT NULL,
	channel_id TEXT NOT NULL,
	message_content TEXT NOT NULL,
	guild_id TEXT NOT NULL,
	PRIMARY KEY (channel_id, guild_id)
);
`

var Plugin = &shared.Plugin{
	Name: "Sticky Messages",
	Init: func(d *sql.DB) error {
		db = d
		_, err := d.Exec(dbCreateQuery)
		if err != nil {
			return err
		}
		shared.BotConfigs = append(shared.BotConfigs, bot.WithEventListenerFunc(messageCreate))
		return nil
	},
	Commands: []shared.Command{
		{
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
		},
		{
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
		},
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
