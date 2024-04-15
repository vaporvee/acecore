package main

import (
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/sirupsen/logrus"
)

func guildMemberJoin(e *events.GuildMemberJoin) {
	role := getAutoJoinRole(e.GuildID.String(), e.Member.User.Bot)
	if role != "" {
		err := e.Client().Rest().AddMemberRole(e.GuildID, e.Member.User.ID, snowflake.MustParse(role))
		if err != nil {
			logrus.Error(err)
		}
	}
}
