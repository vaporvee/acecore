package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Definition   discordgo.ApplicationCommand
	Interact     func(s *discordgo.Session, i *discordgo.InteractionCreate)
	Autocomplete func(s *discordgo.Session, i *discordgo.InteractionCreate)
	ModalSubmit  func(s *discordgo.Session, i *discordgo.InteractionCreate)
	ModalID      string
}

var commands []Command = []Command{tag_command, short_get_tag_command, dadjoke_command, ping_command, ask_command}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	for _, guild := range event.Guilds {
		for _, command := range commands {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, guild.ID, &command.Definition)
			if err != nil {
				fmt.Println("error creating command,", err)
				continue
			}
		}
	}
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
		case discordgo.InteractionModalSubmit: //g has no modal so it crashes
			if command.ModalSubmit != nil && strings.HasPrefix(i.ModalSubmitData().CustomID, command.ModalID) {
				command.ModalSubmit(s, i)
			}
		}
	}
}

func removeCommandFromAllGuilds(s *discordgo.Session) {
	for _, guild := range s.State.Guilds {
		existingCommands, err := s.ApplicationCommands(s.State.User.ID, guild.ID)
		if err != nil {
			fmt.Printf("error fetching existing commands for guild %s: %v\n", guild.Name, err)
			continue
		}

		for _, existingCommand := range existingCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, guild.ID, existingCommand.ID)
			if err != nil {
				fmt.Printf("error deleting command %s for guild %s: %v\n", existingCommand.Name, guild.Name, err)
			}
		}
	}
}
