package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/sirupsen/logrus"
)

var cmd_cat = Command{
	Definition: discord.SlashCommandCreate{
		Name:        "cat",
		Description: "Random cat pictures",
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeGuild,
			discord.InteractionContextTypePrivateChannel,
			discord.InteractionContextTypeBotDM,
		},
		IntegrationTypes: []discord.ApplicationIntegrationType{
			discord.ApplicationIntegrationTypeGuildInstall,
			discord.ApplicationIntegrationTypeUserInstall,
		},
	},
	Interact: func(e *events.ApplicationCommandInteractionCreate) {
		cat, err := GetCatImageURL()
		if err == nil {
			err = e.CreateMessage(discord.NewMessageCreateBuilder().
				AddEmbeds(discord.NewEmbedBuilder().SetImage(cat).SetColor(hexToDecimal(color["primary"])).Build()).
				Build())
			if err != nil {
				logrus.Error(err)
			}
		}
	},
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
