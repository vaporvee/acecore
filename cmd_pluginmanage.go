package main

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/json"
	"github.com/sirupsen/logrus"
	"github.com/vaporvee/acecore/custom"
	"github.com/vaporvee/acecore/shared"
)

var cmd_pluginmanage shared.Command = shared.Command{
	Definition: discord.SlashCommandCreate{
		Name:                     "plugin",
		Description:              "Manage the plugins for this bot.",
		DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionAdministrator),
		Options: []discord.ApplicationCommandOption{
			&discord.ApplicationCommandOptionSubCommand{
				Name:        "list",
				Description: "List all installed plugins for this bot.",
			},
		},
	},
	Interact: func(e *events.ApplicationCommandInteractionCreate) {
		app, err := e.Client().Rest().GetCurrentApplication()
		if err != nil {
			logrus.Error(err)
			return
		}
		if app.Owner.ID == e.User().ID {
			switch *e.SlashCommandInteractionData().SubCommandName {
			case "list":
				var fields []discord.EmbedField
				for _, name := range pluginNames {
					fields = append(fields, discord.EmbedField{Name: name})
				}
				e.CreateMessage(discord.NewMessageCreateBuilder().
					SetEmbeds(discord.NewEmbedBuilder().
						SetTitle("Plugins").SetDescription("These are the currently installed plugins for this bot.").SetFields(fields...).SetColor(custom.GetColor("primary")).
						Build()).SetEphemeral(true).Build())
			}
		} else {
			e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("You are not the owner of this bot.").SetEphemeral(true).Build())
		}
	},
}
