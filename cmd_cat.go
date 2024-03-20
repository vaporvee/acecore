package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var cmd_cat Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:        "cat",
		Description: "Random cat pictures",
	},
	Interact: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		cat, err := GetCatImageURL()
		if err == nil {
			err := respondEmbed(i.Interaction, discordgo.MessageEmbed{
				Type:  discordgo.EmbedTypeImage,
				Color: hexToDecimal(color["primary"]),
				Image: &discordgo.MessageEmbedImage{
					URL: cat,
				}}, false)
			if err != nil {
				logrus.Error(err)
			}
		} else {
			logrus.Error(err)
		}
	},
	AllowDM: true,
}

type CatImage struct {
	ID string `json:"_id"`
}

func GetCatImageURL() (string, error) {
	resp, err := http.Get("https://cataas.com/cat?json=true")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var images CatImage
	err = json.Unmarshal(body, &images)

	return "https://cataas.com/cat/" + images.ID, err
}
