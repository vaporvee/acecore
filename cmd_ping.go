package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/sirupsen/logrus"
	"github.com/vaporvee/acecore/custom"
)

var cmd_ping Command = Command{
	Definition: discord.SlashCommandCreate{
		Name:        "ping",
		Description: "Returns the ping of the bot",
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
		start := time.Now()

		client := http.Client{
			Timeout: 5 * time.Second,
		}

		resp, err := client.Get("https://discord.com/api/v9/gateway/bot")
		if err != nil {
			logrus.Error(err)
			return
		}
		defer resp.Body.Close()

		ping := time.Since(start)
		var pingColor string
		if ping.Milliseconds() < 200 {
			pingColor = "green"
		} else if ping.Milliseconds() < 400 {
			pingColor = "yellow"
		} else {
			pingColor = "red"
		}
		app, err := e.Client().Rest().GetCurrentApplication()
		if err != nil {
			logrus.Error(err)
		}
		err = e.CreateMessage(discord.NewMessageCreateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetTitle(app.Bot.Username + " ping").
				SetDescription(fmt.Sprintf("# %.2fms", ping.Seconds()*1000)).
				SetColor(custom.GetColor(pingColor)).Build()).SetEphemeral(true).Build())
		if err != nil {
			logrus.Error(err)
		}
	},
}
