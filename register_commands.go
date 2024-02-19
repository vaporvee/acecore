package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Definition discordgo.ApplicationCommand
	Interact   func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	commands := []*discordgo.ApplicationCommand{
		&tag_command.Definition,
		&short_get_tag_command.Definition,
	}

	for _, guild := range event.Guilds {
		for _, command := range commands {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, guild.ID, command)
			if err != nil {
				fmt.Println("error creating command,", err)
				continue // Continue to the next guild
			}
		}
	}
}

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "tag":
		tag_command.Interact(s, i)
	case "g":
		short_get_tag_command.Interact(s, i)
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
