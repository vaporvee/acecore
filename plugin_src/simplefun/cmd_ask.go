package main

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/sirupsen/logrus"
	"github.com/vaporvee/acecore/custom"
	"github.com/vaporvee/acecore/shared"
)

var cmd_ask = shared.Command{
	Definition: discord.SlashCommandCreate{
		Name:        "ask",
		Description: "Ask anything and get a gif as response!",
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
			&discord.ApplicationCommandOptionString{
				Name:        "question",
				Description: "The question you want to ask",
				Required:    true,
			},
		},
	},
	Interact: func(e *events.ApplicationCommandInteractionCreate) {
		err := e.CreateMessage(discord.NewMessageCreateBuilder().
			AddEmbeds(discord.NewEmbedBuilder().SetImage(simpleGetFromAPI("image", "https://yesno.wtf/api").(string)).SetColor(custom.GetColor("primary")).Build()).
			Build())
		if err != nil {
			logrus.Error(err)
		}
	},
}
