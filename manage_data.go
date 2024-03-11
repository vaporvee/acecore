package main

import (
	"log"

	"github.com/google/uuid"
)

func initTables() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS tags (
		tag_id TEXT NOT NULL,
		tag_name TEXT NOT NULL,
		tag_content TEXT NOT NULL,
		guild_id TEXT NOT NULL,
		PRIMARY KEY (tag_id, guild_id)
	);
	CREATE TABLE IF NOT EXISTS sticky (
		message_id TEXT NOT NULL,
		channel_id TEXT NOT NULL,
		message_content TEXT NOT NULL,
		guild_id TEXT NOT NULL,
		PRIMARY KEY (channel_id, guild_id)
	);
	CREATE TABLE IF NOT EXISTS custom_forms (
		form_type TEXT NOT NULL,
		title TEXT NOT NULL,
		json JSON NOT NULL,
		guild_id TEXT NOT NULL,
		PRIMARY KEY (form_type, guild_id)
	);
	CREATE TABLE IF NOT EXISTS form_manage (
		form_manage_id TEXT NOT NULL,
		form_type TEXT NOT NULL,
		overwrite_title TEXT,
		message_id TEXT NOT NULL,
		channel_id TEXT NOT NULL,
		guild_id TEXT NOT NULL,
		result_channel_id TEXT NOT NULL,
		accept_channel_id TEXT,
		mods_can_comment BOOL,
		PRIMARY KEY (form_manage_id, form_type)
	);
	CREATE TABLE IF NOT EXISTS autojoinroles (
		guild_id TEXT NOT NULL,
		bot_role TEXT,
		user_role TEXT,
		PRIMARY KEY (guild_id)
	);
	CREATE TABLE IF NOT EXISTS autopublish (
		guild_id TEXT NOT NULL,
		news_channel_id TEXT NOT NULL,
		PRIMARY KEY (guild_id, news_channel_id)
	)
	`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

type FormResult struct {
	OverwriteTitle  string
	ResultChannelID string
	AcceptChannelID string
	ModsCanComment  bool
}

func addTag(guildID, tagName, tagContent string) bool {
	var exists bool = true
	//TODO: add modify command
	id := uuid.New()
	for exists {
		id = uuid.New()
		err := db.QueryRow("SELECT EXISTS (SELECT  1 FROM tags WHERE guild_id = $1 AND tag_id = $2)", guildID, id).Scan(&exists)
		if err != nil {
			log.Println(err)
		}
	}
	_, err := db.Exec("INSERT INTO tags (guild_id, tag_name, tag_content, tag_id) VALUES ($1, $2, $3, $4)", guildID, tagName, tagContent, id)
	if err != nil {
		log.Println(err)
	}

	return exists
}
func removeTag(guildID string, tagID string) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT  1 FROM tags WHERE guild_id = $1 AND tag_id = $2)", guildID, tagID).Scan(&exists)
	if err != nil {
		log.Println(err)
	}
	if exists {
		_, err = db.Exec("DELETE FROM tags WHERE guild_id = $1 AND tag_id = $2", guildID, tagID)
		if err != nil {
			log.Println(err)
		}
	}
}
func getTagIDs(guildID string) ([]string, error) {
	var IDs []string
	rows, err := db.Query("SELECT tag_id FROM tags WHERE guild_id = $1", guildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		IDs = append(IDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return IDs, nil
}
func getTagName(guildID string, tagID string) string {
	var tagName string
	db.QueryRow("SELECT tag_name FROM tags WHERE guild_id = $1 AND tag_id = $2", guildID, tagID).Scan(&tagName)
	return tagName
}
func getTagContent(guildID string, tagID string) string {
	var tagContent string
	db.QueryRow("SELECT tag_content FROM tags WHERE guild_id = $1 AND tag_id = $2", guildID, tagID).Scan(&tagContent)
	return tagContent
}

func addSticky(guildID string, channelID string, messageContent string, messageID string) bool {
	exists := hasSticky(guildID, channelID)
	if exists {
		_, err := db.Exec("UPDATE sticky SET message_content = $1 WHERE guild_id = $2 AND channel_id = $3", messageContent, guildID, channelID)
		if err != nil {
			log.Println(err)
		}
	} else {
		_, err := db.Exec("INSERT INTO sticky (guild_id, channel_id, message_id, message_content) VALUES ($1, $2, $3, $4)", guildID, channelID, messageID, messageContent)
		if err != nil {
			log.Println(err)
		}
	}
	return exists
}

func hasSticky(guildID string, channelID string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT  1 FROM sticky WHERE guild_id = $1 AND channel_id = $2)", guildID, channelID).Scan(&exists)
	if err != nil {
		log.Println(err)
	}
	return exists
}

func getStickyMessageID(guildID string, channelID string) string {
	var messageID string
	exists := hasSticky(guildID, channelID)
	if exists {
		err := db.QueryRow("SELECT message_id FROM sticky WHERE guild_id = $1 AND channel_id = $2", guildID, channelID).Scan(&messageID)
		if err != nil {
			log.Println(err)
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
			log.Println(err)
		}
	}
	return messageID
}

func updateStickyMessageID(guildID string, channelID string, messageID string) {
	exists := hasSticky(guildID, channelID)
	if exists {
		_, err := db.Exec("UPDATE sticky SET message_id = $1 WHERE guild_id = $2 AND channel_id = $3", messageID, guildID, channelID)
		if err != nil {
			log.Println(err)
		}
	}
}

func removeSticky(guildID string, channelID string) {
	_, err := db.Exec("DELETE FROM sticky WHERE guild_id = $1 AND channel_id = $2", guildID, channelID)
	if err != nil {
		log.Println(err)
	}
}

func getFormManageIdExists(id uuid.UUID) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM form_manage WHERE form_manage_id = $1)", id).Scan(&exists)
	if err != nil {
		log.Println(err)
	}
	return exists
}

func addFormButton(guildID string, channelID string, messageID string, formManageID string, formType string, resultChannelID string, overwriteTitle string, acceptChannelID string, modsCanComment bool) {
	_, err := db.Exec("INSERT INTO form_manage (guild_id, form_manage_id, channel_id, message_id, form_type, result_channel_id, overwrite_title, accept_channel_id, mods_can_comment) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", guildID, formManageID, channelID, messageID, formType, resultChannelID, overwriteTitle, acceptChannelID, modsCanComment)
	if err != nil {
		log.Println(err)
	}
}

func getFormManageIDs() []string {
	if db == nil {
		return nil
	}
	var IDs []string
	rows, err := db.Query("SELECT form_manage_id FROM form_manage")
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Println(err)
			return nil
		}
		IDs = append(IDs, id)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil
	}
	return IDs
}

func getFormType(formManageID string) string {
	var formType string
	err := db.QueryRow("SELECT form_type FROM form_manage WHERE form_manage_id = $1", formManageID).Scan(&formType)
	if err != nil {
		log.Println(err)
	}
	return formType
}

func getFormResultValues(formManageID string) FormResult {
	var result FormResult
	err := db.QueryRow("SELECT overwrite_title FROM form_manage WHERE form_manage_id = $1", formManageID).Scan(&result.OverwriteTitle)
	if err != nil {
		log.Println(err)
	}
	err = db.QueryRow("SELECT result_channel_id FROM form_manage WHERE form_manage_id = $1", formManageID).Scan(&result.ResultChannelID)
	if err != nil {
		log.Println(err)
	}
	err = db.QueryRow("SELECT accept_channel_id FROM form_manage WHERE form_manage_id = $1", formManageID).Scan(&result.AcceptChannelID)
	if err != nil {
		log.Println(err)
	}
	err = db.QueryRow("SELECT mods_can_comment FROM form_manage WHERE form_manage_id = $1", formManageID).Scan(&result.ModsCanComment)
	if err != nil {
		log.Println(err)
	}
	return result
}

func getFormOverwriteTitle(formManageID string) string {
	var overwriteTitle string
	err := db.QueryRow("SELECT overwrite_title FROM form_manage WHERE form_manage_id = $1", formManageID).Scan(&overwriteTitle)
	if err != nil {
		log.Println(err)
	}
	return overwriteTitle
}

func removeForm(formManageID string) {
	_, err := db.Exec("DELETE FROM form_manage WHERE form_manage_id = $1", formManageID)
	if err != nil {
		log.Println(err)
	}
}

func setAutoJoinRole(guildID string, option string, roleID string) bool {
	var role_exists bool
	var autojoinroles_exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT  1 FROM autojoinroles WHERE guild_id = $1)", guildID).Scan(&autojoinroles_exists)
	if err != nil {
		log.Println(err)
	}
	err = db.QueryRow("SELECT EXISTS (SELECT  1 FROM autojoinroles WHERE guild_id = $1 AND "+option+"_role IS NOT NULL AND "+option+"_role != '')", guildID).Scan(&role_exists)
	if err != nil {
		log.Println(err)
	}
	if autojoinroles_exists {
		_, err = db.Exec("UPDATE autojoinroles SET "+option+"_role = $1 WHERE guild_id = $2", roleID, guildID)
		if err != nil {
			log.Println(err)
		}
	} else {
		_, err = db.Exec("INSERT INTO autojoinroles (guild_id, "+option+"_role) VALUES ($1, $2)", guildID, roleID)
		if err != nil {
			log.Println(err)
		}
	}
	return role_exists
}

func purgeUnusedAutoJoinRoles(guildID string) {
	_, err := db.Exec("DELETE FROM autojoinroles WHERE guild_id = $1 AND user_role = '' OR user_role IS NULL AND bot_role = '' OR bot_role IS NULL", guildID)
	if err != nil {
		log.Println(err)
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
	err := db.QueryRow("SELECT "+isBotString+"_role FROM autojoinroles WHERE guild_id = $1", guildID).Scan(&role)
	if err != nil {
		log.Println(err, guildID)
	}
	return role
}

func toggleAutoPublish(guildID string, newsChannelID string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM autopublish WHERE guild_id = $1 AND news_channel_id = $2)", guildID, newsChannelID).Scan(&exists)
	if err != nil {
		log.Print(err)
	}
	if exists {
		_, err := db.Exec("DELETE FROM autopublish WHERE guild_id = $1 AND news_channel_id = $2", guildID, newsChannelID)
		if err != nil {
			log.Print(err)
		}
	} else {
		_, err := db.Exec("INSERT INTO autopublish (guild_id, news_channel_id) VALUES ($1, $2)", guildID, newsChannelID)
		if err != nil {
			log.Print(err)
		}
	}
	return exists
}

func isAutopublishEnabled(guildID string, newsChannelID string) bool {
	var enabled bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM autopublish WHERE guild_id = $1 AND news_channel_id = $2)", guildID, newsChannelID).Scan(&enabled)
	if err != nil {
		log.Print(err)
	}
	return enabled
}

func tryDeleteUnusedMessage(messageID string) {
	_, err := db.Exec("DELETE FROM form_manage WHERE message_id = $1", messageID)
	if err != nil {
		log.Println(err)
	}
}

func getAllSavedMessages() []MessageIDs {
	var savedMessages []MessageIDs
	rows, err := db.Query("SELECT message_id, channel_id FROM form_manage")
	if err != nil {
		log.Print(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var messageID, channelID string
		if err := rows.Scan(&messageID, &channelID); err != nil {
			log.Print(err)
			continue
		}
		savedMessages = append(savedMessages, MessageIDs{ID: messageID, ChannelID: channelID})
	}
	if err := rows.Err(); err != nil {
		log.Print(err)
	}
	return savedMessages
}
