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

var tags Tags
var filename string = "data.json"

func readTags() {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read tags: %v", err)
	}
	err = json.Unmarshal(bytes, &tags)
	if err != nil {
		log.Fatalf("Failed to read tags: %v", err)
	}
}

func writeTags() {
	jsonBytes, err := json.MarshalIndent(&tags, "", "  ")
	if err != nil {
		log.Fatalf("Failed to write tags: %v", err)
	}
	err = os.WriteFile(filename, jsonBytes, 0644)
	if err != nil {
		log.Fatalf("Failed to write tags: %v", err)
	}
}

func addTag(tags *Tags, tagKey string, tagValue string) {
	readTags()
	tags.Tags[tagKey] = tagValue
	writeTags()
}

func removeTag(tags *Tags, tagKey string) {
	readTags()
	delete(tags.Tags, tagKey)
	writeTags()
}

func (tags Tags) getTagKeys() []string {
	readTags()
	keys := make([]string, 0, len(tags.Tags))
	for k := range tags.Tags {
		keys = append(keys, k)
	}
	return keys
}

func modifyTag(tags *Tags, tagKey string, newTagValue string) {
	if _, exists := tags.Tags[tagKey]; exists {
		tags.Tags[tagKey] = newTagValue
	}
}

func debugTags() {
	addTag(&tags, "new_command", "a new command description")
}
