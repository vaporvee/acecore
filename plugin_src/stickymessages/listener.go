package main

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/sirupsen/logrus"
	"github.com/vaporvee/acecore/custom"
)

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
}
