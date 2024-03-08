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
				respond(i.Interaction, "Autopublishing is now disabled on <#"+i.ChannelID+">", true)
			} else {
				respond(i.Interaction, "Autopublishing is now enabled on <#"+i.ChannelID+">", true)
			}
		} else {
			respond(i.Interaction, "This is not an announcement channel!", true)
		}
	},
}
