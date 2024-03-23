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
	AllowDM             bool
}

var commands []Command = []Command{cmd_form, cmd_tag, cmd_tag_short, cmd_dadjoke, cmd_ping, cmd_ask, cmd_sticky, cmd_cat, cmd_autojoinroles, cmd_autopublish, context_sticky, context_tag}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	logrus.Info("Starting up...")
	findAndDeleteUnusedMessages()
	removeOldCommandFromAllGuilds(s)
	var existingCommandNames []string
	existingCommands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		logrus.Errorf("error fetching existing global commands: %v", err)
	} else {
		for _, existingCommand := range existingCommands {
			existingCommandNames = append(existingCommandNames, existingCommand.Name)
		}
	}
	if slices.Contains(os.Args, "--clean") {
		guilds := s.State.Guilds
		if err != nil {
			logrus.Errorf("error retrieving guilds: %v", err)
		}

		for _, guild := range guilds {
			_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, guild.ID, []*discordgo.ApplicationCommand{})
			if err != nil {
				logrus.Errorf("error deleting guild commands: %v", err)
			}
		}
	}
	for _, command := range commands {
		if !slices.Contains(existingCommandNames, command.Definition.Name) || slices.Contains(os.Args, "--update="+command.Definition.Name) || slices.Contains(os.Args, "--update=all") || slices.Contains(os.Args, "--clean") {
			cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", &command.Definition)
			if err != nil {
				logrus.Errorf("error creating global command '%s': %v", cmd.Name, err)
			} else {
				logrus.Infof("Added global command '%s'", cmd.Name)
			}
		}
	}
	logrus.Info("Successfully started the Bot!")
}

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	for _, command := range commands {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if command.Interact != nil && i.ApplicationCommandData().Name == command.Definition.Name {
				if !command.AllowDM && i.Interaction.GuildID == "" {
					respond(i.Interaction, "This command is not available in DMs.", true)
				} else {
					command.Interact(s, i)
				}
			}
		case discordgo.InteractionApplicationCommandAutocomplete:
			if command.Autocomplete != nil && i.ApplicationCommandData().Name == command.Definition.Name {
				if !command.AllowDM && i.Interaction.GuildID == "" {
					err := bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionApplicationCommandAutocompleteResult,
						Data: &discordgo.InteractionResponseData{
							Choices: nil,
						},
					})
					if err != nil {
						logrus.Error(err)
					}
				} else {
					command.Autocomplete(s, i)
				}
			}
		case discordgo.InteractionModalSubmit:
			if !command.AllowDM && i.Interaction.GuildID == "" {
				respond(i.Interaction, "This modal is not available in DMs.", true)
			} else {
				if command.ModalSubmit != nil {
					// FIXME: Makes it dynamic i don't know why it isn't otherwise
					if command.Definition.Name == "form" {
						command.ModalIDs = getFormButtonIDs()
					}
					var hasID bool = false
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
			if !command.AllowDM && i.Interaction.GuildID == "" {
				respond(i.Interaction, "This component is not available in DMs.", true)
			} else {
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
}

func removeOldCommandFromAllGuilds(s *discordgo.Session) {
	existingCommands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		logrus.Errorf("error fetching existing commands: %v\n", err)
		var commandIDs []string
		for _, command := range commands {
			commandIDs = append(commandIDs, command.Definition.Name)
		}
		for _, existingCommand := range existingCommands {
			if !slices.Contains(commandIDs, existingCommand.Name) {
				logrus.Infof("Deleting command '%s'", existingCommand.Name)
				err := s.ApplicationCommandDelete(s.State.User.ID, "", existingCommand.ID)
				if err != nil {
					logrus.Errorf("error deleting command %s: %v", existingCommand.Name, err)
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
