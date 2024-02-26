package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Definition        discordgo.ApplicationCommand
	Interact          func(s *discordgo.Session, i *discordgo.InteractionCreate)
	ComponentInteract func(s *discordgo.Session, i *discordgo.InteractionCreate)
	ComponentIDs      []string
	Autocomplete      func(s *discordgo.Session, i *discordgo.InteractionCreate)
	ModalSubmit       func(s *discordgo.Session, i *discordgo.InteractionCreate)
	ModalIDs          []string
}

var commands []Command = []Command{form_command, tag_command, short_get_tag_command, dadjoke_command, ping_command, ask_command, sticky_command, cat_command}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Print("\nStarting up... (May take longer when Discord rate limits the bot)")
	removeOldCommandFromAllGuilds(s)
	var existingCommandNames []string
	for _, guild := range event.Guilds {
		existingCommands, err := s.ApplicationCommands(s.State.User.ID, guild.ID)
		for _, existingCommand := range existingCommands {
			existingCommandNames = append(existingCommandNames, existingCommand.Name)
		}
		if err != nil {
			fmt.Printf("error fetching existing commands for guild %s: %v\n", guild.Name, err)
			continue
		}
		for _, command := range commands {
			if !slices.Contains(existingCommandNames, command.Definition.Name) || slices.Contains(os.Args, "--update="+command.Definition.Name) {
				cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guild.ID, &command.Definition)
				fmt.Printf("\nAdded command \"%s\"", cmd.Name)
				if err != nil {
					fmt.Println("error creating command,", err)
					continue
				}
			}
		}
	}
	fmt.Print("\nSuccessfully started the Bot!")
}

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	for _, command := range commands {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if command.Interact != nil && i.ApplicationCommandData().Name == command.Definition.Name {
				command.Interact(s, i)
			}
		case discordgo.InteractionApplicationCommandAutocomplete:
			if command.Autocomplete != nil && i.ApplicationCommandData().Name == command.Definition.Name {
				command.Autocomplete(s, i)
			}
		case discordgo.InteractionModalSubmit:
			var hasID bool = false
			if command.ModalSubmit != nil {
				for _, modalID := range command.ModalIDs {
					if strings.HasPrefix(i.ModalSubmitData().CustomID, modalID) {
						hasID = true
					}
				}
				if hasID {
					command.ModalSubmit(s, i)
				}
			}
		case discordgo.InteractionMessageComponent:
			if command.ComponentInteract != nil {
				if command.Definition.Name == "form" {
					command.ComponentIDs = getFormButtonIDs()
				} // FIXME: Makes it dynamic i don't know why it isn't otherwise
				if slices.Contains(command.ComponentIDs, i.MessageComponentData().CustomID) {
					command.ComponentInteract(s, i)
				}
			}
		}
	}
}

func removeOldCommandFromAllGuilds(s *discordgo.Session) {
	for _, guild := range s.State.Guilds {
		existingCommands, err := s.ApplicationCommands(s.State.User.ID, guild.ID)
		if err != nil {
			fmt.Printf("error fetching existing commands for guild %s: %v\n", guild.Name, err)
			continue
		}
		var commandIDs []string
		for _, command := range commands {
			commandIDs = append(commandIDs, command.Definition.Name)
		}
		for _, existingCommand := range existingCommands {
			if !slices.Contains(commandIDs, existingCommand.Name) {
				fmt.Printf("\nDeleting command \"%s\"", existingCommand.Name)
				err := s.ApplicationCommandDelete(s.State.User.ID, guild.ID, existingCommand.ID)
				if err != nil {
					fmt.Printf("error deleting command %s for guild %s: %v\n", existingCommand.Name, guild.Name, err)
				}
			}
		}
	}
}
