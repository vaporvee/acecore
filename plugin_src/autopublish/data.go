package main

import (
	"github.com/sirupsen/logrus"
)

func toggleAutoPublish(guildID string, newsChannelID string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM autopublish WHERE guild_id = $1 AND news_channel_id = $2)", guildID, newsChannelID).Scan(&exists)
	if err != nil {
		logrus.Error(err)
	}
	if exists {
		_, err := db.Exec("DELETE FROM autopublish WHERE guild_id = $1 AND news_channel_id = $2", guildID, newsChannelID)
		if err != nil {
			logrus.Error(err)
		}
	} else {
		_, err := db.Exec("INSERT INTO autopublish (guild_id, news_channel_id) VALUES ($1, $2)", guildID, newsChannelID)
		if err != nil {
			logrus.Error(err)
		}
	}
	return exists
}

func isAutopublishEnabled(guildID string, newsChannelID string) bool {
	var enabled bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM autopublish WHERE guild_id = $1 AND news_channel_id = $2)", guildID, newsChannelID).Scan(&enabled)
	if err != nil {
		logrus.Error(err)
	}
	return enabled
}
