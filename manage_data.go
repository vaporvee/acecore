package main

import (
	"encoding/json"
	"log"
	"os"
)

//DATA WILL ONLY BE USED AS JSON FILE FOR TESTING. SYSTEM WILL BE REPLACED

type Tags struct {
	Tags map[string]string `json:"tags"`
}

func readTags(filename string) (*Tags, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var tags Tags
	err = json.Unmarshal(bytes, &tags)
	if err != nil {
		return nil, err
	}

	return &tags, nil
}

func writeTags(filename string, tags *Tags) error {
	jsonBytes, err := json.MarshalIndent(tags, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, jsonBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func addTag(tags *Tags, tagKey string, tagValue string) {
	tags.Tags[tagKey] = tagValue
}

func removeTag(tags *Tags, tagKey string) {
	delete(tags.Tags, tagKey)
}

func modifyTag(tags *Tags, tagKey string, newTagValue string) {
	if _, exists := tags.Tags[tagKey]; exists {
		tags.Tags[tagKey] = newTagValue
	}
}

func debugTags() {
	tags, err := readTags("data.json")
	if err != nil {
		log.Fatalf("Failed to read tags: %v", err)
	}

	addTag(tags, "new_command", "a new command description")
	removeTag(tags, "test_command")
	modifyTag(tags, "another_test_command", "updated command description")

	err = writeTags("data.json", tags)
	if err != nil {
		log.Fatalf("Failed to write tags: %v", err)
	}
}
