package main

import (
	"github.com/bwmarrin/discordgo"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != s.State.User.ID {
		if hasSticky(m.GuildID, m.ChannelID) {
			s.ChannelMessageDelete(m.ChannelID, getStickyMessageID(m.GuildID, m.ChannelID))
			stickyMessage, _ := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
				Type:        discordgo.EmbedTypeArticle,
				Title:       ":pushpin: Sticky message",
				Color:       hexToDecimal(color["primary"]),
				Description: getStickyMessageContent(m.GuildID, m.ChannelID),
			})
			updateStickyMessageID(m.GuildID, m.ChannelID, stickyMessage.ID)
		}
	}
}
