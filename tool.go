package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type ModalJsonField struct {
	Label       string `json:"label"`
	IsParagraph bool   `json:"is_paragraph"`
	Value       string `json:"value"`
	Required    bool   `json:"required"`
	Placeholder string `json:"placeholder"`
	MinLength   int    `json:"min_length"`
	MaxLength   int    `json:"max_length"`
}

type ModalJson struct {
	FormType string           `json:"form_type"`
	Title    string           `json:"title"`
	Form     []ModalJsonField `json:"form"`
}

type MessageIDs struct {
	ID        string
	ChannelID string
}

func jsonStringShowModal(interaction *discordgo.Interaction, manageID string, formID string, overwrite ...string) {
	var modal ModalJson = getModalByFormID(formID)
	var components []discordgo.MessageComponent
	for index, component := range modal.Form {
		var style discordgo.TextInputStyle = discordgo.TextInputShort
		if component.IsParagraph {
			style = discordgo.TextInputParagraph
		}
		components = append(components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    fmt.Sprint(index),
					Label:       component.Label,
					Style:       style,
					Placeholder: component.Placeholder,
					Required:    component.Required,
					MaxLength:   component.MaxLength,
					MinLength:   component.MinLength,
					Value:       component.Value,
				},
			},
		})
	}
	if overwrite != nil && overwrite[0] != "" {
		modal.Title = overwrite[0]
	}
	var err error
	if modal.Title != "" && components != nil {
		err = bot.InteractionRespond(interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID:   manageID + ":" + interaction.Member.User.ID,
				Title:      modal.Title,
				Components: components,
			},
		})
	}
	if err != nil {
		logrus.Error(err)
	}
}

// Why does the golang compiler care about commands??
//
//go:embed form_templates/*.json
var formTemplates embed.FS

func getModalByFormID(formID string) ModalJson {
	var modal ModalJson
	entries, err := formTemplates.ReadDir("form_templates")
	if err != nil {
		logrus.Error(err)
		return modal
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), formID) {
			jsonFile, err := formTemplates.ReadFile("form_templates/" + entry.Name())
			if err != nil {
				logrus.Error(err)
				continue
			}
			err = json.Unmarshal(jsonFile, &modal)
			if err != nil {
				logrus.Error(err)
				continue
			}
			break
		}
	}
	return modal
}

func getHighestRole(guildID string) (*discordgo.Role, error) {
	botMember, err := bot.GuildMember(guildID, bot.State.User.ID)
	if err != nil {
		return nil, err
	}
	roles, err := bot.GuildRoles(guildID)
	if err != nil {
		return nil, err
	}

	var highestRole *discordgo.Role
	for _, roleID := range botMember.Roles {
		for _, role := range roles {
			if role.ID == roleID {
				if highestRole == nil || role.Position > highestRole.Position {
					highestRole = role
				}
				break
			}
		}
	}
	return highestRole, nil
}

func int64Ptr(i int64) *int64 {
	return &i
}

func hexToDecimal(hexColor string) int {
	// Remove the hash symbol if it's present
	hexColor = strings.TrimPrefix(hexColor, "#")
	decimal, err := strconv.ParseInt(hexColor, 16, 64)
	if err != nil {
		return 0
	}
	return int(decimal)
}

func simpleGetFromAPI(key string, url string) interface{} {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error("Error creating request:", err)
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("Error making request:", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("Error reading response body:", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		logrus.Error("Error decoding JSON:", err)
	}
	return result[key]
}

func respond(interaction *discordgo.Interaction, content string, ephemeral bool) error {
	var flag discordgo.MessageFlags
	if ephemeral {
		flag = discordgo.MessageFlagsEphemeral
	}
	err := bot.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   flag,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func respondEmbed(interaction *discordgo.Interaction, embed discordgo.MessageEmbed, ephemeral bool) error {
	var flag discordgo.MessageFlags
	if ephemeral {
		flag = discordgo.MessageFlagsEphemeral
	}
	err := bot.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: flag,
			Embeds: []*discordgo.MessageEmbed{
				&embed,
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func findAndDeleteUnusedMessages() {
	for _, message := range getAllSavedMessages() {
		_, err := bot.ChannelMessage(message.ChannelID, message.ID)
		if err != nil {
			tryDeleteUnusedMessage(message.ID)
		}
	}
}
