package main

import (
	"database/sql"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/vaporvee/acecore/shared"
)

var db *sql.DB

var Plugin = &shared.Plugin{
	Name: "testplugin",
	Init: func(d *sql.DB) error {
		db = d
		return nil
	},
	Commands: []shared.Command{
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
