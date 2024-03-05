package main

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ModalJsonField struct {
	Label       string `json:"label"`
	IsParagraph bool   `json:"is_paragraph"`
	Value       string `json:"value"`
	Required    bool   `json:"required"`
	MinLength   int    `json:"min_length"`
	MaxLength   int    `json:"max_length"`
}

type ModalJson struct {
	FormType string           `json:"form_type"`
	Title    string           `json:"title"`
	Form     []ModalJsonField `json:"form"`
}

func jsonStringShowModal(jsonString string, id string) {
	var modal ModalJson
	json.Unmarshal([]byte(jsonString), &modal)
}

func getHighestRole(s *discordgo.Session, guildID string) (*discordgo.Role, error) {
	botMember, err := s.GuildMember(guildID, s.State.User.ID)
	if err != nil {
		return nil, err
	}
	roles, err := s.GuildRoles(guildID)
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

func respond(s *discordgo.Session, interaction *discordgo.Interaction, content string, ephemeral bool) {
	var flag discordgo.MessageFlags
	if ephemeral {
		flag = discordgo.MessageFlagsEphemeral
	}
	s.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   flag,
		},
	})
}

func respondEmbed(s *discordgo.Session, interaction *discordgo.Interaction, embed discordgo.MessageEmbed, ephemeral bool) {
	var flag discordgo.MessageFlags
	if ephemeral {
		flag = discordgo.MessageFlagsEphemeral
	}
	s.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: flag,
			Embeds: []*discordgo.MessageEmbed{
				&embed,
			},
		},
	})
}
