package main

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/sirupsen/logrus"
)

func messageCreate(e *events.MessageCreate) {
	channel, err := e.Client().Rest().GetChannel(e.Message.ChannelID)
	if err != nil {
		logrus.Error(err)
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
