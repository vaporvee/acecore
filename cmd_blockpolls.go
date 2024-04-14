package main

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/json"
	"github.com/sirupsen/logrus"
	"github.com/vaporvee/acecore/struct_cmd"
)

var cmd_blockpolls struct_cmd.Command = struct_cmd.Command{
	Definition: discord.SlashCommandCreate{
		Name:                     "block-polls",
		Description:              "Block polls from beeing posted in this channel.",
		DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionManageChannels),
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeGuild,
			discord.InteractionContextTypePrivateChannel},
		IntegrationTypes: []discord.ApplicationIntegrationType{
			discord.ApplicationIntegrationTypeGuildInstall},
		Options: []discord.ApplicationCommandOption{
			&discord.ApplicationCommandOptionSubCommand{
				Name:        "toggle",
				Description: "Toggle blocking polls from beeing posted in this channel.",
				Options: []discord.ApplicationCommandOption{
					&discord.ApplicationCommandOptionBool{
						Name:        "global",
						Description: "If polls are blocked server wide or only in the current channel.",
					},
					&discord.ApplicationCommandOptionRole{
						Name:        "allowed-role",
						Description: "The role that bypasses this block role.",
					},
				},
			},
			/*&discord.ApplicationCommandOptionSubCommand{
				Name:        "list",
				Description: "List the current block polls rules for this server.",
			},*/
		},
	},
	Interact: func(e *events.ApplicationCommandInteractionCreate) {
		switch *e.SlashCommandInteractionData().SubCommandName {
		case "toggle":
			isGlobal := isGlobalBlockPolls(e.GuildID().String())
			if isGlobal && !e.SlashCommandInteractionData().Bool("global") {
				e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Polls are currently globally blocked. Disable global blocking to enable channel specific blocking.").SetEphemeral(true).Build())
			} else {
				exists, isGlobal := toggleBlockPolls(e.GuildID().String(), e.Channel().ID().String(), e.SlashCommandInteractionData().Bool("global"), e.SlashCommandInteractionData().Role("allowed-role").ID.String())
				if exists {
					if e.SlashCommandInteractionData().Bool("global") {
						err := e.CreateMessage(discord.NewMessageCreateBuilder().
							SetContent("Polls are now globally unblocked.").SetEphemeral(true).
							Build())
						if err != nil {
							logrus.Error(err)
						}
					} else {
						err := e.CreateMessage(discord.NewMessageCreateBuilder().
							SetContent("Polls are now unblocked in " + discord.ChannelMention(e.Channel().ID())).SetEphemeral(true).
							Build())
						if err != nil {
							logrus.Error(err)
						}
					}
				} else {
					if isGlobal {
						err := e.CreateMessage(discord.NewMessageCreateBuilder().
							SetContent("Polls are now globally blocked.").SetEphemeral(true).
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
				}
			}
			/*case "list":
			list := listBlockPolls(e.GuildID().String())*/
		}
	},
}
