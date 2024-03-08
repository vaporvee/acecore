package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

var cmd_cat Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:        "cat",
		Description: "Random cat pictures",
	},
	Interact: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		respondEmbed(i.Interaction, discordgo.MessageEmbed{
			Type:  discordgo.EmbedTypeImage,
			Color: hexToDecimal(color["primary"]),
			Image: &discordgo.MessageEmbedImage{
				URL: GetCatImageURL(),
			}}, false)
	},
}

type CatImage struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func GetCatImageURL() string {
	resp, err := http.Get("https://api.thecatapi.com/v1/images/search?format=json")
	if err != nil {
		log.Fatal("Error making GET request:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	var images []CatImage
	err = json.Unmarshal(body, &images)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
	}

	if len(images) == 0 {
		log.Fatal("No images found")
	}

	return images[0].URL
}
