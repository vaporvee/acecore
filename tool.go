package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
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

func jsonStringShowModal(userID string, manageID string, formID string, overwrite ...string) discord.InteractionResponse {
	var modal ModalJson = getModalByFormID(formID)
	var components []discord.ContainerComponent
	for index, component := range modal.Form {
		var style discord.TextInputStyle = discord.TextInputStyleShort
		if component.IsParagraph {
			style = discord.TextInputStyleParagraph
		}
		components = append(components, discord.ActionRowComponent{
			discord.TextInputComponent{
				CustomID:    fmt.Sprint(index),
				Label:       component.Label,
				Style:       style,
				Placeholder: component.Placeholder,
				Required:    component.Required,
				MaxLength:   component.MaxLength,
				MinLength:   &component.MinLength,
				Value:       component.Value,
			},
		})
	}
	if overwrite != nil && overwrite[0] != "" {
		modal.Title = overwrite[0]
	}

	return discord.InteractionResponse{
		Type: discord.InteractionResponseTypeModal,
		Data: &discord.ModalCreate{
			CustomID:   manageID + ":" + userID,
			Title:      modal.Title,
			Components: components,
		},
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

func getHighestRole(guildID string) (*discord.Role, error) {
	botmember, err := client.Rest().GetMember(snowflake.MustParse(guildID), app.Bot.ID)
	if err != nil {
		return nil, err
	}
	roles, err := client.Rest().GetRoles(snowflake.MustParse(guildID))
	if err != nil {
		return nil, err
	}
	var highestRole *discord.Role
	for _, roleID := range botmember.RoleIDs {
		for _, role := range roles {
			if role.ID == roleID {
				if highestRole == nil || role.Position > highestRole.Position {
					highestRole = &role
				}
				break
			}
		}
	}
	return highestRole, nil
}

func ptr(s string) *string { return &s }

func hexToDecimal(hexColor string) int {
	hexColor = strings.TrimPrefix(hexColor, "#")
	decimal, err := strconv.ParseInt(hexColor, 16, 64)
	if err != nil {
		return 0
	}
	return int(decimal)
}

func decimalToHex(decimal int) string {
	hexString := strconv.FormatInt(int64(decimal), 16)
	return hexString
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

func findAndDeleteUnusedMessages() {
	for _, message := range getAllSavedMessages() {
		_, err := client.Rest().GetMessage(snowflake.MustParse(message.ChannelID), snowflake.MustParse(message.ID))
		if err != nil {
			tryDeleteUnusedMessage(message.ID)
		}
	}
}

func isIDRole(guildID snowflake.ID, id snowflake.ID) bool {
	_, err1 := client.Rest().GetMember(guildID, id)
	if err1 == nil {
		return false
	}
	roles, err2 := client.Rest().GetRoles(guildID)
	if err2 == nil {
		for _, role := range roles {
			if role.ID == id {
				return true
			}
		}
	}

	logrus.Error(err1)
	logrus.Error(err2)
	return false
}
