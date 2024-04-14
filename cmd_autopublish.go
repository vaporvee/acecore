package main

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/json"
	"github.com/sirupsen/logrus"
)

var cmd_autopublish Command = Command{
	Definition: discord.SlashCommandCreate{
		Name:                     "autopublish",
		Description:              "Toggle automatically publishing every post in a announcement channel",
		DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionManageChannels),
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeGuild,
			discord.InteractionContextTypePrivateChannel},
		IntegrationTypes: []discord.ApplicationIntegrationType{
			discord.ApplicationIntegrationTypeGuildInstall},
	},
	Interact: func(e *events.ApplicationCommandInteractionCreate) {
		channel := e.Channel()
		if channel.Type() == discord.ChannelTypeGuildNews {
			if toggleAutoPublish(e.GuildID().String(), e.Channel().ID().String()) {
				err := e.CreateMessage(discord.NewMessageCreateBuilder().
					SetContent("Autopublishing is now disabled on " + discord.ChannelMention(e.Channel().ID())).SetEphemeral(true).
					Build())
				if err != nil {
					logrus.Error(err)
				}
			} else {
				err := e.CreateMessage(discord.NewMessageCreateBuilder().
					SetContent("Autopublishing is now enabled on " + discord.ChannelMention(e.Channel().ID())).SetEphemeral(true).
					Build())
				if err != nil {
					logrus.Error(err)
				}
			}
		} else {
			err := e.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("This is not an announcement channel!").SetEphemeral(true).
				Build())
			if err != nil {
				logrus.Error(err)
			}
		}
	},
}
