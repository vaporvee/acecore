package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var sticky_command Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:                     "sticky",
		Description:              "Stick messages to the bottom of the current channel",
		DefaultMemberPermissions: int64Ptr(discordgo.PermissionManageMessages),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add",
				Description: "Stick messages to the bottom of the current channel",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "remove",
				Description: "Remove the sticky message of the current channel",
			},
		},
	},
	Interact: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.ApplicationCommandData().Options[0].Name {
		case "add":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseModal,
				Data: &discordgo.InteractionResponseData{
					CustomID: "sticky_modal",
					Title:    "Sticky message",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:    "sticky_modal_text",
									Label:       "Text",
									Style:       discordgo.TextInputParagraph,
									Placeholder: "The message you want to stick to the bottom of this channel",
									Required:    true,
									MaxLength:   2000,
									Value:       "",
								},
							},
						},
					},
				},
			})
		case "remove":
			if hasSticky(i.GuildID, i.ChannelID) {
				s.ChannelMessageDelete(i.ChannelID, getStickyMessageID(i.GuildID, i.ChannelID))
				removeSticky(i.GuildID, i.ChannelID)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "The sticky message was removed from this channel!",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "This channel has no sticky message!",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			}
		}
	},
	ModalID: "sticky_modal",
	ModalSubmit: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		text := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		message, err := s.ChannelMessageSendEmbed(i.ChannelID, &discordgo.MessageEmbed{
			Type:        discordgo.EmbedTypeArticle,
			Title:       ":pushpin: Sticky message",
			Color:       hexToDecimal(color["primary"]),
			Description: text,
		})
		if err != nil {
			log.Println(err)
		}
		if addSticky(i.GuildID, i.ChannelID, text, message.ID) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Sticky message in this channel was updated!",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		} else {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Message sticked to the channel!",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
	},
}