package shared

import (
	"database/sql"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type Command struct {
	Ready               func(e *events.Ready)
	Definition          discord.ApplicationCommandCreate
	Interact            func(e *events.ApplicationCommandInteractionCreate)
	Autocomplete        func(e *events.AutocompleteInteractionCreate)
	ComponentInteract   func(e *events.ComponentInteractionCreate)
	ModalSubmit         func(e *events.ModalSubmitInteractionCreate)
	ComponentIDs        []string
	ModalIDs            []string
	DynamicModalIDs     func() []string
	DynamicComponentIDs func() []string
}

type Plugin struct {
	Name     string
	Init     func(d *sql.DB) error
	Commands []Command
}

var BotConfigs []bot.ConfigOpt
