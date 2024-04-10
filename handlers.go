package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/sirupsen/logrus"
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
	AllowDM             bool
}

var commands []Command = []Command{cmd_tag, cmd_tag_short, context_tag, cmd_sticky, context_sticky, cmd_ping, cmd_userinfo, cmd_form, cmd_ask, cmd_cat, cmd_dadjoke, cmd_ticket_form, cmd_autopublish, cmd_autojoinroles}

func ready(e *events.Ready) {
	logrus.Info("Starting up...")
	findAndDeleteUnusedMessages(e.Client())
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
			if !command.AllowDM && e.ApplicationCommandInteraction.GuildID().String() == "" {
				err := e.CreateMessage(discord.NewMessageCreateBuilder().
					SetContent("This command is not available in DMs.").SetEphemeral(true).
					Build())
				if err != nil {
					logrus.Error(err)
				}
			} else {
				command.Interact(e)
			}

		}

	}
}

func autocompleteInteractionCreate(e *events.AutocompleteInteractionCreate) {
	for _, command := range commands {
		if command.Autocomplete != nil && e.Data.CommandName == command.Definition.CommandName() {
			if !command.AllowDM && e.AutocompleteInteraction.GuildID().String() == "" {
				err := e.AutocompleteResult(nil)
				if err != nil {
					logrus.Error(err)
				}
			} else {
				command.Autocomplete(e)
			}
		}
	}
}

func componentInteractionCreate(e *events.ComponentInteractionCreate) {
	for _, command := range commands {
		if !command.AllowDM && e.ComponentInteraction.GuildID().String() == "" {
			e.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("This component is not available in DMs.").SetEphemeral(true).
				Build())
		} else {
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
}

func modalSubmitInteractionCreate(e *events.ModalSubmitInteractionCreate) {
	for _, command := range commands {
		if !command.AllowDM && e.ModalSubmitInteraction.GuildID().String() == "" {
			e.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("This modal is not available in DMs.").SetEphemeral(true).
				Build())
		} else {
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

func messageCreate(e *events.MessageCreate) {
	if len(e.Message.Embeds) == 0 || e.Message.Embeds[0].Footer == nil || e.Message.Embeds[0].Footer.Text != "📌 Sticky message" {
		if hasSticky(e.Message.GuildID.String(), e.Message.ChannelID.String()) {
			stickymessageID := getStickyMessageID(e.Message.GuildID.String(), e.Message.ChannelID.String())
			err := e.Client().Rest().DeleteMessage(e.ChannelID, snowflake.MustParse(stickymessageID))
			stickyMessage, _ := e.Client().Rest().CreateMessage(e.ChannelID, discord.MessageCreate{
				Embeds: []discord.Embed{
					{
						Footer: &discord.EmbedFooter{
							Text: "📌 Sticky message",
						},
						Color:       hexToDecimal(color["primary"]),
						Description: getStickyMessageContent(e.Message.GuildID.String(), e.Message.ChannelID.String()),
					},
				},
			})
			if err != nil {
				logrus.Error(err)
			}
			updateStickyMessageID(e.Message.GuildID.String(), e.Message.ChannelID.String(), stickyMessage.ID.String())
		}
	}
	channel, err := e.Client().Rest().GetChannel(e.Message.ChannelID)
	if err != nil {
		logrus.Error(err)
	}
	if channel != nil && channel.Type() == discord.ChannelTypeGuildNews {
		logrus.Debug("HERE")
		if isAutopublishEnabled(e.GuildID.String(), e.ChannelID.String()) {
			_, err := e.Client().Rest().CrosspostMessage(e.ChannelID, e.MessageID)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}

func messageDelete(e *events.MessageDelete) { //TODO: also clear on bot start when message doesn't exist
	tryDeleteUnusedMessage(e.MessageID.String())
}

func guildMemberJoin(e *events.GuildMemberJoin) {
	logrus.Debug("TESSST")
	role := getAutoJoinRole(e.GuildID.String(), e.Member.User.Bot)
	if role != "" {
		err := e.Client().Rest().AddMemberRole(e.GuildID, e.Member.User.ID, snowflake.MustParse(role))
		if err != nil {
			logrus.Error(err)
		}
	}
}
