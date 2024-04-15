package main

import (
	"github.com/sirupsen/logrus"
)

func setAutoJoinRole(guildID string, option string, roleID string) bool {
	var role_exists bool
	var autojoinroles_exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT  1 FROM autojoinroles WHERE guild_id = $1)", guildID).Scan(&autojoinroles_exists)
	if err != nil {
		logrus.Error(err)
	}
	err = db.QueryRow("SELECT EXISTS (SELECT  1 FROM autojoinroles WHERE guild_id = $1 AND "+option+"_role IS NOT NULL AND "+option+"_role != '')", guildID).Scan(&role_exists)
	if err != nil {
		logrus.Error(err)
	}
	if autojoinroles_exists {
		_, err = db.Exec("UPDATE autojoinroles SET "+option+"_role = $1 WHERE guild_id = $2", roleID, guildID)
		if err != nil {
			logrus.Error(err)
		}
	} else {
		_, err = db.Exec("INSERT INTO autojoinroles (guild_id, "+option+"_role) VALUES ($1, $2)", guildID, roleID)
		if err != nil {
			logrus.Error(err)
		}
	}
	return role_exists
}

func purgeUnusedAutoJoinRoles(guildID string) {
	_, err := db.Exec("DELETE FROM autojoinroles WHERE guild_id = $1 AND user_role = '' OR user_role IS NULL AND bot_role = '' OR bot_role IS NULL", guildID)
	if err != nil {
		logrus.Error(err)
	}
}

func getAutoJoinRole(guildID string, isBot bool) string {
	var isBotString string
	var role string
	if isBot {
		isBotString = "bot"
	} else {
		isBotString = "user"
	}
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM autojoinroles WHERE guild_id = $1)", guildID).Scan(&exists)
	if err != nil {
		logrus.Error(err)
		return role
	}
	if exists {
		err = db.QueryRow("SELECT "+isBotString+"_role FROM autojoinroles WHERE guild_id = $1", guildID).Scan(&role)
		if err != nil {
			logrus.Error(err, guildID)
		}
	}
	return role
}
