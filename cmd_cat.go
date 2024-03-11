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
	ID string `json:"_id"`
}

func GetCatImageURL() string {
	resp, err := http.Get("https://cataas.com/cat?json=true")
	if err != nil {
		log.Print("Error making GET request:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print("Error reading response body:", err)
	}

	var images CatImage
	err = json.Unmarshal(body, &images)
	if err != nil {
		log.Print("Error unmarshalling JSON:", err)
	}

	return "https://cataas.com/cat/" + images.ID
}
