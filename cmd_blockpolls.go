package main

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/sirupsen/logrus"
)

var cmd_blockpolls Command = Command{
	Definition: discord.SlashCommandCreate{
		Name:        "block-polls",
		Description: "Toggle blocking polls from beeing posted in this channel.",
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeGuild,
			discord.InteractionContextTypePrivateChannel},
		IntegrationTypes: []discord.ApplicationIntegrationType{
			discord.ApplicationIntegrationTypeGuildInstall},
	},
	Interact: func(e *events.ApplicationCommandInteractionCreate) {
		if toggleBlockPolls(e.GuildID().String(), e.Channel().ID().String()) {
			err := e.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("Polls are now unblocked in " + discord.ChannelMention(e.Channel().ID())).SetEphemeral(true).
				Build())
			if err != nil {
				logrus.Error(err)
			}
		} else {
			err := e.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("Polls are now blocked in " + discord.ChannelMention(e.Channel().ID())).SetEphemeral(true).
				Build())
			if err != nil {
				logrus.Error(err)
			}
		}
	},
}
