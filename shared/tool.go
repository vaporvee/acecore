package shared

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/bot"
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

// Why does the golang compiler care about commands??
//
//go:embed form_templates/*.json
var FormTemplates embed.FS

func GetModalByFormID(formID string) ModalJson {
	var modal ModalJson
	if formID == "" {
		return modal
	}
	entries, err := FormTemplates.ReadDir("form_templates")
	if err != nil {
		logrus.Error(err)
		return modal
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), formID) {
			jsonFile, err := FormTemplates.ReadFile("form_templates/" + entry.Name())
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

func JsonStringBuildModal(userID string, manageID string, formID string, overwrite ...string) discord.ModalCreate {
	var modal ModalJson = GetModalByFormID(formID)
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

	return discord.ModalCreate{
		CustomID:   "form:" + manageID + ":" + userID,
		Title:      modal.Title,
		Components: components,
	}

}

func IsIDRole(c bot.Client, guildID snowflake.ID, id snowflake.ID) bool {
	_, err1 := c.Rest().GetMember(guildID, id)
	if err1 == nil {
		return false
	}
	roles, err2 := c.Rest().GetRoles(guildID)
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
