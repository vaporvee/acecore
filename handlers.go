package main

import (
	"slices"
	"strings"

	"github.com/disgoorg/disgo/events"

	"github.com/vaporvee/acecore/shared"
)

var commands []shared.Command

func applicationCommandInteractionCreate(e *events.ApplicationCommandInteractionCreate) {
	for _, command := range commands {
		if command.Interact != nil && e.Data.CommandName() == command.Definition.CommandName() {
			command.Interact(e)
		}
	}

}

func autocompleteInteractionCreate(e *events.AutocompleteInteractionCreate) {
	for _, command := range commands {
		if command.Autocomplete != nil && e.Data.CommandName == command.Definition.CommandName() {
			command.Autocomplete(e)
		}
	}
}

func componentInteractionCreate(e *events.ComponentInteractionCreate) {
	for _, command := range commands {
		if command.ComponentInteract != nil {
			if slices.Contains(command.ComponentIDs, e.Data.CustomID()) || slices.ContainsFunc(command.DynamicComponentIDs(), func(id string) bool {
				var customID string
				if strings.ContainsAny(e.Data.CustomID(), ";") {
					customID = strings.TrimSuffix(e.Data.CustomID(), ";"+strings.Split(e.Data.CustomID(), ";")[1])
				} else {
					customID = e.Data.CustomID()
				}
				return id == customID
			}) {
				command.ComponentInteract(e)
			}
		}
	}
}

func modalSubmitInteractionCreate(e *events.ModalSubmitInteractionCreate) {
	for _, command := range commands {
		if command.ModalSubmit != nil {
			var hasID bool = false
			var modalIDs []string
			if command.ModalIDs != nil {
				modalIDs = command.ModalIDs
			}
			if command.DynamicModalIDs != nil {
				modalIDs = append(command.ModalIDs, command.DynamicModalIDs()...)
			}
			for _, modalID := range modalIDs {
				if strings.HasPrefix(e.Data.CustomID, modalID) {
					hasID = true
					break
				}
			}
			if hasID {
				command.ModalSubmit(e)
				return // I have no idea why it crashes without that return
			}
		}
	}
}
