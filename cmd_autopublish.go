package main

import "github.com/bwmarrin/discordgo"

var cmd_autopublish Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:        "autopublish",
		Description: "Toggle automatically publishing every post in a announcement channel",
	},
	Interact: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		channel, _ := s.State.Channel(i.ChannelID)
		if channel.Type == discordgo.ChannelTypeGuildNews {
			if toggleAutoPublish(i.GuildID, i.ChannelID) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Autopublishing is now disabled on <#" + i.ChannelID + ">",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Autopublishing is now enabled on <#" + i.ChannelID + ">",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			}
		} else {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This is not an announcement channel!",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
	},
}
