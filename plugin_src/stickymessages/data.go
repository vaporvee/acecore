package main

import (
	"github.com/sirupsen/logrus"
)

func addSticky(guildID string, channelID string, messageContent string, messageID string) {
	_, err := db.Exec("INSERT INTO sticky (guild_id, channel_id, message_id, message_content) VALUES ($1, $2, $3, $4)", guildID, channelID, messageID, messageContent)
	if err != nil {
		logrus.Error(err)
	}

}

func hasSticky(guildID string, channelID string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT  1 FROM sticky WHERE guild_id = $1 AND channel_id = $2)", guildID, channelID).Scan(&exists)
	if err != nil {
		logrus.Error(err)
	}
	return exists
}

func getStickyMessageID(guildID string, channelID string) string {
	var messageID string
	exists := hasSticky(guildID, channelID)
	if exists {
		err := db.QueryRow("SELECT message_id FROM sticky WHERE guild_id = $1 AND channel_id = $2", guildID, channelID).Scan(&messageID)
		if err != nil {
			logrus.Error(err)
		}
	}
	return messageID
}
func getStickyMessageContent(guildID string, channelID string) string {
	var messageID string
	exists := hasSticky(guildID, channelID)
	if exists {
		err := db.QueryRow("SELECT message_content FROM sticky WHERE guild_id = $1 AND channel_id = $2", guildID, channelID).Scan(&messageID)
		if err != nil {
			logrus.Error(err)
		}
	}
	return messageID
}

func updateStickyMessageID(guildID string, channelID string, messageID string) {
	exists := hasSticky(guildID, channelID)
	if exists {
		_, err := db.Exec("UPDATE sticky SET message_id = $1 WHERE guild_id = $2 AND channel_id = $3", messageID, guildID, channelID)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func removeSticky(guildID string, channelID string) {
	_, err := db.Exec("DELETE FROM sticky WHERE guild_id = $1 AND channel_id = $2", guildID, channelID)
	if err != nil {
		logrus.Error(err)
	}
}
