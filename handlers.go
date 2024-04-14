package main

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"runtime"
	"slices"
	"strings"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/sirupsen/logrus"
	"github.com/vaporvee/acecore/custom"
	"github.com/vaporvee/acecore/struct_cmd"
)

var commands []struct_cmd.Command = []struct_cmd.Command{cmd_tag, cmd_tag_short, context_tag, cmd_sticky, context_sticky, cmd_ping, cmd_userinfo, cmd_addemoji, cmd_form, cmd_ask, cmd_cat, cmd_dadjoke, cmd_ticket_form, cmd_blockpolls, cmd_autopublish, cmd_autojoinroles}

func ready(e *events.Ready) {
	logrus.Info("Starting up...")
	findAndDeleteUnusedMessages(e.Client())
	removeOldCommandFromAllGuilds(e.Client())
	err := loadPlugins("plugins/", e)
	if err != nil {
		logrus.Error(err)
	}
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

func loadPlugins(directory string, e *events.Ready) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	// Determine the appropriate file extension for dynamic libraries
	var ext string
	switch runtime.GOOS {
	case "windows":
		ext = ".dll"
	case "linux":
		ext = ".so"
	case "darwin":
		ext = ".dylib"
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ext {
			p, err := plugin.Open(filepath.Join(directory, file.Name()))
			if err != nil {
				return err
			}

			symPlugin, err := p.Lookup("Plugin")
			if err != nil {
				logrus.Errorf("Error looking up symbol 'Plugin' in %s: %v", file.Name(), err)
				continue
			}

			pluginPtr, ok := symPlugin.(**struct_cmd.Plugin)
			if !ok {
				logrus.Errorf("Plugin does not match expected type")
				continue
			}

			plugin := *pluginPtr
			if plugin.Name == "" {
				logrus.Warn("Plugin is unnamed")
			}
			if plugin.Commands != nil {
				commands = append(commands, plugin.Commands...)
			} else {
				logrus.Errorf("Plugin %s has no commands set", plugin.Name)
				continue
			}
			if plugin.Register != nil {
				err = plugin.Register(e)
				if err == nil {
					logrus.Infof("Successfully appended plugin %s for registration", plugin.Name)
				} else {
					logrus.Errorf("Error registering plugin %s commands: %v", plugin.Name, err)
					continue
				}
			}

		}
	}

	return nil
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
						Color:       custom.GetColor("primary"),
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
	if channel != nil {
		isBlockPollsEnabledGlobal := isGlobalBlockPolls(e.GuildID.String())
		isBlockPollsEnabled, allowedRole := getBlockPollsEnabled(e.GuildID.String(), e.Message.ChannelID.String())
		var hasAllowedRole bool
		if allowedRole != "" {
			hasAllowedRole = slices.Contains(e.Message.Member.RoleIDs, snowflake.MustParse(allowedRole))
		}
		if (isBlockPollsEnabledGlobal || isBlockPollsEnabled) && !hasAllowedRole && messageIsPoll(e.Message.ChannelID.String(), e.Message.ID.String(), e.Client()) {
			e.Client().Rest().DeleteMessage(e.Message.ChannelID, e.Message.ID)
		}
		if channel.Type() == discord.ChannelTypeGuildNews {
			if isAutopublishEnabled(e.GuildID.String(), e.ChannelID.String()) {
				_, err := e.Client().Rest().CrosspostMessage(e.ChannelID, e.MessageID)
				if err != nil {
					logrus.Error(err)
					return
				}
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
