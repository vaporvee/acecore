package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var cmd_ping Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Returns the ping of the bot",
	},
	Interact: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		start := time.Now()

		client := http.Client{
			Timeout: 5 * time.Second,
		}

		resp, err := client.Get("https://discord.com/api/" + discordgo.APIVersion + "/gateway/bot")
		if err != nil {
			logrus.Error(err)
			return
		}
		defer resp.Body.Close()

		ping := time.Since(start)
		var ping_color string
		if ping.Milliseconds() < 200 {
			ping_color = "green"
		} else if ping.Milliseconds() < 400 {
			ping_color = "yellow"
		} else {
			ping_color = "red"
		}
		err = respondEmbed(i.Interaction, discordgo.MessageEmbed{
			Title:       s.State.User.Username + " ping",
			Description: fmt.Sprintf("# %.2fms", ping.Seconds()*1000),
			Type:        discordgo.EmbedTypeArticle,
			Color:       hexToDecimal(color[ping_color])}, true)
		if err != nil {
			logrus.Error(err)
		}
	},
	AllowDM: true,
}
