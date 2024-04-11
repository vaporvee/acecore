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
				AddEmbeds(discord.NewEmbedBuilder().SetDescription(cat.Fact).SetImage(cat.Image).SetColor(hexToDecimal(color["primary"])).Build()).
				Build())
			if err != nil {
				logrus.Error(err)
			}
		}
	},
}

type Cat struct {
	Image string `json:"image"`
	Fact  string `json:"fact"`
}

func GetCatImageURL() (Cat, error) {
	resp, err := http.Get("https://some-random-api.com/animal/cat")
	var cat Cat
	if err != nil {
		return cat, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return cat, err
	}
	err = json.Unmarshal(body, &cat)
	return cat, err
}
