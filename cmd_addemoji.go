package main

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/json"
	"github.com/sirupsen/logrus"
	"github.com/vaporvee/acecore/cmd"
)

var cmd_addemoji cmd.Command = cmd.Command{
	Definition: discord.SlashCommandCreate{
		Name:                     "add-emoji",
		Description:              "Add an external emoji directly to the server.",
		DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionManageGuildExpressions),
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
		emojiArray := strings.Split(emojistring, ":")
		var emojiName string
		var emojiID string
		var emojiFileName string
		if len(emojiArray) > 1 {
			emojiName = strings.TrimSuffix(emojiArray[1], ">")
			emojiID = strings.TrimSuffix(emojiArray[2], ">")
		}
		imageType, emojiReadBit64 := getEmoji(emojiID)
		emojiData, err := discord.NewIcon(imageType, emojiReadBit64)
		if err != nil {
			logrus.Error(err)
		}
		_, err = e.Client().Rest().CreateEmoji(*e.GuildID(), discord.EmojiCreate{
			Name:  emojiName,
			Image: *emojiData,
		})
		if err != nil {
			if strings.HasPrefix(err.Error(), "50035") {
				e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Failed adding emoji. Did you provide a correct one?").SetEphemeral(true).Build())
				return
			}
			if strings.HasPrefix(err.Error(), "50138") {
				e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Failed adding emoji. Unable to resize the emoji image.").SetEphemeral(true).Build())
				return
			}
			logrus.Error(err)
			return
		}
		if imageType == discord.IconTypeGIF {
			emojiFileName = emojiName + ".gif"
		} else {
			emojiFileName = emojiName + ".png"
		}
		_, emojiRead := getEmoji(emojiID) // for some reason any []bit variable thats used with NewIcon gets corrupted even when its redeclared in a new variable
		err = e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContentf("Emoji %s sucessfully added to this server!", emojiName).SetFiles(discord.NewFile(emojiFileName, "", emojiRead)).SetEphemeral(true).
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
		return discord.IconTypeGIF, bytes.NewReader(imageData)
	} else {
		return discord.IconTypePNG, bytes.NewReader(imageData)
	}
}
