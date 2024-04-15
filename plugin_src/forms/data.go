package main

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type FormResult struct {
	OverwriteTitle    string
	ResultChannelID   string
	AcceptChannelID   string
	CommentCategoryID string
	ModeratorID       string
}

type MessageIDs struct {
	ID        string
	ChannelID string
}

func getFormManageIdExists(id uuid.UUID) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM form_manage WHERE form_manage_id = $1)", id).Scan(&exists)
	if err != nil {
		logrus.Error(err)
	}
	return exists
}

func addFormButton(guildID string, channelID string, messageID string, formManageID string, formType string, resultChannelID string, overwriteTitle string, acceptChannelID string, commentCategory string, moderator_id string) {
	_, err := db.Exec(
		`INSERT INTO form_manage (
			guild_id, 
			form_manage_id, 
			channel_id, 
			message_id, 
			form_type, 
			result_channel_id, 
			overwrite_title, 
			accept_channel_id, 
			comment_category,
			moderator_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		guildID, formManageID, channelID, messageID, formType, resultChannelID, overwriteTitle, acceptChannelID, commentCategory, moderator_id)
	if err != nil {
		logrus.Error(err)
	}
}

func getFormManageIDs() []string {
	if db == nil {
		return nil
	}
	var IDs []string
	rows, err := db.Query("SELECT form_manage_id FROM form_manage")
	if err != nil {
		logrus.Error(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			logrus.Error(err)
			return nil
		}
		IDs = append(IDs, id)
	}

	if err := rows.Err(); err != nil {
		logrus.Error(err)
		return nil
	}
	return IDs
}

func getFormType(formManageID string) string {
	var formType string
	err := db.QueryRow("SELECT form_type FROM form_manage WHERE form_manage_id = $1", formManageID).Scan(&formType)
	if err != nil {
		logrus.Error(err)
	}
	return formType
}

func getFormResultValues(formManageID string) FormResult {
	var result FormResult
	err := db.QueryRow("SELECT overwrite_title FROM form_manage WHERE form_manage_id = $1", formManageID).Scan(&result.OverwriteTitle)
	if err != nil {
		logrus.Error(err)
	}
	err = db.QueryRow("SELECT result_channel_id FROM form_manage WHERE form_manage_id = $1", formManageID).Scan(&result.ResultChannelID)
	if err != nil {
		logrus.Error(err)
	}
	err = db.QueryRow("SELECT accept_channel_id FROM form_manage WHERE form_manage_id = $1", formManageID).Scan(&result.AcceptChannelID)
	if err != nil {
		logrus.Error(err)
	}
	err = db.QueryRow("SELECT comment_category FROM form_manage WHERE form_manage_id = $1", formManageID).Scan(&result.CommentCategoryID)
	if err != nil {
		logrus.Error(err)
	}
	err = db.QueryRow("SELECT moderator_id FROM form_manage WHERE form_manage_id = $1", formManageID).Scan(&result.ModeratorID)
	if err != nil {
		logrus.Error(err)
	}
	return result
}

func getFormOverwriteTitle(formManageID string) string {
	var overwriteTitle string
	err := db.QueryRow("SELECT overwrite_title FROM form_manage WHERE form_manage_id = $1", formManageID).Scan(&overwriteTitle)
	if err != nil {
		logrus.Error(err)
	}
	return overwriteTitle
}

func updateFormCommentCategory(formManageID string, comment_category string) {
	_, err := db.Exec("UPDATE form_manage SET comment_category = $1 WHERE form_manage_id = $2", comment_category, formManageID)
	if err != nil {
		logrus.Error(err)
	}
}

func tryDeleteUnusedMessage(messageID string) {
	_, err := db.Exec("DELETE FROM form_manage WHERE message_id = $1", messageID)
	if err != nil {
		logrus.Error(err)
	}
}

func getAllSavedMessages() []MessageIDs {
	var savedMessages []MessageIDs
	rows, err := db.Query("SELECT message_id, channel_id FROM form_manage")
	if err != nil {
		logrus.Error(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var messageID, channelID string
		if err := rows.Scan(&messageID, &channelID); err != nil {
			logrus.Error(err)
			continue
		}
		savedMessages = append(savedMessages, MessageIDs{ID: messageID, ChannelID: channelID})
	}
	if err := rows.Err(); err != nil {
		logrus.Error(err)
	}
	return savedMessages
}
