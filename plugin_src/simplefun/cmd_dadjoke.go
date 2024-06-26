package main

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/sirupsen/logrus"
	"github.com/vaporvee/acecore/shared"
)

var cmd_dadjoke = shared.Command{
	Definition: discord.SlashCommandCreate{
		Name:        "dadjoke",
		Description: "Gives you a random joke that is as bad as your dad would tell them",
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeGuild,
			discord.InteractionContextTypePrivateChannel,
			discord.InteractionContextTypeBotDM,
		},
		IntegrationTypes: []discord.ApplicationIntegrationType{
			discord.ApplicationIntegrationTypeGuildInstall,
			discord.ApplicationIntegrationTypeUserInstall,
		},
	},
	Interact: func(e *events.ApplicationCommandInteractionCreate) {
		joke := simpleGetFromAPI("joke", "https://icanhazdadjoke.com/").(string)
		err := e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent(joke).
			Build())
		if err != nil {
			logrus.Error(err)
		}
	},
}
