package main

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/sirupsen/logrus"
)

var cmd_getemoji Command = Command{
	Definition: discord.SlashCommandCreate{
		Name:        "add-emoji",
		Description: "Add an external emoji directly to the server.",
		Options: []discord.ApplicationCommandOption{
			&discord.ApplicationCommandOptionString{
				Name:        "emoji",
				Description: "The emoji you want to add",
				Required:    true,
			},
		},
	},
	Interact: func(e *events.ApplicationCommandInteractionCreate) {
		emojiRegex := regexp.MustCompile(`<(.+):(\d+)>`)
		emojistring := emojiRegex.FindString(e.SlashCommandInteractionData().String("emoji"))
		logrus.Debug(emojistring)
		emojiArray := strings.Split(emojistring, ":")
		logrus.Debug(emojiArray)
		var emojiName string
		var emojiID string
		var emojiFileName string
		if len(emojiArray) > 1 {
			emojiName = strings.TrimSuffix(emojiArray[1], ">")
			emojiID = strings.TrimSuffix(emojiArray[2], ">")
		}
		imageType, emojiRead := getEmoji(emojiID)
		emojiData, err := io.ReadAll(emojiRead)
		if err != nil {
			logrus.Error(err)
		}
		_, err = e.Client().Rest().CreateEmoji(*e.GuildID(), discord.EmojiCreate{
			Name: emojiName,
			Image: discord.Icon{
				Type: imageType,
				Data: emojiData,
			},
		})
		if err != nil {
			logrus.Error(err)
		}
		if imageType == discord.IconTypeGIF {
			emojiFileName = emojiName + ".gif"
		} else {
			emojiFileName = emojiName + ".png"
		}
		err = e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContentf("Emoji %s sucessfully added to this server!", emojiName).SetFiles(discord.NewFile(emojiFileName, "The emoji that was picked", emojiRead)).SetEphemeral(true).
			Build())
		if err != nil {
			logrus.Error(err)
		}
	},
}

func getEmoji(emojiID string) (discord.IconType, io.Reader) {
	resp, err := http.Get("https://cdn.discordapp.com/emojis/" + emojiID)
	if err != nil {
		logrus.Error(err)
		return discord.IconTypePNG, nil
	}
	defer resp.Body.Close()
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return discord.IconTypePNG, nil
	}
	isAnimated := isGIFImage(imageData)
	if isAnimated {
		logrus.Debug("GIF")
		return discord.IconTypeGIF, bytes.NewReader(imageData)
	} else {
		logrus.Debug("PNG")
		return discord.IconTypePNG, bytes.NewReader(imageData)
	}
}
