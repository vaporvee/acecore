package main

import (
	"slices"

	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/sirupsen/logrus"
)

func messageCreate(e *events.MessageCreate) {
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
	}
}
