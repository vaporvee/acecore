package main

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/vaporvee/acecore/cmd"
)

var Plugin = &cmd.Plugin{
	Name: "Info",
	Commands: []cmd.Command{
		{
			Definition: discord.SlashCommandCreate{
				Name:        "info",
				Description: "Gives you information about a user or this bot.",
				Contexts: []discord.InteractionContextType{
					discord.InteractionContextTypeGuild,
					discord.InteractionContextTypePrivateChannel,
					discord.InteractionContextTypeBotDM,
				},
				IntegrationTypes: []discord.ApplicationIntegrationType{
					discord.ApplicationIntegrationTypeGuildInstall,
					discord.ApplicationIntegrationTypeUserInstall,
				},
				Options: []discord.ApplicationCommandOption{
					&discord.ApplicationCommandOptionSubCommand{
						Name:        "user",
						Description: "Gives you information about a user and its profile images.",
						Options: []discord.ApplicationCommandOption{
							&discord.ApplicationCommandOptionUser{
								Name:        "user",
								Description: "The user you need information about.",
								Required:    true,
							},
						},
					},
					&discord.ApplicationCommandOptionSubCommand{
						Name:        "bot-service",
						Description: "Gives you information about this bot's server service.",
					},
				},
			},
			Interact: func(e *events.ApplicationCommandInteractionCreate) {
				switch *e.SlashCommandInteractionData().SubCommandName {
				case "user":
					user(e)
				case "bot-service":

				}

			},
		},
	},
}
