package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/sirupsen/logrus"

	"github.com/vaporvee/acecore/shared"
)

var commands []shared.Command

func ready(e *events.Ready) {
	logrus.Info("Starting up...")
	removeOldCommandFromAllGuilds(e.Client())
	var existingCommandNames []string
	existingCommands, err := e.Client().Rest().GetGlobalCommands(e.Client().ApplicationID(), false)
	if err != nil {
		logrus.Errorf("error fetching existing global commands: %v", err)
	} else {
		for _, existingCommand := range existingCommands {
			existingCommandNames = append(existingCommandNames, existingCommand.Name())
		}
	}
	globalCommands := []discord.ApplicationCommandCreate{}
	for _, command := range commands {
		if !slices.Contains(existingCommandNames, command.Definition.CommandName()) || slices.Contains(os.Args, "--update-all") || slices.Contains(os.Args, "--clean") {
			globalCommands = append(globalCommands, command.Definition)
			logrus.Infof("Appending command \"%s\"", command.Definition.CommandName())
		}
	}
	if len(globalCommands) > 0 {
		logrus.Infof("Attempting to add global commands %s", fmt.Sprint(globalCommands))
		_, err = e.Client().Rest().SetGlobalCommands(e.Client().ApplicationID(), globalCommands)
		if err != nil {
			logrus.Errorf("error creating global commands '%s'", err)
		} else {
			logrus.Infof("Added global commands sucessfully!")
		}
	}
	logrus.Info("Successfully started the Bot!")
}

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

func removeOldCommandFromAllGuilds(c bot.Client) {
	app, err := c.Rest().GetCurrentApplication()
	if err != nil {
		logrus.Error(err)
	}
	globalCommands, err := c.Rest().GetGlobalCommands(app.Bot.ID, false)
	if err != nil {
		logrus.Errorf("error fetching existing global commands: %v", err)
		return
	}
	var commandNames []string
	for _, command := range commands {
		commandNames = append(commandNames, command.Definition.CommandName())
	}
	for _, existingCommand := range globalCommands {
		if !slices.Contains(commandNames, existingCommand.Name()) {
			logrus.Infof("Deleting command '%s'", existingCommand.Name())
			err := c.Rest().DeleteGlobalCommand(c.ApplicationID(), existingCommand.ID())
			if err != nil {
				logrus.Errorf("error deleting command %s: %v", existingCommand.Name(), err)
			}
		}
	}
}
