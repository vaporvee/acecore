package main

import (
	"log"

	"github.com/google/uuid"
)

func initTable() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS tags (
		tag_id TEXT NOT NULL,
		tag_name TEXT NOT NULL,
		tag_content TEXT,
		guild_id TEXT NOT NULL,
		PRIMARY KEY (tag_id,guild_id)
	);`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func addTag(guildID, tagName, tagContent string) bool {
	var exists bool
	id := uuid.New()
	err := db.QueryRow("SELECT EXISTS (SELECT  1 FROM tags WHERE guild_id = $1 AND tag_id = $2)", guildID, id).Scan(&exists)
	if err != nil {
		log.Println(err)
	}
	// TODO: add modify command
	for exists {
		id = uuid.New()
		err = db.QueryRow("SELECT EXISTS (SELECT  1 FROM tags WHERE guild_id = $1 AND tag_id = $2)", guildID, id).Scan(&exists)
		if err != nil {
			log.Println(err)
		}
	}
	_, err = db.Exec("INSERT INTO tags (guild_id, tag_name, tag_content, tag_id) VALUES ($1, $2, $3, $4)", guildID, tagName, tagContent, id)
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
