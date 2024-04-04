package main

/*
import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var cmd_sticky Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:                     "sticky",
		Description:              "Stick or unstick messages to the bottom of the current channel",
		DefaultMemberPermissions: int64Ptr(discordgo.PermissionManageMessages),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "message",
				Description: "The message you want to stick to the bottom of this channel",
				Required:    false,
			},
		},
	},
	Interact: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if len(i.ApplicationCommandData().Options) == 0 {
			if hasSticky(i.GuildID, i.ChannelID) {
				err := s.ChannelMessageDelete(i.ChannelID, getStickyMessageID(i.GuildID, i.ChannelID))
				if err != nil {
					logrus.Error(err)
				}
				removeSticky(i.GuildID, i.ChannelID)
				err = respond(i.Interaction, "The sticky message was removed from this channel!", true)
				if err != nil {
					logrus.Error(err)
				}
			} else {
				err := respond(i.Interaction, "This channel has no sticky message!", true)
				if err != nil {
					logrus.Error(err)
				}
			}
		} else {
			inputStickyMessage(i)
		}
	},
}

var context_sticky Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:                     "Stick to channel",
		Type:                     discordgo.MessageApplicationCommand,
		DefaultMemberPermissions: int64Ptr(discordgo.PermissionManageMessages),
	},
	Interact: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		inputStickyMessage(i)
	},
}

func inputStickyMessage(i *discordgo.InteractionCreate) {
	var messageText string
	if len(i.ApplicationCommandData().Options) == 0 {
		messageText = i.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID].Content //TODO add more data then just content
	} else {
		messageText = i.ApplicationCommandData().Options[0].StringValue()
	}
	if messageText == "" {
		err := respond(i.Interaction, "Can't add empty sticky messages!", true)
		if err != nil {
			logrus.Error(err)
		}
	} else {
		message, err := bot.ChannelMessageSendEmbed(i.ChannelID, &discordgo.MessageEmbed{
			Type: discordgo.EmbedTypeArticle,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "📌 Sticky message",
			},
			Color:       hexToDecimal(color["primary"]),
			Description: messageText,
		})
		if err != nil {
			logrus.Error(err)
		}

		if hasSticky(i.GuildID, i.ChannelID) {
			err = bot.ChannelMessageDelete(i.ChannelID, getStickyMessageID(i.GuildID, i.ChannelID))
			if err != nil {
				logrus.Error(err, getStickyMessageID(i.GuildID, i.ChannelID))
			}
			removeSticky(i.GuildID, i.ChannelID)
			addSticky(i.GuildID, i.ChannelID, messageText, message.ID)
			err = respond(i.Interaction, "Sticky message in this channel was updated!", true)
			if err != nil {
				logrus.Error(err)
			}
		} else {
			addSticky(i.GuildID, i.ChannelID, messageText, message.ID)
			err := respond(i.Interaction, "Message sticked to the channel!", true)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}
*/
