package main

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/sirupsen/logrus"
	"github.com/vaporvee/acecore/struct_cmd"
)

var Plugin = &struct_cmd.Plugin{
	Name: "testplugin",
	Register: func(e *events.Ready) error {
		app, err := e.Client().Rest().GetCurrentApplication()
		if err != nil {
			return err
		}
		logrus.Infof("%s has a working plugin called \"testplugin\"", app.Bot.Username)
		return nil
	},
	Commands: []struct_cmd.Command{
		{
			Definition: discord.SlashCommandCreate{
				Name:        "testplugincommand",
				Description: "Tesing if plugins work",
			},
			Interact: func(e *events.ApplicationCommandInteractionCreate) {
				e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Plugins are working!").SetEphemeral(true).Build())
			},
		},
	},
}
