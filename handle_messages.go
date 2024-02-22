package main

import (
	"github.com/bwmarrin/discordgo"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Embeds) == 0 || m.Embeds[0].Footer == nil || m.Embeds[0].Footer.Text != "📌 Sticky message" {
		if hasSticky(m.GuildID, m.ChannelID) {
			s.ChannelMessageDelete(m.ChannelID, getStickyMessageID(m.GuildID, m.ChannelID))
			stickyMessage, _ := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
				Type: discordgo.EmbedTypeArticle,
				Footer: &discordgo.MessageEmbedFooter{
					Text: "📌 Sticky message",
				},
				Color:       hexToDecimal(color["primary"]),
				Description: getStickyMessageContent(m.GuildID, m.ChannelID),
			})
			updateStickyMessageID(m.GuildID, m.ChannelID, stickyMessage.ID)
		}
	}
}
