package main

import (
	"log"

	"github.com/sirupsen/logrus"
)

type BlockPoll struct {
	ChannelID   string
	Global      bool
	AllowedRole string
}

func isGlobalBlockPolls(guildID string) bool {
	var globalexists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM blockpolls WHERE guild_id = $1 AND global = true)", guildID).Scan(&globalexists)
	if err != nil {
		logrus.Error(err)
	}
	return globalexists
}

func toggleBlockPolls(guildID string, channelID string, global bool, allowedRole string) (e bool, isGlobal bool) {
	globalexists := isGlobalBlockPolls(guildID)
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM blockpolls WHERE guild_id = $1 AND channel_id = $2)", guildID, channelID).Scan(&exists)
	if err != nil {
		logrus.Error(err)
	}
	if globalexists {
		_, err := db.Exec("DELETE FROM blockpolls WHERE guild_id = $1 AND global = true", guildID)
		if err != nil {
			logrus.Error(err)
		}
		return true, true
	} else if global {
		_, err = db.Exec("DELETE FROM blockpolls WHERE guild_id = $1", guildID)
		if err != nil {
			logrus.Error(err)
		}
		_, err := db.Exec("INSERT INTO blockpolls (guild_id, global, channel_id, allowed_role) VALUES ($1, true, $2, $3)", guildID, channelID, allowedRole)
		if err != nil {
			logrus.Error(err)
		}
		return false, true
	} else if exists && !globalexists {
		_, err := db.Exec("DELETE FROM blockpolls WHERE guild_id = $1 AND channel_id = $2", guildID, channelID)
		if err != nil {
			logrus.Error(err)
		}
		return true, false
	} else if !globalexists {
		_, err := db.Exec("INSERT INTO blockpolls (guild_id, channel_id, allowed_role) VALUES ($1, $2, $3)", guildID, channelID, allowedRole)
		if err != nil {
			logrus.Error(err)
		}
		return false, false
	} else {
		return false, false
	}
}

func listBlockPolls(guildID string) []BlockPoll {
	var list []BlockPoll
	rows, err := db.Query("SELECT channel_id, global, allowed_role FROM blockpolls WHERE guild_id = $1", guildID)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var bp BlockPoll
		err := rows.Scan(&bp.ChannelID, &bp.Global, &bp.AllowedRole)
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, bp)
	}
	return list
}

func getBlockPollsEnabled(guildID string, channelID string) (isEnabled bool, allowedRole string) {
	var enabled bool
	var v_allowedRole string
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM blockpolls WHERE guild_id = $1 AND channel_id = $2)", guildID, channelID).Scan(&enabled)
	if err != nil {
		logrus.Error(err)
	}
	err = db.QueryRow("SELECT allowed_role FROM blockpolls WHERE guild_id = $1 AND channel_id = $2", guildID, channelID).Scan(&v_allowedRole)
	if err != nil && err.Error() != "sql: no rows in result set" {
		logrus.Error(err)
	}
	return enabled, v_allowedRole
}
