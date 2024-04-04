package main

import (
	"slices"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/sirupsen/logrus"
)

type Command struct {
	Definition          discord.SlashCommandCreate
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

var commands []Command = []Command{cmd_tag, cmd_tag_short, context_tag /*, cmd_form, cmd_ticket_form, cmd_dadjoke, cmd_ping, cmd_ask, cmd_sticky, cmd_cat, cmd_autojoinroles, cmd_autopublish, context_sticky, cmd_userinfo*/}

func ready(e *events.Ready) {
	logrus.Info("Starting up...")
	//findAndDeleteUnusedMessages()
	removeOldCommandFromAllGuilds()
	/*
		var existingCommandNames []string
		existingCommands, err := client.Rest().GetGlobalCommands(client.ApplicationID(), false)
		if err != nil {
			logrus.Errorf("error fetching existing global commands: %v", err)
		} else {
			for _, existingCommand := range existingCommands {
				existingCommandNames = append(existingCommandNames, existingCommand.Name())
			}
		}
		for _, command := range commands {
			if !slices.Contains(existingCommandNames, command.Definition.Name) || slices.Contains(os.Args, "--update="+command.Definition.Name) || slices.Contains(os.Args, "--update=all") || slices.Contains(os.Args, "--clean") {
				cmd, err := client.Rest().CreateGlobalCommand(client.ApplicationID(), command.Definition)
				if err != nil {
					logrus.Errorf("error creating global command '%s': %v", cmd.Name(), err)
				} else {
					logrus.Infof("Added global command '%s'", cmd.Name())
				}
			}
		}
		logrus.Info("Successfully started the Bot!")
	*/
}

func applicationCommandInteractionCreate(e *events.ApplicationCommandInteractionCreate) {
	for _, command := range commands {
		if command.Interact != nil && e.SlashCommandInteractionData().CommandName() == command.Definition.Name {
			if !command.AllowDM && e.SlashCommandInteractionData().GuildID().String() == "" {
				e.CreateMessage(discord.NewMessageCreateBuilder().
					SetContent("This command is not available in DMs.").SetEphemeral(true).
					Build())
			} else {
				command.Interact(e)
			}
		}
	}
}

func autocompleteInteractionCreate(e *events.AutocompleteInteractionCreate) {
	for _, command := range commands {
		if command.Autocomplete != nil && e.Data.CommandName == command.Definition.Name {
			if !command.AllowDM && e.GuildID().String() == "" {
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
		if !command.AllowDM && e.GuildID().String() == "" {
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
		if !command.AllowDM && e.GuildID().String() == "" {
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

func removeOldCommandFromAllGuilds() {
	logrus.Debug(app.Bot.ID.String())
	globalCommands, err := client.Rest().GetGlobalCommands(app.Bot.ID, false)
	logrus.Debug("HERE") //doesnt get called
	if err != nil {
		logrus.Errorf("error fetching existing global commands: %v", err)
		return
	}
	var commandNames []string
	for _, command := range commands {
		commandNames = append(commandNames, command.Definition.Name)
	}
	for _, existingCommand := range globalCommands {
		if slices.Contains(commandNames, existingCommand.Name()) {
			logrus.Infof("Deleting command '%s'", existingCommand.Name())
			err := client.Rest().DeleteGlobalCommand(client.ApplicationID(), existingCommand.ID())
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
	channel, _ := e.Channel()
	if channel.Type() == discord.ChannelTypeGuildNews {
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
	role := getAutoJoinRole(e.GuildID.String(), e.Member.User.Bot)
	if role != "" {
		err := e.Client().Rest().AddMemberRole(e.GuildID, e.Member.User.ID, snowflake.MustParse(role))
		if err != nil {
			logrus.Error(err)
		}
	}
}
