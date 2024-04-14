package main

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/vaporvee/acecore/struct_cmd"
)

var Plugin = &struct_cmd.Plugin{
	Name: "testplugin",
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
