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
	var commands []Command = []Command{tag_command, short_get_tag_command}
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

/*
func hasManageServerPermissions(s *discordgo.Session, userID string, guildID string) bool {
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		fmt.Printf("Error fetching guild member: %v\n", err)
		return false
	}

	guild, err := s.Guild(guildID)
	if err != nil {
		fmt.Printf("Error fetching guild: %v\n", err)
		return false
	}

	if guild.OwnerID == userID {
		return true
	}

	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			fmt.Printf("Error fetching role: %v\n", err)
			continue
		}
		if role.Permissions&discordgo.PermissionManageServer != 0 || role.Permissions&discordgo.PermissionAdministrator != 0 {
			return true
		}
	}

	return false
}
*/
