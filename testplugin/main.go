package main

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/vaporvee/acecore/struct_cmd"
)

var Plugin = &struct_cmd.Plugin{
	Register: func(client bot.Client, commands *[]struct_cmd.Command) error {
		*commands = append(*commands, struct_cmd.Command{
			Definition: discord.SlashCommandCreate{
				Name:        "TESTPLUGIN",
				Description: "TESTPLUGIN",
			},
			Interact: func(e *events.ApplicationCommandInteractionCreate) {
				e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("TEST").SetEphemeral(true).Build())
			},
		})
		return nil
	},
}
