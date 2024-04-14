package cmd

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type Command struct {
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
	Commands []Command
	Register func(e *events.Ready) error
}
