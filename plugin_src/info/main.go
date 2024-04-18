package main

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/vaporvee/acecore/shared"
)

var Plugin = &shared.Plugin{
	Name: "Info",
	Commands: []shared.Command{
		{
			Definition: discord.SlashCommandCreate{
				Name:        "info",
				Description: "Gives you information about a user or this app.",
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
						Name:        "app-service",
						Description: "Gives you information about this app's server service.",
					},
				},
			},
			Interact: func(e *events.ApplicationCommandInteractionCreate) {
				switch *e.SlashCommandInteractionData().SubCommandName {
				case "user":
					user(e)
				case "app-service":

				}

			},
		},
	},
}
