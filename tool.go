package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
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
	if overwrite[0] != "" {
		modal.Title = overwrite[0]
	}
	err := bot.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   manageID + ":" + interaction.Member.User.ID,
			Title:      modal.Title,
			Components: components,
		},
	})
	if err != nil {
		log.Print(err)
	}
}

func getModalByFormID(formID string) ModalJson {
	var modal ModalJson
	//TODO: add custom forms
	entries, err := os.ReadDir("./form_templates")
	if err != nil {
		log.Print(err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), formID) {
			json_file, err := os.ReadFile("./form_templates/" + entry.Name())
			if err != nil {
				log.Print(err)
			}
			json.Unmarshal(json_file, &modal)
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

func respond(interaction *discordgo.Interaction, content string, ephemeral bool) {
	var flag discordgo.MessageFlags
	if ephemeral {
		flag = discordgo.MessageFlagsEphemeral
	}
	bot.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   flag,
		},
	})
}

func respondEmbed(interaction *discordgo.Interaction, embed discordgo.MessageEmbed, ephemeral bool) {
	var flag discordgo.MessageFlags
	if ephemeral {
		flag = discordgo.MessageFlagsEphemeral
	}
	bot.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: flag,
			Embeds: []*discordgo.MessageEmbed{
				&embed,
			},
		},
	})
}

func checkMessageNotExists(channelID, messageID string) bool {
	_, err := bot.ChannelMessage(channelID, messageID)
	if err != nil {
		return true
	}
	return false
}

func findAndDeleteUnusedMessages() {
	for _, message := range getAllSavedMessages() {
		if checkMessageNotExists(message.ChannelID, message.ID) {
			tryDeleteUnusedMessage(message.ID)
		}
	}
}
