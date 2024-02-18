package main

import (
	"log"
)

func addTag(guildID, tagName, tagContent string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT  1 FROM tags WHERE guild_id = $1 AND tag_name = $2)", guildID, tagName).Scan(&exists)
	if err != nil {
		log.Println(err)
	}
	// If the tag exists it updates it but TODO: needs to return a discord message to use the modify command with autocomplete
	if exists {
		_, err = db.Exec("UPDATE tags SET tag_content = $1 WHERE guild_id = $2 AND tag_name = $3", tagContent, guildID, tagName)
		if err != nil {
			log.Println(err)
		}
	} else {
		_, err = db.Exec("INSERT INTO tags (guild_id, tag_name, tag_content) VALUES ($1, $2, $3)", guildID, tagName, tagContent)
		if err != nil {
			log.Println(err)
		}
	}
	return exists
}

func removeTag(guildID, tagContent string) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT  1 FROM tags WHERE guild_id = $1 AND tag_content = $2)", guildID, tagContent).Scan(&exists) // that is so dumb next commit i will make the tag IDs
	if err != nil {
		log.Println(err)
	}
	if exists {
		_, err = db.Exec("DELETE FROM tags WHERE guild_id = $1 AND tag_content = $2", guildID, tagContent)
		if err != nil {
			log.Println(err)
		}
	}
}

func getTagKeys(guildID string) ([]string, error) {
	var keys []string
	rows, err := db.Query("SELECT tag_name FROM tags WHERE guild_id = $1", guildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}

func getTag(guildID, tagName string) (string, error) {
	var tagContent string
	err := db.QueryRow("SELECT tag_content FROM tags WHERE guild_id = $1 AND tag_name = $2", guildID, tagName).Scan(&tagContent)
	return tagContent, err
}
