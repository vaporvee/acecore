package main

import (
	"os"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type Command struct {
	Definition          discordgo.ApplicationCommand
	Interact            func(s *discordgo.Session, i *discordgo.InteractionCreate)
	ComponentInteract   func(s *discordgo.Session, i *discordgo.InteractionCreate)
	Autocomplete        func(s *discordgo.Session, i *discordgo.InteractionCreate)
	ModalSubmit         func(s *discordgo.Session, i *discordgo.InteractionCreate)
	ComponentIDs        []string
	ModalIDs            []string
	DynamicComponentIDs func() []string
	DynamicModalIDs     func() []string
}

var commands []Command = []Command{cmd_form, cmd_tag, cmd_tag_short, cmd_dadjoke, cmd_ping, cmd_ask, cmd_sticky, cmd_cat, cmd_autojoinroles, cmd_autopublish, context_sticky, context_tag}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	logrus.Info("Starting up...")
	findAndDeleteUnusedMessages()
	removeOldCommandFromAllGuilds(s)
	var existingCommandNames []string
	for _, guild := range event.Guilds {
		existingCommands, err := s.ApplicationCommands(s.State.User.ID, guild.ID)
		for _, existingCommand := range existingCommands {
			existingCommandNames = append(existingCommandNames, existingCommand.Name)
		}
		if err != nil {
			logrus.Errorf("error fetching existing commands for guild %s: %v", guild.Name, err)
			continue
		}
		for _, command := range commands {
			if !slices.Contains(existingCommandNames, command.Definition.Name) || slices.Contains(os.Args, "--update="+command.Definition.Name) {
				cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guild.ID, &command.Definition)
				logrus.Infof("Added command '%s'", cmd.Name)
				if err != nil {
					logrus.Error("error creating command,", err)
					continue
				}
			}
		}
	}
	logrus.Info("Successfully started the Bot!")
}

func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	for _, command := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, event.Guild.ID, &command.Definition)
		if err != nil {
			logrus.Errorf("error creating command for guild %s: %v", event.Guild.Name, err)
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
		case discordgo.InteractionModalSubmit:
			if command.ModalSubmit != nil {
				// FIXME: Makes it dynamic i don't know why it isn't otherwise
				if command.Definition.Name == "form" {
					command.ModalIDs = getFormButtonIDs()
				}
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
			logrus.Errorf("error fetching existing commands for guild %s: %v\n", guild.Name, err)
			continue
		}
		var commandIDs []string
		for _, command := range commands {
			commandIDs = append(commandIDs, command.Definition.Name)
		}
		for _, existingCommand := range existingCommands {
			if !slices.Contains(commandIDs, existingCommand.Name) {
				logrus.Infof("Deleting command '%s'", existingCommand.Name)
				err := s.ApplicationCommandDelete(s.State.User.ID, guild.ID, existingCommand.ID)
				if err != nil {
					logrus.Errorf("error deleting command %s for guild %s: %v", existingCommand.Name, guild.Name, err)
				}
			}
		}
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Embeds) == 0 || m.Embeds[0].Footer == nil || m.Embeds[0].Footer.Text != "📌 Sticky message" {
		if hasSticky(m.GuildID, m.ChannelID) {
			err := s.ChannelMessageDelete(m.ChannelID, getStickyMessageID(m.GuildID, m.ChannelID))
			stickyMessage, _ := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
				Type: discordgo.EmbedTypeArticle,
				Footer: &discordgo.MessageEmbedFooter{
					Text: "📌 Sticky message",
				},
				Color:       hexToDecimal(color["primary"]),
				Description: getStickyMessageContent(m.GuildID, m.ChannelID),
			})
			if err != nil {
				logrus.Error(err)
			}
			updateStickyMessageID(m.GuildID, m.ChannelID, stickyMessage.ID)
		}
	}
	channel, _ := s.Channel(m.ChannelID)
	if channel.Type == discordgo.ChannelTypeGuildNews {
		if isAutopublishEnabled(m.GuildID, m.ChannelID) {
			_, err := s.ChannelMessageCrosspost(m.ChannelID, m.ID)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}

func messageDelete(s *discordgo.Session, m *discordgo.MessageDelete) { //TODO: also clear on bot start when message doesn't exist
	tryDeleteUnusedMessage(m.ID)
}

func guildMemberJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	role := getAutoJoinRole(m.GuildID, m.User.Bot)
	if role != "" {
		err := s.GuildMemberRoleAdd(m.GuildID, m.User.ID, role)
		if err != nil {
			logrus.Error(err)
		}
	}
}
