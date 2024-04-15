package main

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func addTag(guildID, tagName, tagContent string) bool {
	var exists bool = true
	//TODO: add modify command
	id := uuid.New()
	for exists {
		id = uuid.New()
		err := db.QueryRow("SELECT EXISTS (SELECT  1 FROM tags WHERE guild_id = $1 AND tag_id = $2)", guildID, id).Scan(&exists)
		if err != nil {
			logrus.Error(err)
		}
	}
	_, err := db.Exec("INSERT INTO tags (guild_id, tag_name, tag_content, tag_id) VALUES ($1, $2, $3, $4)", guildID, tagName, tagContent, id)
	if err != nil {
		logrus.Error(err)
	}

	return exists
}
func removeTag(guildID string, tagID string) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT  1 FROM tags WHERE guild_id = $1 AND tag_id = $2)", guildID, tagID).Scan(&exists)
	if err != nil {
		logrus.Error(err)
	}
	if exists {
		_, err = db.Exec("DELETE FROM tags WHERE guild_id = $1 AND tag_id = $2", guildID, tagID)
		if err != nil {
			logrus.Error(err)
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
